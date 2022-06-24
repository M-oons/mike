package main

import (
	"github.com/m-oons/mike/config"
	"github.com/m-oons/mike/devices"
	"github.com/m-oons/mike/hotkeys"
	"github.com/m-oons/mike/player"
	"github.com/m-oons/mike/tray"
)

func main() {
	config.LoadConfig()
	devices.InitializeCOM()
	go player.SetupPlayer()
	go tray.CreateTray()
	hotkeys.RegisterHotkeys()
}
