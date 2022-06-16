package main

import (
	"github.com/go-ole/go-ole"
)

func main() {
	if err := ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED); err != nil {
		return
	}

	RegisterHotkeys()
	CreateTray()
}
