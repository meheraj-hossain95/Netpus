//go:build !windows && !linux

package installer

import "fmt"

func installWindows(execPath string) error {
	return fmt.Errorf("not supported on this platform")
}

func uninstallWindows() error {
	return fmt.Errorf("not supported on this platform")
}

func installLinux(execPath string) error {
	return fmt.Errorf("not supported on this platform")
}

func uninstallLinux() error {
	return fmt.Errorf("not supported on this platform")
}
