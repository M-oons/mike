package audio

import (
	"bytes"
	"fmt"
	"math"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/m-oons/mike/assets"
	"github.com/m-oons/mike/config"
)

type sound struct {
	Buffer *beep.Buffer
	Format beep.Format
}

type player struct {
	config      config.ConfigSounds
	sounds      map[string]*sound
	initialized bool
}

func NewPlayer(config config.ConfigSounds) *player {
	return &player{
		config: config,
		sounds: make(map[string]*sound),
	}
}

func (p *player) Setup() error {
	if p.initialized {
		return nil
	}

	if err := p.loadSound("mute", assets.MuteSound); err != nil {
		return fmt.Errorf("error loading mute sound: %w", err)
	}
	if err := p.loadSound("unmute", assets.UnmuteSound); err != nil {
		return fmt.Errorf("error loading unmute sound: %w", err)
	}

	sound := p.sounds["mute"] // assume all sounds have the same sample rate
	if err := speaker.Init(sound.Format.SampleRate, sound.Format.SampleRate.N(time.Second/10)); err != nil {
		return fmt.Errorf("error initializing speaker: %w", err)
	}
	p.initialized = true

	return nil
}

func (p *player) PlaySound(name string) {
	if !p.config.Enabled || !p.initialized {
		return
	}

	sound, ok := p.sounds[name]
	if !ok {
		return
	}

	streamer := sound.Buffer.Streamer(0, sound.Buffer.Len())
	vol := effects.Volume{
		Streamer: streamer,
		Base:     2,
		Volume:   math.Log2(float64(p.config.Volume) / 100),
		Silent:   false,
	}
	speaker.Play(&vol)
}

func (p *player) loadSound(name string, data []byte) error {
	streamer, format, err := decodeWav(data)
	if err != nil {
		return err
	}

	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	if err := streamer.Close(); err != nil {
		return err
	}

	p.sounds[name] = &sound{
		Buffer: buffer,
		Format: format,
	}

	return nil
}

func decodeWav(data []byte) (streamer beep.StreamSeekCloser, format beep.Format, err error) {
	reader := bytes.NewReader(data)
	return wav.Decode(reader)
}
