package hotkeys

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"github.com/m-oons/mike/internal/config"
	"golang.org/x/sys/windows"
)

const (
	WM_QUIT   = 0x0012
	WM_HOTKEY = 0x0312
)

type MSG struct {
	HWnd    uintptr
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      struct {
		X int32
		Y int32
	}
}

type AudioService interface {
	Mute() error
	Unmute() error
	ToggleMute() (bool, error)
}

type manager struct {
	audioService       AudioService
	config             []config.ConfigHotkey
	user32             *syscall.DLL
	kernel32           *syscall.DLL
	registerHotkey     *syscall.Proc
	unregisterHotkey   *syscall.Proc
	getMessage         *syscall.Proc
	postThreadMessage  *syscall.Proc
	getCurrentThreadId *syscall.Proc
	registeredHotkeys  map[uintptr]hotkey
}

func NewManager(audioService AudioService, config []config.ConfigHotkey) *manager {
	return &manager{
		audioService:      audioService,
		config:            config,
		registeredHotkeys: make(map[uintptr]hotkey),
	}
}

func (m *manager) Start(ctx context.Context) error {
	if err := m.loadAPI(); err != nil {
		return fmt.Errorf("error loading Windows API: %w", err)
	}
	defer m.user32.Release()
	defer m.kernel32.Release()

	if err := m.registerAll(); err != nil {
		m.unregisterAll()
		return fmt.Errorf("error registering hotkeys: %w", err)
	}
	defer m.unregisterAll()

	threadId, _, _ := m.getCurrentThreadId.Call()

	go func() {
		<-ctx.Done()
		m.postThreadMessage.Call(threadId, WM_QUIT, 0, 0)
	}()

	return m.loop()
}

func (m *manager) loadAPI() error {
	var err error
	m.user32, err = syscall.LoadDLL("user32")
	if err != nil {
		return fmt.Errorf("error loading 'user32' DLL: %w", err)
	}

	m.kernel32, err = syscall.LoadDLL("kernel32")
	if err != nil {
		return fmt.Errorf("error loading 'kernel32' DLL: %w", err)
	}

	m.registerHotkey, err = m.user32.FindProc("RegisterHotKey")
	if err != nil {
		return fmt.Errorf("error finding user32 procedure 'RegisterHotKey': %w", err)
	}

	m.unregisterHotkey, err = m.user32.FindProc("UnregisterHotKey")
	if err != nil {
		return fmt.Errorf("error finding user32 procedure 'UnregisterHotKey': %w", err)
	}

	m.getMessage, err = m.user32.FindProc("GetMessageW")
	if err != nil {
		return fmt.Errorf("error finding user32 procedure 'GetMessageW': %w", err)
	}

	m.postThreadMessage, err = m.user32.FindProc("PostThreadMessageW")
	if err != nil {
		return fmt.Errorf("error finding user32 procedure 'PostThreadMessageW': %w", err)
	}

	m.getCurrentThreadId, err = m.kernel32.FindProc("GetCurrentThreadId")
	if err != nil {
		return fmt.Errorf("error finding kernel32 procedure 'GetCurrentThreadId': %w", err)
	}

	return nil
}

func (m *manager) registerAll() error {
	registeredCombos := make(map[string]struct{})
	var specificHotkeys []config.ConfigHotkey
	var genericHotkeys []config.ConfigHotkey

	// separate hotkeys into groups - specific/generic
	for _, confHotkey := range m.config {
		if !confHotkey.Ctrl && !confHotkey.Shift && !confHotkey.Alt && !confHotkey.Win { // no modifiers
			genericHotkeys = append(genericHotkeys, confHotkey)
		} else {
			specificHotkeys = append(specificHotkeys, confHotkey)
		}
	}

	// register specific hotkeys first
	for _, confHotkey := range specificHotkeys {
		hotkey := hotkey{
			action: confHotkey.Action,
			key:    confHotkey.Key,
			ctrl:   confHotkey.Ctrl,
			shift:  confHotkey.Shift,
			alt:    confHotkey.Alt,
			win:    confHotkey.Win,
		}

		keyCode := hotkey.code()
		if keyCode == -1 { // skip invalid hotkey
			continue
		}

		if err := m.register(hotkey); err != nil {
			return err
		}

		registeredCombo := fmt.Sprintf("%d|%d", keyCode, hotkey.modifiers())
		registeredCombos[registeredCombo] = struct{}{}
	}

	// register generic hotkeys - skip already registered modifier combinations
	for _, confHotkey := range genericHotkeys {
		hotkey := hotkey{
			action: confHotkey.Action,
			key:    confHotkey.Key,
		}

		keyCode := hotkey.code()
		if keyCode == -1 {
			continue
		}

		// brute force register all possible modifier combinations
		for combo := range modCombinations {
			comboHotkey := hotkey
			comboHotkey.ctrl = (combo & modCtrl) != 0
			comboHotkey.shift = (combo & modShift) != 0
			comboHotkey.alt = (combo & modAlt) != 0
			comboHotkey.win = (combo & modWin) != 0

			registeredCombo := fmt.Sprintf("%d|%d", keyCode, comboHotkey.modifiers())
			if _, ok := registeredCombos[registeredCombo]; ok { // modifier combination already registered
				continue
			}

			if err := m.register(comboHotkey); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *manager) register(hotkey hotkey) error {
	hotkeyID := uintptr(len(m.registeredHotkeys) + 1)
	ret, _, err := m.registerHotkey.Call(
		0,
		hotkeyID,
		uintptr(hotkey.modifiers()),
		uintptr(hotkey.code()),
	)

	if ret == 0 && !errors.Is(err, windows.ERROR_HOTKEY_ALREADY_REGISTERED) {
		return err
	}

	m.registeredHotkeys[hotkeyID] = hotkey

	return nil
}

func (m *manager) unregisterAll() {
	for hotkeyID := range m.registeredHotkeys {
		m.unregisterHotkey.Call(
			0,
			hotkeyID,
		)
	}
}

func (m *manager) loop() error {
	var msg MSG

	for {
		ret, _, _ := m.getMessage.Call(
			uintptr(unsafe.Pointer(&msg)),
			0,
			0,
			0,
		)

		if ret == 0 || msg.Message == WM_QUIT { // quit
			return nil
		}

		if ret == ^uintptr(0) { // error
			continue
		}

		if msg.Message == WM_HOTKEY {
			if hotkey, ok := m.registeredHotkeys[msg.WParam]; ok {
				switch strings.ToLower(hotkey.action) {
				case "mute":
					m.audioService.Mute()

				case "unmute":

					m.audioService.Unmute()
				case "toggle":
					m.audioService.ToggleMute()
				}
			}
		}
	}
}
