//go:build linux

package autostart

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	desktopEntry = `[Desktop Entry]
Type=Application
Name=Octopus Network Monitor
Comment=Real-time network bandwidth monitoring
Exec=%s
Icon=network-monitor
Terminal=false
Categories=Network;Monitor;
X-GNOME-Autostart-enabled=true
`
)

// Enable enables autostart on Linux via XDG autostart
func Enable(execPath string) error {
	autostartDir := filepath.Join(os.Getenv("HOME"), ".config", "autostart")
	if err := os.MkdirAll(autostartDir, 0755); err != nil {
		return fmt.Errorf("failed to create autostart directory: %w", err)
	}

	desktopFile := filepath.Join(autostartDir, "octopus.desktop")
	content := fmt.Sprintf(desktopEntry, execPath)

	err := os.WriteFile(desktopFile, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write desktop entry: %w", err)
	}

	return nil
}

// Disable disables autostart on Linux
func Disable() error {
	desktopFile := filepath.Join(os.Getenv("HOME"), ".config", "autostart", "octopus.desktop")

	err := os.Remove(desktopFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove desktop entry: %w", err)
	}

	return nil
}

// IsEnabled checks if autostart is enabled
func IsEnabled() (bool, error) {
	desktopFile := filepath.Join(os.Getenv("HOME"), ".config", "autostart", "octopus.desktop")
	_, err := os.Stat(desktopFile)

	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}
