package tray

import (
	"context"

	"github.com/getlantern/systray"
	"github.com/m-oons/mike/internal/assets"
	"github.com/m-oons/mike/internal/events"
)

type Service interface {
	ToggleMute() (bool, error)
	AddMuteStateListener(listener events.MuteStateListener)
}

type manager struct {
	service Service
	cancel  context.CancelFunc
}

func NewManager(service Service, cancel context.CancelFunc) *manager {
	return &manager{
		service: service,
		cancel:  cancel,
	}
}

func (m *manager) Start(ctx context.Context) {
	systray.Run(m.onReady(ctx), m.onExit)
}

func (m *manager) OnMuteStateChanged(muted bool) {
	if muted {
		m.setMuteIcon()
	} else {
		m.setUnmuteIcon()
	}
}

func (m *manager) setMuteIcon() {
	systray.SetTemplateIcon(assets.MuteIcon, assets.MuteIcon)
}

func (m *manager) setUnmuteIcon() {
	systray.SetTemplateIcon(assets.UnmuteIcon, assets.UnmuteIcon)
}

func (m *manager) onReady(ctx context.Context) func() {
	return func() {
		systray.SetTitle("Mike")
		systray.SetTooltip("Mike")
		systray.AddMenuItem("Mike", "Mike").Disable()
		systray.AddSeparator()
		muteItem := systray.AddMenuItem("Toggle Mute", "Toggle Mute")
		systray.AddSeparator()
		quitItem := systray.AddMenuItem("Quit", "Quit")

		m.service.AddMuteStateListener(m)

		// listen for menu item clicks
		for {
			select {
			case <-muteItem.ClickedCh:
				m.service.ToggleMute()

			case <-quitItem.ClickedCh:
				m.cancel()
				return

			case <-ctx.Done():
				systray.Quit()
				return
			}
		}
	}
}

func (m *manager) onExit() {}
