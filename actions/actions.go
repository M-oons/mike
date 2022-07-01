package actions

import (
	"github.com/m-oons/mike/assets"
	"github.com/m-oons/mike/config"
	"github.com/m-oons/mike/devices"
	"github.com/m-oons/mike/player"
)

func Mute() {
	mic := devices.GetCurrentMicrophone()
	if mic == nil {
		return
	}

	if !mic.Mute() {
		return
	}

	tryPlaySound("mute")
	assets.SetMuteIcon()
}

func Unmute() {
	mic := devices.GetCurrentMicrophone()
	if mic == nil {
		return
	}

	if !mic.Unmute() {
		return
	}

	tryPlaySound("unmute")
	assets.SetUnmuteIcon()
}

func ToggleMute() {
	mic := devices.GetCurrentMicrophone()
	if mic == nil {
		return
	}

	muted := mic.ToggleMute()
	if muted {
		tryPlaySound("mute")
		assets.SetMuteIcon()
	} else {
		tryPlaySound("unmute")
		assets.SetUnmuteIcon()
	}
}

func tryPlaySound(sound string) {
	if !config.Current.Sounds.Enabled {
		return
	}

	player.PlaySound(sound, config.Current.Sounds.Volume)
}
