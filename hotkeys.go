package main

import (
	"strings"
	"syscall"
	"time"
	"unsafe"
)

type Hotkey struct {
	Action   string `json:"action"`
	Key      string `json:"key"`
	Ctrl     bool   `json:"ctrl"`
	Shift    bool   `json:"shift"`
	Alt      bool   `json:"alt"`
	Win      bool   `json:"win"`
	NoRepeat bool   `json:"norepeat"`
}

type MSG struct {
	HWND   uintptr
	UINT   uintptr
	WPARAM int16
	LPARAM int64
	DWORD  int32
	POINT  struct {
		X int64
		Y int64
	}
}

const (
	ModAlt      = 0x0001
	ModCtrl     = 0x0002
	ModShift    = 0x0004
	ModWin      = 0x0008
	ModNoRepeat = 0x4000
)

var keys map[string]int = map[string]int{
	"0":   48,
	"1":   49,
	"2":   50,
	"3":   51,
	"4":   52,
	"5":   53,
	"6":   54,
	"7":   55,
	"8":   56,
	"9":   57,
	"A":   65,
	"B":   66,
	"C":   67,
	"D":   68,
	"E":   69,
	"F":   70,
	"G":   71,
	"H":   72,
	"I":   73,
	"J":   74,
	"K":   75,
	"L":   76,
	"M":   77,
	"N":   78,
	"O":   79,
	"P":   80,
	"Q":   81,
	"R":   82,
	"S":   83,
	"T":   84,
	"U":   85,
	"V":   86,
	"W":   87,
	"X":   88,
	"Y":   89,
	"Z":   90,
	"F1":  112,
	"F2":  113,
	"F3":  114,
	"F4":  115,
	"F5":  116,
	"F6":  117,
	"F7":  118,
	"F8":  119,
	"F9":  120,
	"F10": 121,
	"F11": 122,
	"F12": 123,
	"F13": 124,
	"F14": 125,
	"F15": 126,
	"F16": 127,
	"F17": 128,
	"F18": 129,
	"F19": 130,
	"F20": 131,
	"F21": 132,
	"F22": 133,
	"F23": 134,
	"F24": 135,
}

func RegisterHotkeys() {
	user32 := syscall.MustLoadDLL("user32")
	defer user32.Release()

	reghotkey := user32.MustFindProc("RegisterHotKey")
	getmsg := user32.MustFindProc("GetMessageW")

	for i, hotkey := range Conf.Hotkeys {
		keycode, ok := keys[strings.ToUpper(hotkey.Key)]
		if !ok {
			continue
		}
		reghotkey.Call(0, uintptr(i+1), uintptr(hotkey.Modifiers()), uintptr(keycode))
	}

	for {
		msg := &MSG{}
		getmsg.Call(uintptr(unsafe.Pointer(msg)), 0, 0, 0)

		if msg.WPARAM <= int16(len(Conf.Hotkeys)) {
			hotkey := Conf.Hotkeys[msg.WPARAM-1]
			switch strings.ToLower(hotkey.Action) {
			case "mute":
				Mute()
			case "unmute":
				Unmute()
			case "toggle":
				ToggleMute()
			}
		}

		time.Sleep(time.Millisecond * 50)
	}
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
