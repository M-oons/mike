package main

import (
	"strings"
	"syscall"
	"time"
	"unsafe"
)

type Hotkey struct {
	Modifiers int
	KeyCode   int
}

type MSG struct {
	HWND   uintptr
	UINT   uintptr
	WPARAM int16
	LPARAM int64
	DWORD  int32
	POINT  struct{ X, Y int64 }
}

const (
	ModAlt = 1 << iota
	ModCtrl
	ModShift
	ModWin
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

	hotkeys := []*Hotkey{
		getHotkey("O", true, false, true, false),
	}

	for i, hotkey := range hotkeys {
		reghotkey.Call(0, uintptr(i+1), uintptr(hotkey.Modifiers), uintptr(hotkey.KeyCode))
	}

	for {
		msg := &MSG{}
		getmsg.Call(uintptr(unsafe.Pointer(msg)), 0, 0, 0)

		switch msg.WPARAM {
		case 1:
			mic := GetCurrentMicrophone()
			mic.ToggleMute()
		}

		time.Sleep(time.Millisecond * 50)
	}
}

func getHotkey(key string, ctrl bool, shift bool, alt bool, windows bool) *Hotkey {
	if keycode, ok := keys[strings.ToUpper(key)]; ok {
		modifiers := getModifiers(ctrl, shift, alt, windows)
		return &Hotkey{Modifiers: modifiers, KeyCode: keycode}
	}
	return nil
}

func getModifiers(ctrl bool, shift bool, alt bool, windows bool) int {
	modifiers := 0
	if ctrl {
		modifiers += ModCtrl
	}
	if shift {
		modifiers += ModShift
	}
	if alt {
		modifiers += ModAlt
	}
	if windows {
		modifiers += ModWin
	}
	return modifiers
}
