//go:build windows

package autostart

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

const (
	regPath = `Software\Microsoft\Windows\CurrentVersion\Run`
	appName = "Netpus"
)

// Enable enables autostart on Windows via registry
func Enable(execPath string) error {
	key, err := registry.OpenKey(registry.CURRENT_USER, regPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	// Quote the path to handle spaces
	quotedPath := fmt.Sprintf(`"%s"`, execPath)

	err = key.SetStringValue(appName, quotedPath)
	if err != nil {
		return fmt.Errorf("failed to set registry value: %w", err)
	}

	return nil
}

// Disable disables autostart on Windows
func Disable() error {
	key, err := registry.OpenKey(registry.CURRENT_USER, regPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	err = key.DeleteValue(appName)
	if err != nil && err != registry.ErrNotExist {
		return fmt.Errorf("failed to delete registry value: %w", err)
	}

	return nil
}

// IsEnabled checks if autostart is enabled
func IsEnabled() (bool, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, regPath, registry.QUERY_VALUE)
	if err != nil {
		return false, nil
	}
	defer key.Close()

	_, _, err = key.GetStringValue(appName)
	if err == registry.ErrNotExist {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}
