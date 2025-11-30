//go:build windows

package installer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"golang.org/x/sys/windows/registry"
)

const (
	uninstallKey = `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\Netpus`
	appVersion   = "1.0.0"
)

// installWindows creates desktop and start menu shortcuts on Windows
func installWindows(execPath string) error {
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	// Create desktop shortcut
	desktop := os.Getenv("USERPROFILE") + "\\Desktop"
	if err := createShortcut(execPath, filepath.Join(desktop, "Netpus.lnk")); err != nil {
		return fmt.Errorf("failed to create desktop shortcut: %w", err)
	}

	// Create start menu shortcut
	startMenu := os.Getenv("APPDATA") + "\\Microsoft\\Windows\\Start Menu\\Programs"
	if err := createShortcut(execPath, filepath.Join(startMenu, "Netpus.lnk")); err != nil {
		return fmt.Errorf("failed to create start menu shortcut: %w", err)
	}

	// Add to Windows "Apps & Features" (Add/Remove Programs)
	if err := addUninstallEntry(execPath); err != nil {
		// Non-fatal, just log
		fmt.Printf("Warning: Failed to add uninstall entry: %v\n", err)
	}

	return nil
}

// addUninstallEntry adds Netpus to Windows "Apps & Features"
func addUninstallEntry(execPath string) error {
	key, _, err := registry.CreateKey(registry.CURRENT_USER, uninstallKey, registry.ALL_ACCESS)
	if err != nil {
		return fmt.Errorf("failed to create registry key: %w", err)
	}
	defer key.Close()

	installDir := filepath.Dir(execPath)
	uninstallCmd := fmt.Sprintf(`"%s" --uninstall`, execPath)

	// Set required values for Apps & Features
	if err := key.SetStringValue("DisplayName", "Netpus Network Monitor"); err != nil {
		return err
	}
	if err := key.SetStringValue("DisplayVersion", appVersion); err != nil {
		return err
	}
	if err := key.SetStringValue("Publisher", "Netpus"); err != nil {
		return err
	}
	if err := key.SetStringValue("InstallLocation", installDir); err != nil {
		return err
	}
	if err := key.SetStringValue("UninstallString", uninstallCmd); err != nil {
		return err
	}
	if err := key.SetStringValue("DisplayIcon", execPath); err != nil {
		return err
	}
	if err := key.SetDWordValue("NoModify", 1); err != nil {
		return err
	}
	if err := key.SetDWordValue("NoRepair", 1); err != nil {
		return err
	}

	// Estimate size in KB
	if fi, err := os.Stat(execPath); err == nil {
		sizeKB := uint32(fi.Size() / 1024)
		key.SetDWordValue("EstimatedSize", sizeKB)
	}

	return nil
}

// uninstallWindows removes shortcuts, registry entry, and optionally app data on Windows
func uninstallWindows() error {
	desktop := os.Getenv("USERPROFILE") + "\\Desktop\\Netpus.lnk"
	startMenu := os.Getenv("APPDATA") + "\\Microsoft\\Windows\\Start Menu\\Programs\\Netpus.lnk"

	os.Remove(desktop)
	os.Remove(startMenu)

	// Remove from Apps & Features
	registry.DeleteKey(registry.CURRENT_USER, uninstallKey)

	// Remove autostart registry entry
	autorunKey, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`, registry.ALL_ACCESS)
	if err == nil {
		autorunKey.DeleteValue("Netpus")
		autorunKey.Close()
	}

	// Remove app data folder
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData != "" {
		appDataDir := filepath.Join(localAppData, "Netpus")
		os.RemoveAll(appDataDir)
	}

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
	oleutil.PutProperty(idispatch, "Description", "Netpus Network Monitor")
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
