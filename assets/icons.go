package assets

import (
	_ "embed"

	"github.com/getlantern/systray"
)

//go:embed icons/mute.ico
var MuteIcon []byte

//go:embed icons/unmute.ico
var UnmuteIcon []byte

func SetMuteIcon() {
	systray.SetTemplateIcon(MuteIcon, MuteIcon)
}

func SetUnmuteIcon() {
	systray.SetTemplateIcon(UnmuteIcon, UnmuteIcon)
}
