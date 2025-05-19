package tray

import (
	"os"

	"github.com/getlantern/systray"
	"github.com/m-oons/mike/assets"
	"github.com/m-oons/mike/core"
)

func Create() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTitle("Mike")
	systray.SetTooltip("Mike")

	titleItem := systray.AddMenuItem("Mike", "Mike")
	titleItem.Disable()

	systray.AddSeparator()

	muteItem := systray.AddMenuItem("Toggle Mute", "Toggle Mute")

	systray.AddSeparator()

	quitItem := systray.AddMenuItem("Quit", "Quit")

	// set icon based on initial mute state
	muted := core.IsMuted()
	if muted {
		assets.SetMuteIcon()
	} else {
		assets.SetUnmuteIcon()
	}

	// listen for menu item clicks
	for {
		select {
		case <-muteItem.ClickedCh:
			core.ToggleMute()
		case <-quitItem.ClickedCh:
			systray.Quit()
		}
	}
}

func onExit() {
	core.Close()
	os.Exit(0)
}
