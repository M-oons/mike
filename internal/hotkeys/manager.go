package hotkeys

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/m-oons/mike/internal/config"
	"github.com/m-oons/mike/internal/windows/kernel32"
	"github.com/m-oons/mike/internal/windows/user32"
	"golang.org/x/sys/windows"
)

const (
	WM_QUIT   = 0x0012
	WM_HOTKEY = 0x0312
)

type AudioService interface {
	Mute() error
	Unmute() error
	ToggleMute() (bool, error)
}

type manager struct {
	audioService AudioService
	config       []config.ConfigHotkey

	registeredHotkeys map[uintptr]hotkey
}

func NewManager(audioService AudioService, config []config.ConfigHotkey) *manager {
	return &manager{
		audioService:      audioService,
		config:            config,
		registeredHotkeys: make(map[uintptr]hotkey),
	}
}

func (m *manager) Start(ctx context.Context) error {
	if err := m.registerAll(); err != nil {
		m.unregisterAll()
		return fmt.Errorf("error registering hotkeys: %w", err)
	}
	defer m.unregisterAll()

	threadID, _ := kernel32.GetCurrentThreadID()

	go func() {
		<-ctx.Done()
		user32.PostThreadMessage(threadID, WM_QUIT)
	}()

	return m.loop()
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
	ret, err := user32.RegisterHotKey(
		hotkeyID,
		hotkey.code(),
		hotkey.modifiers(),
	)

	if ret == 0 && !errors.Is(err, windows.ERROR_HOTKEY_ALREADY_REGISTERED) {
		return err
	}

	m.registeredHotkeys[hotkeyID] = hotkey

	return nil
}

func (m *manager) unregisterAll() {
	for hotkeyID := range m.registeredHotkeys {
		user32.UnregisterHotKey(hotkeyID)
	}
}

func (m *manager) loop() error {
	for {
		msg, ret, _ := user32.GetMessage()

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
