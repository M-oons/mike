package hotkeys

import (
	"strings"
	"syscall"
	"unsafe"

	"github.com/m-oons/mike/actions"
	"github.com/m-oons/mike/config"
)

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
	"0":       48,
	"1":       49,
	"2":       50,
	"3":       51,
	"4":       52,
	"5":       53,
	"6":       54,
	"7":       55,
	"8":       56,
	"9":       57,
	";":       59,
	"=":       61,
	"a":       65,
	"b":       66,
	"c":       67,
	"d":       68,
	"e":       69,
	"f":       70,
	"g":       71,
	"h":       72,
	"i":       73,
	"j":       74,
	"k":       75,
	"l":       76,
	"m":       77,
	"n":       78,
	"o":       79,
	"p":       80,
	"q":       81,
	"r":       82,
	"s":       83,
	"t":       84,
	"u":       85,
	"v":       86,
	"w":       87,
	"x":       88,
	"y":       89,
	"z":       90,
	"numpad0": 96,
	"numpad1": 97,
	"numpad2": 98,
	"numpad3": 99,
	"numpad4": 100,
	"numpad5": 101,
	"numpad6": 102,
	"numpad7": 103,
	"numpad8": 104,
	"numpad9": 105,
	"numpad*": 106,
	"numpad+": 107,
	"numpad-": 109,
	"numpad.": 110,
	"numpad/": 111,
	"f1":      112,
	"f2":      113,
	"f3":      114,
	"f4":      115,
	"f5":      116,
	"f6":      117,
	"f7":      118,
	"f8":      119,
	"f9":      120,
	"f10":     121,
	"f11":     122,
	"f12":     123,
	"f13":     124,
	"f14":     125,
	"f15":     126,
	"f16":     127,
	"f17":     128,
	"f18":     129,
	"f19":     130,
	"f20":     131,
	"f21":     132,
	"f22":     133,
	"f23":     134,
	"f24":     135,
	"#":       163,
	"-":       173,
	",":       188,
	".":       190,
	"/":       191,
	"`":       192,
	"[":       219,
	"\\":      220,
	"]":       221,
	"'":       222,
}

func RegisterHotkeys() {
	user32 := syscall.MustLoadDLL("user32")

	reghotkey := user32.MustFindProc("RegisterHotKey")
	getmsg := user32.MustFindProc("GetMessageW")

	user32.Release()

	for i, confkey := range config.Current.Hotkeys {
		hotkey := Hotkey{
			Action:   confkey.Action,
			Key:      confkey.Key,
			Ctrl:     confkey.Ctrl,
			Shift:    confkey.Shift,
			Alt:      confkey.Alt,
			Win:      confkey.Win,
			NoRepeat: confkey.NoRepeat,
		}
		reghotkey.Call(0, uintptr(i+1), uintptr(hotkey.Modifiers()), uintptr(hotkey.KeyCode()))
	}

	for {
		msg := &MSG{}
		getmsg.Call(uintptr(unsafe.Pointer(msg)), 0, 0, 0)

		if msg.WPARAM <= int16(len(config.Current.Hotkeys)) {
			hotkey := config.Current.Hotkeys[msg.WPARAM-1]
			switch strings.ToLower(hotkey.Action) {
			case "mute":
				actions.Mute()
			case "unmute":
				actions.Unmute()
			case "toggle":
				actions.ToggleMute()
			}
		}
	}
}
