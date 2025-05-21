package main

import (
	"github.com/m-oons/mike/config"
	"github.com/m-oons/mike/core"
	"github.com/m-oons/mike/hotkeys"
	"github.com/m-oons/mike/player"
	"github.com/m-oons/mike/tray"
)

func main() {
	config.Load()

	go tray.Create()

	core.Setup()
	defer core.Close()

	player.Setup()

	hotkeys.Register()
}
