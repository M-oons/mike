package controllers

type Controller interface {
	Init() error
	Mute() error
	Unmute() error
	ToggleMute() (bool, error)
	IsMuted() (bool, error)
	Close() error
}
