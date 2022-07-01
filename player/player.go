package player

import (
	"bytes"
	"math"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/m-oons/mike/assets"
)

var sounds map[string]Sound = make(map[string]Sound)

func SetupPlayer() {
	loadSound("mute", assets.MuteSound)
	loadSound("unmute", assets.UnmuteSound)
}

func PlaySound(name string, volume int) {
	sound, ok := sounds[name]
	if !ok {
		return
	}

	streamer := sound.Buffer.Streamer(0, sound.Buffer.Len())
	vol := getVolume(streamer, volume)
	speaker.Init(sound.Format.SampleRate, sound.Format.SampleRate.N(time.Second/10))
	speaker.Play(vol)
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

func getVolume(streamer beep.Streamer, volume int) *effects.Volume {
	vol := math.Log2(float64(volume) / 100)
	return &effects.Volume{
		Streamer: streamer,
		Base:     2,
		Volume:   vol,
		Silent:   false,
	}
}

func decodeWav(data []byte) (streamer beep.StreamSeekCloser, format beep.Format, err error) {
	reader := bytes.NewReader(data)
	return wav.Decode(reader)
}
