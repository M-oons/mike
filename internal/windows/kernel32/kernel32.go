package kernel32

import "syscall"

var (
	kernel32           = syscall.MustLoadDLL("kernel32")
	getCurrentThreadId = kernel32.MustFindProc("GetCurrentThreadId")
)

func GetCurrentThreadID() (uintptr, error) {
	threadID, _, err := getCurrentThreadId.Call()
	return threadID, err
}

func Close() {
	kernel32.Release()
}
