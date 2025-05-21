package hotkeys

import (
	"strings"
	"syscall"
	"unsafe"

	"github.com/m-oons/mike/config"
	"github.com/m-oons/mike/core"
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
	ModNoRepeat = 0x4000 // holding down the hotkey won't continuously trigger keybind
)

func Register() {
	user32 := syscall.MustLoadDLL("user32")

	reghotkey := user32.MustFindProc("RegisterHotKey")
	getmsg := user32.MustFindProc("GetMessageW")

	user32.Release()

	for i, confkey := range config.Current.Hotkeys {
		hotkey := Hotkey{
			Action: confkey.Action,
			Key:    confkey.Key,
			Ctrl:   confkey.Ctrl,
			Shift:  confkey.Shift,
			Alt:    confkey.Alt,
			Win:    confkey.Win,
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
				core.Mute()
			case "unmute":
				core.Unmute()
			case "toggle":
				core.ToggleMute()
			}
		}
	}
}
