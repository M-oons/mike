package actions

import (
	"github.com/m-oons/mike/assets"
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

	player.PlaySound("mute")
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

	player.PlaySound("unmute")
	assets.SetUnmuteIcon()
}

func ToggleMute() {
	mic := devices.GetCurrentMicrophone()
	if mic == nil {
		return
	}

	muted := mic.ToggleMute()
	if muted {
		player.PlaySound("mute")
		assets.SetMuteIcon()
	} else {
		player.PlaySound("unmute")
		assets.SetUnmuteIcon()
	}
}
