package audio

import (
	"fmt"

	"github.com/m-oons/mike/internal/events"
)

type Controller interface {
	Setup() error
	Mute() error
	Unmute() error
	ToggleMute() (bool, error)
	IsMuted() (bool, error)
	Close()
}

type Player interface {
	Setup() error
	PlaySound(name string)
}

type service struct {
	controller  Controller
	player      Player
	listeners   []events.MuteStateListener
	initialized bool
}

func NewService(controller Controller, player Player) *service {
	return &service{
		controller: controller,
		player:     player,
		listeners:  make([]events.MuteStateListener, 0),
	}
}

func (s *service) AddMuteStateListener(listener events.MuteStateListener) {
	s.listeners = append(s.listeners, listener)

	if muted, err := s.controller.IsMuted(); err == nil {
		listener.OnMuteStateChanged(muted)
	}
}

func (s *service) Setup() error {
	if s.initialized {
		return nil
	}

	if err := s.controller.Setup(); err != nil {
		return fmt.Errorf("error setting up audio controller: %w", err)
	}

	if err := s.player.Setup(); err != nil {
		return fmt.Errorf("error setting up sound player: %w", err)
	}

	s.initialized = true

	return nil
}

func (s *service) Mute() error {
	if !s.initialized {
		return nil
	}

	if err := s.controller.Mute(); err != nil {
		return fmt.Errorf("error muting audio controller: %w", err)
	}

	s.notifyListeners(true)
	s.player.PlaySound("mute")

	return nil
}

func (s *service) Unmute() error {
	if !s.initialized {
		return nil
	}

	if err := s.controller.Unmute(); err != nil {
		return fmt.Errorf("error unmuting audio controller: %w", err)
	}

	s.notifyListeners(false)
	s.player.PlaySound("unmute")

	return nil
}

func (s *service) ToggleMute() (bool, error) {
	if !s.initialized {
		return false, nil
	}

	muted, err := s.controller.ToggleMute()
	if err != nil {
		return false, fmt.Errorf("error toggling mute state for audio controller: %w", err)
	}

	s.notifyListeners(muted)

	if muted {
		s.player.PlaySound("mute")
	} else {
		s.player.PlaySound("unmute")
	}

	return muted, nil
}

func (s *service) IsMuted() (bool, error) {
	if !s.initialized {
		return false, nil
	}

	muted, err := s.controller.IsMuted()
	if err != nil {
		return false, fmt.Errorf("error getting mute state for audio controller: %w", err)
	}

	return muted, nil
}

func (s *service) Close() {
	s.controller.Close()
}

func (s *service) notifyListeners(muted bool) {
	for _, listener := range s.listeners {
		listener.OnMuteStateChanged(muted)
	}
}
