package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"octopus/internal/autostart"
	"octopus/internal/database"
	"octopus/internal/monitor"
	"octopus/internal/tray"
	"octopus/internal/utils"
)

// App struct
type App struct {
	ctx     context.Context
	db      *database.DB
	monitor *monitor.Monitor
	tray    *tray.Tray
	config  *utils.Config
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize database
	dbPath := utils.GetDatabasePath()
	db, err := database.New(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	a.db = db

	// Load configuration
	config, err := utils.LoadConfig(db)
	if err != nil {
		log.Printf("Failed to load config: %v, using defaults", err)
		config = utils.DefaultConfig()
	}
	a.config = config

	// Initialize monitor
	a.monitor = monitor.New(db)
	if err := a.monitor.Start(ctx); err != nil {
		log.Printf("Failed to start monitor: %v", err)
	}

	// Initialize system tray
	a.tray = tray.New(a)
	go a.tray.Setup()

	// Start background tasks
	go a.updateTrayTooltip()
	go a.periodicCleanup()
	go a.hourlyCleanup()
	go a.minuteCleanup()
}

// domReady is called after front-end resources have been loaded
func (a *App) domReady(ctx context.Context) {
	// Optional: perform actions after DOM is ready
}

// beforeClose is called when the application is about to quit
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	// Minimize to tray instead of closing
	runtime.WindowHide(ctx)
	return true // Prevent close
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	if a.monitor != nil {
		a.monitor.Stop()
	}
	if a.db != nil {
		a.db.Close()
	}
	if a.tray != nil {
		a.tray.Destroy()
	}
}

// GetNetworkStats returns current network statistics
func (a *App) GetNetworkStats() map[string]*monitor.NetworkStat {
	if a.monitor == nil {
		return make(map[string]*monitor.NetworkStat)
	}
	return a.monitor.GetStats()
}

// GetMonitorStatus returns the current monitor status
func (a *App) GetMonitorStatus() monitor.MonitorStatus {
	if a.monitor == nil {
		return monitor.MonitorStatus{}
	}
	return a.monitor.GetMonitorStatus()
}

// GetTodayStats returns today's total upload and download
func (a *App) GetTodayStats() map[string]interface{} {
	today := time.Now().Format("2006-01-02")
	summary, err := a.db.GetDailySummary(today)
	if err != nil {
		return map[string]interface{}{
			"upload":   int64(0),
			"download": int64(0),
		}
	}
	return map[string]interface{}{
		"upload":   summary.TotalUpload,
		"download": summary.TotalDownload,
	}
}

// Get24HourUsage returns the last 24 hours usage statistics
func (a *App) Get24HourUsage() map[string]interface{} {
	stats, err := a.db.Get24HourUsage()
	if err != nil {
		return map[string]interface{}{
			"upload":   int64(0),
			"download": int64(0),
		}
	}
	return map[string]interface{}{
		"upload":   stats["upload"],
		"download": stats["download"],
	}
}

// GetNetworkUsageStats returns aggregated network usage statistics
func (a *App) GetNetworkUsageStats() []database.AppUsageStat {
	days := a.config.DataRetention
	stats, err := a.db.GetAppUsageWithRetention(days)
	if err != nil {
		log.Printf("Failed to get usage stats: %v", err)
		return []database.AppUsageStat{}
	}
	return stats
}

// GetHistoricalData returns daily summaries for the specified number of days
func (a *App) GetHistoricalData(days int) []database.DailySummary {
	summaries, err := a.db.GetRecentSummaries(days)
	if err != nil {
		log.Printf("Failed to get historical data: %v", err)
		return []database.DailySummary{}
	}
	return summaries
}

// GetSettings returns current settings
func (a *App) GetSettings() utils.Config {
	return *a.config
}

// UpdateSettings updates application settings
func (a *App) UpdateSettings(settings utils.Config) error {
	// Validate settings
	if settings.Theme != "auto" && settings.Theme != "light" && settings.Theme != "dark" {
		return fmt.Errorf("invalid theme: %s", settings.Theme)
	}
	if settings.DataRetention < -1 {
		return fmt.Errorf("invalid data retention: %d", settings.DataRetention)
	}

	// Handle autostart change
	if settings.AutoStart != a.config.AutoStart {
		execPath, err := utils.GetExecutablePath()
		if err == nil {
			if settings.AutoStart {
				autostart.Enable(execPath)
			} else {
				autostart.Disable()
			}
		}
	}

	// Save settings
	if err := settings.Save(a.db); err != nil {
		return err
	}

	a.config = &settings
	return nil
}

// PauseMonitoring pauses the network monitoring
func (a *App) PauseMonitoring() {
	if a.monitor != nil {
		a.monitor.Pause()
	}
	if a.tray != nil {
		a.tray.UpdatePauseState(true)
	}
}

