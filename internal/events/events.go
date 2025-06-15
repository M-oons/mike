package events

type MuteStateListener interface {
	OnMuteStateChanged(muted bool)
}
