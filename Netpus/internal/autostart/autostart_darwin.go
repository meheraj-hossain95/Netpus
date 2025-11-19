//go:build darwin

package autostart

// Enable adds the application to system startup (stub for macOS)
func Enable(execPath string) error {
	// macOS implementation not yet available
	return nil
}

// Disable removes the application from system startup (stub for macOS)
func Disable() error {
	// macOS implementation not yet available
	return nil
}

// IsEnabled checks if autostart is enabled (stub for macOS)
func IsEnabled() bool {
	// macOS implementation not yet available
	return false
}
