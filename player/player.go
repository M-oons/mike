package player

import (
	"bytes"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/m-oons/mike/assets"
)

var sounds map[string]Sound = make(map[string]Sound)

func SetupPlayer() {
	loadSound("mute", assets.MuteSound)
	loadSound("unmute", assets.UnmuteSound)
}

func PlaySound(name string) {
	sound, ok := sounds[name]
	if !ok {
		return
	}

	streamer := sound.Buffer.Streamer(0, sound.Buffer.Len())
	speaker.Init(sound.Format.SampleRate, sound.Format.SampleRate.N(time.Second/10))
	speaker.Play(streamer)
}

func loadSound(name string, data []byte) {
	streamer, format, err := decodeWav(data)
	if err != nil {
		return
	}

	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	streamer.Close()

	sound := Sound{
		Buffer: buffer,
		Format: format,
	}

	sounds[name] = sound
}

func decodeWav(data []byte) (streamer beep.StreamSeekCloser, format beep.Format, err error) {
	reader := bytes.NewReader(data)
	return wav.Decode(reader)
}
