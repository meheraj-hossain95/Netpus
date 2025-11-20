//go:build linux

package tray

import (
	"github.com/getlantern/systray"
)

// Tray represents the system tray icon
type Tray struct {
	app        interface{}
	menuShow   *systray.MenuItem
	menuHide   *systray.MenuItem
	menuPause  *systray.MenuItem
	menuResume *systray.MenuItem
	menuQuit   *systray.MenuItem
}

// AppInterface defines the required methods from the main app
type AppInterface interface {
	ShowWindow()
	HideWindow()
	PauseMonitoring()
	ResumeMonitoring()
	QuitApp()
}

// New creates a new Tray instance
func New(app interface{}) *Tray {
	return &Tray{
		app: app,
	}
}

// Setup initializes the system tray
func (t *Tray) Setup() {
	systray.Run(t.onReady, t.onExit)
}

// onReady is called when systray is ready
func (t *Tray) onReady() {
	// Set icon and tooltip
	// Note: SetTooltip may not be available on all Linux systems
	// systray.SetTooltip("Octopus Network Monitor")

	// Create menu items
	t.menuShow = systray.AddMenuItem("Show Dashboard", "Show the main window")
	t.menuHide = systray.AddMenuItem("Hide Dashboard", "Hide the main window")
	systray.AddSeparator()
	t.menuPause = systray.AddMenuItem("Pause Monitoring", "Pause network monitoring")
	t.menuResume = systray.AddMenuItem("Resume Monitoring", "Resume network monitoring")
	t.menuResume.Hide() // Initially hidden
	systray.AddSeparator()
	t.menuQuit = systray.AddMenuItem("Quit", "Quit the application")

	// Handle menu clicks
	go t.handleMenuClicks()
}

// onExit is called when systray exits
func (t *Tray) onExit() {
	// Cleanup
}

// handleMenuClicks processes menu item clicks
func (t *Tray) handleMenuClicks() {
	app, ok := t.app.(AppInterface)
	if !ok {
		return
	}

	for {
		select {
		case <-t.menuShow.ClickedCh:
			app.ShowWindow()

		case <-t.menuHide.ClickedCh:
			app.HideWindow()

		case <-t.menuPause.ClickedCh:
			app.PauseMonitoring()
			t.menuPause.Hide()
			t.menuResume.Show()

		case <-t.menuResume.ClickedCh:
			app.ResumeMonitoring()
			t.menuResume.Hide()
			t.menuPause.Show()

		case <-t.menuQuit.ClickedCh:
			app.QuitApp()
			systray.Quit()
			return
		}
	}
}

// UpdateTooltip updates the tray icon tooltip
func (t *Tray) UpdateTooltip(text string) {
	// Note: SetTooltip may not be available on all Linux systems
	// Comment out if build fails
	// systray.SetTooltip(text)
}

// UpdatePauseState updates the pause/resume menu items
func (t *Tray) UpdatePauseState(paused bool) {
	if paused {
		t.menuPause.Hide()
		t.menuResume.Show()
	} else {
		t.menuResume.Hide()
		t.menuPause.Show()
	}
}

// Destroy cleans up the tray icon
func (t *Tray) Destroy() {
	systray.Quit()
}

// ShowNotification shows a system notification (reserved for future use)
func (t *Tray) ShowNotification(title, message string) {
	// TODO: Implement notifications
}
