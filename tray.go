package main

import (
	"github.com/getlantern/systray"
)

func CreateTray() {
	systray.Run(onReady, onExit)
}

func onReady() {
	muteIcon := GetIconData("icons/mute.ico")
	unmuteIcon := GetIconData("icons/unmute.ico")

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

	// set icon based on initial mute state
	if currentMicrophone.IsMuted() {
		systray.SetTemplateIcon(muteIcon, muteIcon)
	} else {
		systray.SetTemplateIcon(unmuteIcon, unmuteIcon)
	}

	// listen for menu item clicks
	go func() {
		for {
			select {
			case <-muteItem.ClickedCh:
				if currentMicrophone.ToggleMute() {
					systray.SetTemplateIcon(muteIcon, muteIcon)
				} else {
					systray.SetTemplateIcon(unmuteIcon, unmuteIcon)
				}
			case <-quitItem.ClickedCh:
				systray.Quit()
			}
		}
	}()
}

func onExit() {

}
