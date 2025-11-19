//go:build linux

package installer

import (
	"fmt"
	"os"
	"path/filepath"
)

const desktopEntry = `[Desktop Entry]
Type=Application
Name=Octopus Network Monitor
Comment=Real-time network bandwidth monitoring application
Exec=%s
Icon=network-monitor
Terminal=false
Categories=Network;Monitor;System;
Keywords=network;bandwidth;monitor;traffic;
`

// installLinux creates an application menu entry on Linux
func installLinux(execPath string) error {
	appDir := filepath.Join(os.Getenv("HOME"), ".local", "share", "applications")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return fmt.Errorf("failed to create applications directory: %w", err)
	}

	desktopFile := filepath.Join(appDir, "octopus.desktop")
	content := fmt.Sprintf(desktopEntry, execPath)

	err := os.WriteFile(desktopFile, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write desktop entry: %w", err)
	}

	// Make executable
	if err := os.Chmod(desktopFile, 0755); err != nil {
		return fmt.Errorf("failed to make desktop entry executable: %w", err)
	}

	return nil
}

// uninstallLinux removes the application menu entry on Linux
func uninstallLinux() error {
	desktopFile := filepath.Join(os.Getenv("HOME"), ".local", "share", "applications", "octopus.desktop")

	err := os.Remove(desktopFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove desktop entry: %w", err)
	}

	return nil
}

// Stubs for Windows functions (not used on Linux)
func installWindows(execPath string) error {
	return fmt.Errorf("not supported on Linux")
}

func uninstallWindows() error {
	return fmt.Errorf("not supported on Linux")
}
