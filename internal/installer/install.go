package installer

import (
	"fmt"
	"runtime"
)

// Install creates shortcuts and installs the application
func Install(execPath string) error {
	if runtime.GOOS == "windows" {
		return installWindows(execPath)
	} else if runtime.GOOS == "linux" {
		return installLinux(execPath)
	}
	return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
}

// Uninstall removes shortcuts and uninstalls the application
func Uninstall() error {
	if runtime.GOOS == "windows" {
		return uninstallWindows()
	} else if runtime.GOOS == "linux" {
		return uninstallLinux()
	}
	return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
}