// ResumeMonitoring resumes the network monitoring
func (a *App) ResumeMonitoring() {
	if a.monitor != nil {
		a.monitor.Resume()
	}
	if a.tray != nil {
		a.tray.UpdatePauseState(false)
	}
}

// ClearOldData manually clears all old data from database
func (a *App) ClearOldData() error {
	// Clear all data
	if err := a.db.ClearAllData(); err != nil {
		return err
	}

	// Vacuum to reclaim space
	if err := a.db.Vacuum(); err != nil {
		return fmt.Errorf("failed to vacuum database: %w", err)
	}

	return nil
}

// GetDatabaseSize returns the database file size
func (a *App) GetDatabaseSize() int64 {
	size, err := a.db.GetSize()
	if err != nil {
		return 0
	}
	return size
}

// GetDatabaseStats returns database statistics
func (a *App) GetDatabaseStats() map[string]interface{} {
	size, _ := a.db.GetSize()
	count, _ := a.db.GetRecordCount()
	oldest, _ := a.db.GetOldestRecord()

	return map[string]interface{}{
		"size":    size,
		"records": count,
		"oldest":  oldest,
	}
}

// ShowWindow shows the application window
func (a *App) ShowWindow() {
	runtime.WindowShow(a.ctx)
	runtime.WindowUnminimise(a.ctx)
}

// HideWindow hides the application window
func (a *App) HideWindow() {
	runtime.WindowHide(a.ctx)
}

// ToggleWindow toggles window visibility
func (a *App) ToggleWindow() {
	// This is handled by the tray
}

// QuitApp completely quits the application
func (a *App) QuitApp() {
	runtime.Quit(a.ctx)
}

// updateTrayTooltip updates the tray tooltip every second
func (a *App) updateTrayTooltip() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-a.ctx.Done():
			return
		case <-ticker.C:
			if a.tray != nil && a.monitor != nil {
				stats := a.monitor.GetStats()
				var totalUp, totalDown int64
				for _, stat := range stats {
					totalUp += stat.UploadSpeed
					totalDown += stat.DownloadSpeed
				}
				tooltip := fmt.Sprintf("Octopus\n↑ %s/s ↓ %s/s",
					utils.FormatSpeed(totalUp),
					utils.FormatSpeed(totalDown))
				a.tray.UpdateTooltip(tooltip)
			}
		}
	}
}

// periodicCleanup performs daily database cleanup
func (a *App) periodicCleanup() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-a.ctx.Done():
			return
		case <-ticker.C:
			// Always clean expired records
			a.db.DeleteExpiredRecords()

			// Handle retention policy
			if a.config.DataRetention > 0 {
				cutoff := time.Now().AddDate(0, 0, -a.config.DataRetention).Unix()
				a.db.DeleteOldRecords(cutoff)
			} else if a.config.DataRetention == 0 {
				// 1-minute testing mode
				cutoff := time.Now().Add(-1 * time.Minute).Unix()
				a.db.DeleteOldRecords(cutoff)
			}
			// -1 means forever, don't delete based on age

			a.db.Vacuum()
		}
	}
}

// hourlyCleanup runs every 30 minutes for retention enforcement
func (a *App) hourlyCleanup() {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-a.ctx.Done():
			return
		case <-ticker.C:
			if a.config.DataRetention > 0 {
				cutoff := time.Now().AddDate(0, 0, -a.config.DataRetention).Unix()
				a.db.DeleteOldRecords(cutoff)
			}
		}
	}
}

// minuteCleanup runs every minute for expiration and testing retention
func (a *App) minuteCleanup() {
	// Run cleanup immediately on startup
	a.db.DeleteExpiredRecords()
	if a.config.DataRetention == 0 {
		cutoff := time.Now().Add(-1 * time.Minute).Unix()
		a.db.DeleteOldRecords(cutoff)
	} else if a.config.DataRetention > 0 {
		cutoff := time.Now().AddDate(0, 0, -a.config.DataRetention).Unix()
		a.db.DeleteOldRecords(cutoff)
	}

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-a.ctx.Done():
			return
		case <-ticker.C:
			// Delete expired records (24-hour auto-deletion)
			if err := a.db.DeleteExpiredRecords(); err != nil {
				log.Printf("Failed to delete expired records: %v", err)
			}

			// Handle 1-minute testing retention (DataRetention == 0)
			if a.config.DataRetention == 0 {
				cutoff := time.Now().Add(-1 * time.Minute).Unix()
				if err := a.db.DeleteOldRecords(cutoff); err != nil {
					log.Printf("Failed to delete 1-minute old records: %v", err)
				}
			} else if a.config.DataRetention > 0 {
				// Also enforce retention policy every minute for more immediate cleanup
				cutoff := time.Now().AddDate(0, 0, -a.config.DataRetention).Unix()
				if err := a.db.DeleteOldRecords(cutoff); err != nil {
					log.Printf("Failed to delete old records: %v", err)
				}
			}
		}
	}
}
