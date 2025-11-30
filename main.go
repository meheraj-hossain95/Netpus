package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"netpus/internal/installer"
)

//go:embed all:frontend/dist
var assets embed.FS

var (
	installFlag   = flag.Bool("install", false, "Install Netpus (create shortcuts)")
	uninstallFlag = flag.Bool("uninstall", false, "Uninstall Netpus (remove shortcuts)")
	versionFlag   = flag.Bool("version", false, "Show version information")
)

const version = "1.0.0"

// Windows API for mutex
var (
	kernel32      = syscall.NewLazyDLL("kernel32.dll")
	user32        = syscall.NewLazyDLL("user32.dll")
	createMutexW  = kernel32.NewProc("CreateMutexW")
	getLastError  = kernel32.NewProc("GetLastError")
	findWindowW   = user32.NewProc("FindWindowW")
	showWindow    = user32.NewProc("ShowWindow")
	setForeground = user32.NewProc("SetForegroundWindow")
	postMessageW  = user32.NewProc("PostMessageW")
)

const (
	ERROR_ALREADY_EXISTS = 183
	SW_RESTORE           = 9
	SW_SHOW              = 5
	WM_USER              = 0x0400
	WM_SHOWWINDOW_CUSTOM = WM_USER + 100
)

// checkSingleInstance returns true if this is the first instance
func checkSingleInstance() bool {
	mutexName, _ := syscall.UTF16PtrFromString("Global\\NetpusNetworkMonitor")

	ret, _, _ := createMutexW.Call(0, 0, uintptr(unsafe.Pointer(mutexName)))
	if ret == 0 {
		return true // Failed to create mutex, assume first instance
	}

	lastErr, _, _ := getLastError.Call()
	if lastErr == ERROR_ALREADY_EXISTS {
		// Another instance exists - try to show its window
		bringExistingToFront()
		return false
	}

	return true
}

// bringExistingToFront finds and shows the existing Netpus window
func bringExistingToFront() {
	// Try to find the Netpus window by title
	windowTitle, _ := syscall.UTF16PtrFromString("Netpus Network Monitor")
	hwnd, _, _ := findWindowW.Call(0, uintptr(unsafe.Pointer(windowTitle)))

	if hwnd != 0 {
		showWindow.Call(hwnd, SW_RESTORE)
		showWindow.Call(hwnd, SW_SHOW)
		setForeground.Call(hwnd)
	}
}

func main() {
	flag.Parse()

	// Handle CLI flags
	if *versionFlag {
		fmt.Printf("Netpus Network Monitor v%s\n", version)
		fmt.Printf("Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	if *installFlag {
		execPath, err := os.Executable()
		if err != nil {
			log.Fatalf("Failed to get executable path: %v", err)
		}
		if err := installer.Install(execPath); err != nil {
			log.Fatalf("Installation failed: %v", err)
		}
		fmt.Println("Netpus installed successfully!")
		os.Exit(0)
	}

	if *uninstallFlag {
		if err := installer.Uninstall(); err != nil {
			log.Fatalf("Uninstallation failed: %v", err)
		}
		fmt.Println("Netpus uninstalled successfully!")
		os.Exit(0)
	}

	// Check for single instance using Windows mutex
	if !checkSingleInstance() {
		// Another instance is running, we tried to bring it to front
		fmt.Println("Netpus is already running. Bringing existing window to front.")
		os.Exit(0)
	}

	// Also use lock file as backup
	lockPath := filepath.Join(os.TempDir(), "netpus.lock")
	lockFile, err := os.OpenFile(lockPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err == nil {
		fmt.Fprintf(lockFile, "%d", os.Getpid())
		lockFile.Close()
		defer os.Remove(lockPath)
	}

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	runErr := wails.Run(&options.App{
		Title:  "Netpus Network Monitor",
		Width:  1440,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		OnDomReady:       app.domReady,
		OnBeforeClose:    app.beforeClose,
		OnShutdown:       app.shutdown,
		Bind: []interface{}{
			app,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
		},
	})

	if runErr != nil {
		log.Fatal("Error:", runErr.Error())
	}
}
