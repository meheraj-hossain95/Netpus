package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

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

	// Check for multiple instances using lock file
	lockPath := filepath.Join(os.TempDir(), "netpus.lock")
	lockFile, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		if os.IsExist(err) {
			// Lock file exists, check if process is actually running
			if data, readErr := os.ReadFile(lockPath); readErr == nil {
				fmt.Printf("Netpus is already running (PID: %s)\n", string(data))
				fmt.Println("Only one instance of Netpus can run at a time.")
				os.Exit(1)
			}
		}
		log.Printf("Warning: Could not create lock file: %v", err)
	} else {
		// Write our PID to lock file
		fmt.Fprintf(lockFile, "%d", os.Getpid())
		lockFile.Close()

		// Clean up lock file on exit
		defer func() {
			os.Remove(lockPath)
		}()
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
