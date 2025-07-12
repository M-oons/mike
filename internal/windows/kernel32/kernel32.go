package kernel32

import "golang.org/x/sys/windows"

var (
	kernel32           = windows.NewLazySystemDLL("kernel32")
	getCurrentThreadId = kernel32.NewProc("GetCurrentThreadId")
)

func GetCurrentThreadID() (uintptr, error) {
	threadID, _, err := getCurrentThreadId.Call()
	return threadID, err
}
