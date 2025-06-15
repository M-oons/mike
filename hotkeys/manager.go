package hotkeys

import (
	"context"
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"github.com/m-oons/mike/config"
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
}

func NewManager(audioService AudioService, config []config.ConfigHotkey) *manager {
	return &manager{
		audioService: audioService,
		config:       config,
	}
}

func (m *manager) Start(ctx context.Context) error {
	if err := m.loadAPI(); err != nil {
		return fmt.Errorf("error loading Windows API: %w", err)
	}
	defer m.user32.Release()
	defer m.kernel32.Release()

	if err := m.register(); err != nil {
		m.unregister()
		return fmt.Errorf("error registering hotkeys: %w", err)
	}
	defer m.unregister()

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

func (m *manager) register() error {
	for i, confkey := range m.config {
		hotkey := hotkey{
			action: confkey.Action,
			key:    confkey.Key,
			ctrl:   confkey.Ctrl,
			shift:  confkey.Shift,
			alt:    confkey.Alt,
			win:    confkey.Win,
		}
		ret, _, err := m.registerHotkey.Call(0, uintptr(i+1), uintptr(hotkey.modifiers()), uintptr(hotkey.code()))
		if ret == 0 {
			return err
		}
	}

	return nil
}

func (m *manager) unregister() {
	for i := range m.config {
		m.unregisterHotkey.Call(0, uintptr(i+1))
	}
}

func (m *manager) loop() error {
	var msg MSG

	for {
		ret, _, _ := m.getMessage.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)

		if ret == 0 || msg.Message == WM_QUIT { // quit
			return nil
		}

		if ret == ^uintptr(0) { // error
			continue
		}

		if msg.Message == WM_HOTKEY && int16(msg.WParam) <= int16(len(m.config)) {
			hotkey := m.config[msg.WParam-1]
			switch strings.ToLower(hotkey.Action) {
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
