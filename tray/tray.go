package tray

import (
	"os"

	"github.com/getlantern/systray"
	"github.com/m-oons/mike/actions"
	"github.com/m-oons/mike/assets"
	"github.com/m-oons/mike/devices"
)

func CreateTray() {
	systray.Run(onReady, onExit)
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

	currentMicrophone := devices.GetCurrentMicrophone()
	if currentMicrophone == nil {
		return
	}

	// set icon based on initial mute state
	if currentMicrophone.IsMuted() {
		assets.SetMuteIcon()
	} else {
		assets.SetUnmuteIcon()
	}

	// listen for menu item clicks
	for {
		select {
		case <-muteItem.ClickedCh:
			actions.ToggleMute()
		case <-quitItem.ClickedCh:
			systray.Quit()
		}
	}
}

func onExit() {
	os.Exit(0)
}
