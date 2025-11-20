//go:build windows

package installer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

// installWindows creates desktop and start menu shortcuts on Windows
func installWindows(execPath string) error {
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	// Create desktop shortcut
	desktop := os.Getenv("USERPROFILE") + "\\Desktop"
	if err := createShortcut(execPath, filepath.Join(desktop, "Octopus.lnk")); err != nil {
		return fmt.Errorf("failed to create desktop shortcut: %w", err)
	}

	// Create start menu shortcut
	startMenu := os.Getenv("APPDATA") + "\\Microsoft\\Windows\\Start Menu\\Programs"
	if err := createShortcut(execPath, filepath.Join(startMenu, "Octopus.lnk")); err != nil {
		return fmt.Errorf("failed to create start menu shortcut: %w", err)
	}

	return nil
}

// uninstallWindows removes shortcuts on Windows
func uninstallWindows() error {
	desktop := os.Getenv("USERPROFILE") + "\\Desktop\\Octopus.lnk"
	startMenu := os.Getenv("APPDATA") + "\\Microsoft\\Windows\\Start Menu\\Programs\\Octopus.lnk"

	os.Remove(desktop)
	os.Remove(startMenu)

	return nil
}

// createShortcut creates a Windows shortcut file
func createShortcut(targetPath, shortcutPath string) error {
	oleShellObject, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		return err
	}
	defer oleShellObject.Release()

	wshell, err := oleShellObject.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return err
	}
	defer wshell.Release()

	cs, err := oleutil.CallMethod(wshell, "CreateShortcut", shortcutPath)
	if err != nil {
		return err
	}
	idispatch := cs.ToIDispatch()
	defer idispatch.Release()

	oleutil.PutProperty(idispatch, "TargetPath", targetPath)
	oleutil.PutProperty(idispatch, "Description", "Octopus Network Monitor")
	oleutil.PutProperty(idispatch, "WorkingDirectory", filepath.Dir(targetPath))

	// Set icon if it exists
	iconPath := filepath.Join(filepath.Dir(targetPath), "icon.ico")
	if _, err := os.Stat(iconPath); err == nil {
		oleutil.PutProperty(idispatch, "IconLocation", iconPath)
	}

	_, err = oleutil.CallMethod(idispatch, "Save")
	return err
}

// Stubs for Linux functions (not used on Windows)
func installLinux(execPath string) error {
	return fmt.Errorf("not supported on Windows")
}

func uninstallLinux() error {
	return fmt.Errorf("not supported on Windows")
}
