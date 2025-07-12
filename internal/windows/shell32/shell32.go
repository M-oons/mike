package shell32

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	shell32      = windows.NewLazySystemDLL("shell32")
	shellExecute = shell32.NewProc("ShellExecuteW")
)

func OpenURL(url string) (uintptr, error) {
	verb, _ := windows.UTF16PtrFromString("open")
	urlp, _ := windows.UTF16PtrFromString(url)

	ret, _, err := shellExecute.Call(
		0,
		uintptr(unsafe.Pointer(verb)),
		uintptr(unsafe.Pointer(urlp)),
		0,
		0,
		uintptr(windows.SW_SHOWNORMAL),
	)
	if ret <= 32 {
		return ret, err
	}
	return ret, nil
}
