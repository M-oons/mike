package hotkeys

import "strings"

type Hotkey struct {
	Action string
	Key    string
	Ctrl   bool
	Shift  bool
	Alt    bool
	Win    bool
}

func (hotkey *Hotkey) Modifiers() int {
	modifiers := ModNoRepeat
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
	return modifiers
}

func (hotkey *Hotkey) KeyCode() int {
	keycode, ok := Keys[strings.ToLower(hotkey.Key)]
	if !ok {
		return -1
	}
	return keycode
}
