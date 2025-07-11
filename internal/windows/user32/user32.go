package user32

import (
	"syscall"
	"unsafe"
)

const (
	MB_YESNO        = 0x00000004
	MB_ICONQUESTION = 0x00000020
	IDYES           = 6
)

var (
	user32            = syscall.MustLoadDLL("user32")
	registerHotKey    = user32.MustFindProc("RegisterHotKey")
	unregisterHotKey  = user32.MustFindProc("UnregisterHotKey")
	getMessage        = user32.MustFindProc("GetMessageW")
	postThreadMessage = user32.MustFindProc("PostThreadMessageW")
	messageBox        = user32.MustFindProc("MessageBoxW")
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

func RegisterHotKey(hotkeyID uintptr, keyCode int, modifiers int) (uintptr, error) {
	ret, _, err := registerHotKey.Call(
		0,
		hotkeyID,
		uintptr(modifiers),
		uintptr(keyCode),
	)
	return ret, err
}

func UnregisterHotKey(hotkeyID uintptr) (uintptr, error) {
	ret, _, err := unregisterHotKey.Call(
		0,
		hotkeyID,
	)
	return ret, err
}

func GetMessage() (MSG, uintptr, error) {
	var msg MSG
	ret, _, err := getMessage.Call(
		uintptr(unsafe.Pointer(&msg)),
		0,
		0,
		0,
	)
	return msg, ret, err
}

func PostThreadMessage(threadID uintptr, message uintptr) (uintptr, error) {
	ret, _, err := postThreadMessage.Call(
		threadID,
		message,
		0,
		0,
	)
	return ret, err
}

func ShowMessageBox(title string, message string) bool {
	t, _ := syscall.UTF16PtrFromString(title)
	m, _ := syscall.UTF16PtrFromString(message)
	ret, _, _ := messageBox.Call(
		0,
		uintptr(unsafe.Pointer(m)),
		uintptr(unsafe.Pointer(t)),
		MB_YESNO|MB_ICONQUESTION,
	)
	return ret == IDYES
}

func Close() {
	user32.Release()
}
