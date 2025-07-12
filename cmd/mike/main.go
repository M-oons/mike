package main

//go:generate go-winres make --in ./winres/winres.json --out ./resource

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/m-oons/mike/internal/audio"
	"github.com/m-oons/mike/internal/audio/controllers"
	"github.com/m-oons/mike/internal/config"
	"github.com/m-oons/mike/internal/hotkeys"
	"github.com/m-oons/mike/internal/tray"
	"github.com/m-oons/mike/internal/updater"
	"golang.org/x/sys/windows"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// handle OS signals
	go handleSignals(cancel)

	// check for updates
	go updater.CheckForUpdates()

	// load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	// setup services
	controller := newController(cfg.Controller)
	soundPlayer := audio.NewPlayer(cfg.Sounds)
	audioService := audio.NewService(controller, soundPlayer)
	if err := audioService.Setup(); err != nil {
		log.Fatalf("error setting up audio service: %v", err)
	}
	defer audioService.Close()

	// setup managers
	trayManager := tray.NewManager(audioService, cancel)
	hotkeyManager := hotkeys.NewManager(audioService, cfg.Hotkeys)

	// start managers
	go trayManager.Start(ctx)
	hotkeyManager.Start(ctx)
}

func handleSignals(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, windows.SIGINT, windows.SIGTERM)
	<-sigChan
	cancel()
}

func newController(config config.ConfigController) audio.Controller {
	var controller audio.Controller
	if config.Type == "voicemeeter" {
		controller = controllers.NewVoicemeeterController(config.Voicemeeter)
	} else {
		controller = controllers.NewWindowsController(config.Windows)
	}
	return controller
}
