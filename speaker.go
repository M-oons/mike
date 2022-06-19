package main

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

type Sound struct {
	Buffer *beep.Buffer
	Format beep.Format
}

var sounds map[string]Sound = make(map[string]Sound)

func SetupSpeaker() {
	loadSound("sounds/mute.wav")
	loadSound("sounds/unmute.wav")
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

func loadSound(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}

	streamer, format, err := decodeFile(file)
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

	name := filepath.Base(strings.TrimSuffix(filename, filepath.Ext(filename)))
	sounds[name] = sound
}

func decodeFile(file *os.File) (streamer beep.StreamSeekCloser, format beep.Format, err error) {
	filename := file.Name()
	if strings.HasSuffix(filename, ".wav") {
		return wav.Decode(file)
	}
	return mp3.Decode(file)
}
