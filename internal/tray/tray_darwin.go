//go:build darwin

package tray

// Tray represents the system tray icon (stub for macOS)
type Tray struct {
	app interface{}
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

// Setup initializes the system tray (stub)
func (t *Tray) Setup() {
	// macOS implementation not yet available
}

// UpdateTooltip updates the tray icon tooltip (stub)
func (t *Tray) UpdateTooltip(text string) {
	// macOS implementation not yet available
}

// UpdatePauseState updates the pause/resume menu items (stub)
func (t *Tray) UpdatePauseState(paused bool) {
	// macOS implementation not yet available
}

// Destroy cleans up the tray icon (stub)
func (t *Tray) Destroy() {
	// macOS implementation not yet available
}

// ShowNotification shows a system notification (stub)
func (t *Tray) ShowNotification(title, message string) {
	// macOS implementation not yet available
}
