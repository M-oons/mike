package core

import (
	"github.com/m-oons/mike/assets"
	"github.com/m-oons/mike/config"
	"github.com/m-oons/mike/controllers"
	"github.com/m-oons/mike/player"
)

var controller controllers.Controller

func Setup() {
	switch config.Current.Controller.Type {
	case "windows":
		controller = &controllers.WindowsController{}
	case "voicemeeter":
		controller = &controllers.VoicemeeterController{}
	default:
		controller = &controllers.WindowsController{}
	}

	controller.Init()
}

func Mute() {
	if controller.Mute() != nil {
		return
	}

	player.PlaySound("mute")
	assets.SetMuteIcon()
}

func Unmute() {
	if controller.Unmute() != nil {
		return
	}

	player.PlaySound("unmute")
	assets.SetUnmuteIcon()
}

func ToggleMute() {
	muted, err := controller.ToggleMute()
	if err != nil {
		return
	}

	if muted {
		player.PlaySound("mute")
		assets.SetMuteIcon()
	} else {
		player.PlaySound("unmute")
		assets.SetUnmuteIcon()
	}
}

func IsMuted() bool {
	muted, err := controller.IsMuted()
	if err != nil {
		return false
	}

	return muted
}

func Close() {
	controller.Close()
}
