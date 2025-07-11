package shell32

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	shell32      = syscall.MustLoadDLL("shell32")
	shellExecute = shell32.MustFindProc("ShellExecuteW")
)

func OpenURL(url string) (uintptr, error) {
	verb, _ := syscall.UTF16PtrFromString("open")
	urlp, _ := syscall.UTF16PtrFromString(url)

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
