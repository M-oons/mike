package hotkeys

import "strings"

const (
	modAlt      = 0x0001
	modCtrl     = 0x0002
	modShift    = 0x0004
	modWin      = 0x0008
	modNoRepeat = 0x4000 // holding down the hotkey won't continuously trigger keybind
)

type hotkey struct {
	action string
	key    string
	ctrl   bool
	shift  bool
	alt    bool
	win    bool
}

func (hotkey *hotkey) modifiers() int {
	modifiers := modNoRepeat
	if hotkey.ctrl {
		modifiers += modCtrl
	}
	if hotkey.shift {
		modifiers += modShift
	}
	if hotkey.alt {
		modifiers += modAlt
	}
	if hotkey.win {
		modifiers += modWin
	}

	return modifiers
}

func (hotkey *hotkey) code() int {
	keycode, ok := keys[strings.ToLower(hotkey.key)]
	if !ok {
		return -1
	}

	return keycode
}
