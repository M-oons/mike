package hotkeys

import "strings"

type Hotkey struct {
	Action   string
	Key      string
	Ctrl     bool
	Shift    bool
	Alt      bool
	Win      bool
	NoRepeat bool
}

func (hotkey *Hotkey) Modifiers() int {
	modifiers := 0
	if hotkey.Ctrl {
		modifiers += ModCtrl
	}
	if hotkey.Shift {
		modifiers += ModShift
	}
	if hotkey.Alt {
		modifiers += ModAlt
	}
	if hotkey.Win {
		modifiers += ModWin
	}
	if hotkey.NoRepeat {
		modifiers += ModNoRepeat // holding down the hotkey won't continuously trigger keybind
	}
	return modifiers
}

func (hotkey *Hotkey) KeyCode() int {
	keycode, ok := Keys[strings.ToLower(hotkey.Key)]
	if !ok {
		return -1
	}
	return keycode
}
