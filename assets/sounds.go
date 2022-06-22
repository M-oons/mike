package assets

import _ "embed"

//go:embed sounds/mute.wav
var MuteSound []byte

//go:embed sounds/unmute.wav
var UnmuteSound []byte
