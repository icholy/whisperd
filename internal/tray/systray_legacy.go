//go:build legacy_systray

package tray

import "github.com/getlantern/systray"

// Run starts the systray. ready is called once the tray is initialized.
// This function blocks and must be called from the main goroutine.
func Run(ready func()) {
	if !Enabled {
		if ready != nil {
			ready()
		}
		select {}
	}
	systray.Run(func() {
		SetStatus(Idle)
		if ready != nil {
			ready()
		}
	}, nil)
}

// SetStatus updates the tray icon and tooltip.
func SetStatus(s Status) {
	if !Enabled {
		return
	}
	systray.SetIcon(icons[s])
	systray.SetTooltip(tooltips[s])
}
