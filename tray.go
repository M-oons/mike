package main

import (
	"os"

	"github.com/getlantern/systray"
	"github.com/m-oons/mike/assets"
)

func CreateTray() {
	systray.Run(onReady, onExit)
}

func SetMuteIcon() {
	systray.SetTemplateIcon(assets.MuteIcon, assets.MuteIcon)
}

func SetUnmuteIcon() {
	systray.SetTemplateIcon(assets.UnmuteIcon, assets.UnmuteIcon)
}

func onReady() {
	systray.SetTitle("Mike")
	systray.SetTooltip("Mike")

	titleItem := systray.AddMenuItem("Mike", "Mike")
	titleItem.Disable()

	systray.AddSeparator()

	// micMenu := systray.AddMenuItem("Microphones", "Microphones")

	// populate microphones list
	// microphones := GetMicrophones()
	// for _, mic := range microphones {
	// 	micItem := micMenu.AddSubMenuItem(mic.Name(), mic.Name())
	// 	if mic.Id() == currentMicrophone.Id() {
	// 		micItem.Check()
	// 	}
	// }

	muteItem := systray.AddMenuItem("Toggle Mute", "Toggle Mute")

	systray.AddSeparator()

	quitItem := systray.AddMenuItem("Quit", "Quit")

	currentMicrophone := GetCurrentMicrophone()
	if currentMicrophone == nil {
		return
	}

	// set icon based on initial mute state
	if currentMicrophone.IsMuted() {
		SetMuteIcon()
	} else {
		SetUnmuteIcon()
	}

	// listen for menu item clicks
	for {
		select {
		case <-muteItem.ClickedCh:
			ToggleMute()
		case <-quitItem.ClickedCh:
			systray.Quit()
		}
	}
}

func onExit() {
	os.Exit(0)
}
