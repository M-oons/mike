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

	if config.Current.Sounds {
		player.PlaySound("mute")
	}
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

	if config.Current.Sounds {
		player.PlaySound("unmute")
	}
	assets.SetUnmuteIcon()
}

func ToggleMute() {
	mic := devices.GetCurrentMicrophone()
	if mic == nil {
		return
	}

	muted := mic.ToggleMute()
	sounds := config.Current.Sounds
	if muted {
		if sounds {
			player.PlaySound("mute")
		}
		assets.SetMuteIcon()
	} else {
		if sounds {
			player.PlaySound("unmute")
		}
		assets.SetUnmuteIcon()
	}
}
