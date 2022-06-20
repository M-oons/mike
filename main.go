package main

const (
	AppName   = "mike"
	AppAuthor = "Moons"
	AppUrl    = "https://github.com/m-oons/mike"
)

func main() {
	LoadConfig()
	InitializeCOM()
	go SetupSpeaker()
	go CreateTray()
	RegisterHotkeys()
}

func Mute() {
	mic := GetCurrentMicrophone()
	if mic == nil {
		return
	}

	if !mic.Mute() {
		return
	}

	PlaySound("mute")
	SetMuteIcon()
}

func Unmute() {
	mic := GetCurrentMicrophone()
	if mic == nil {
		return
	}

	if !mic.Unmute() {
		return
	}

	PlaySound("unmute")
	SetUnmuteIcon()
}

func ToggleMute() {
	mic := GetCurrentMicrophone()
	if mic == nil {
		return
	}

	muted := mic.ToggleMute()
	if muted {
		PlaySound("mute")
		SetMuteIcon()
	} else {
		PlaySound("unmute")
		SetUnmuteIcon()
	}
}
