//go:build legacy_systray

package tray

import "github.com/getlantern/systray"

// Run starts the systray. ready is called once the tray is initialized,
// and exit is called when the tray is shutting down.
// This function blocks and must be called from the main goroutine.
func Run(ready func(), exit func()) {
	systray.Run(func() {
		SetStatus(Idle)
		if ready != nil {
			ready()
		}
	}, exit)
}

// SetStatus updates the tray icon and tooltip.
func SetStatus(s Status) {
	systray.SetIcon(icons[s])
	systray.SetTooltip(tooltips[s])
}
