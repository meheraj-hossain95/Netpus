package monitor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"netpus/internal/database"
)

const (
	UPDATE_INTERVAL   = 500 * time.Millisecond // Fast 500ms collection for responsive real-time UI
	BATCH_INTERVAL    = 10 * time.Second       // Database write interval (zero data loss)
	CLEANUP_THRESHOLD = 3 * time.Second        // Inactive process cleanup time (3-second timeout)
)

// NetworkStat represents network statistics for a single application
type NetworkStat struct {
	AppName       string
	ProcessID     int
	UploadSpeed   int64 // Bytes per second
	DownloadSpeed int64 // Bytes per second
	TotalUpload   int64 // Total bytes uploaded
	TotalDownload int64 // Total bytes downloaded
	LastUpdate    time.Time
}

// MonitorStatus represents the current monitor state
type MonitorStatus struct {
	Running        bool      `json:"running"`
	Paused         bool      `json:"paused"`
	UpdateInterval int       `json:"updateInterval"` // Seconds
	LastUpdate     time.Time `json:"lastUpdate"`
}

// Monitor represents the network monitoring system
type Monitor struct {
	ctx         context.Context
	cancel      context.CancelFunc
	db          interface{}
	stats       map[string]*NetworkStat
	statsMux    sync.RWMutex
	batch       []batchRecord
	batchMux    sync.Mutex
	paused      bool
	pauseMux    sync.RWMutex
	lastUpdate  time.Time
	saveEnabled bool
	saveMux     sync.RWMutex
}

type batchRecord struct {
	appName     string
	processID   int
	upload      int64
	download    int64
	timestamp   int64
	isTemporary bool
	expiresAt   int64
}

// New creates a new Monitor instance
func New(db interface{}) *Monitor {
	return &Monitor{
		db:          db,
		stats:       make(map[string]*NetworkStat),
		batch:       make([]batchRecord, 0),
		saveEnabled: true,
	}
}

// Start begins network monitoring
func (m *Monitor) Start(ctx context.Context) error {
	m.ctx, m.cancel = context.WithCancel(ctx)

	// Initialize monitoring with 2 collections
	fmt.Println("Initializing network monitor...")
	if err := m.collect(); err != nil {
		return fmt.Errorf("initial collection failed: %w", err)
	}
	time.Sleep(UPDATE_INTERVAL)
	if err := m.collect(); err != nil {
		return fmt.Errorf("second collection failed: %w", err)
	}

	// Start monitoring loops
	go m.monitorLoop()
	go m.batchWriteLoop()

	fmt.Println("Network monitor started successfully")
	return nil
}

// Stop stops network monitoring and flushes remaining data
func (m *Monitor) Stop() {
	if m.cancel != nil {
		m.cancel()
	}
	m.flushBatch()
}

// monitorLoop is the main monitoring loop
func (m *Monitor) monitorLoop() {
	ticker := time.NewTicker(UPDATE_INTERVAL)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			m.pauseMux.RLock()
			paused := m.paused
			m.pauseMux.RUnlock()

			if !paused {
				if err := m.collect(); err != nil {
					fmt.Printf("Collection error: %v\n", err)
				}
				m.cleanupInactive()
			}
		}
	}
}

// batchWriteLoop handles periodic database writes
func (m *Monitor) batchWriteLoop() {
	ticker := time.NewTicker(BATCH_INTERVAL)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			m.flushBatch()
		}
	}
}

// collect gathers network statistics (platform-specific implementation)
func (m *Monitor) collect() error {
	m.pauseMux.RLock()
	if m.paused {
		m.pauseMux.RUnlock()
		return nil
	}
	m.pauseMux.RUnlock()

	// Get network data from platform-specific implementation
	// NOTE: getNetworkProcesses() now returns DELTA bytes (bytes transferred since last call)
	// distributed proportionally to processes with active connections
	processes, err := getNetworkProcesses()
	if err != nil {
		return err
	}

	now := time.Now()
	m.statsMux.Lock()
	defer m.statsMux.Unlock()

	// Update stats with delta values directly
	for appName, data := range processes {
		// data.uploadBytes and data.downloadBytes are already deltas
		uploadDelta := data.uploadBytes
		downloadDelta := data.downloadBytes

		// Skip if no activity
		if uploadDelta == 0 && downloadDelta == 0 {
			continue
		}

		// Get or create stat entry
		stat, exists := m.stats[appName]
		if !exists {
			stat = &NetworkStat{
				AppName:   appName,
				ProcessID: data.processID,
			}
			m.stats[appName] = stat
		}

		// Calculate time delta for speed calculation
		timeDelta := 1.0 // Default 1 second
		if !stat.LastUpdate.IsZero() {
			timeDelta = now.Sub(stat.LastUpdate).Seconds()
			if timeDelta <= 0 {
				timeDelta = 1.0
			}
		}

		// Calculate speeds (bytes per second)
		stat.UploadSpeed = int64(float64(uploadDelta) / timeDelta)
		stat.DownloadSpeed = int64(float64(downloadDelta) / timeDelta)
		stat.TotalUpload += uploadDelta
		stat.TotalDownload += downloadDelta
		stat.LastUpdate = now

		// Add to batch for database storage
		m.batchMux.Lock()
		expiresAt := now.Add(24 * time.Hour).Unix()
		m.batch = append(m.batch, batchRecord{
			appName:     appName,
			processID:   data.processID,
			upload:      uploadDelta,
			download:    downloadDelta,
			timestamp:   now.Unix(),
			isTemporary: false,
			expiresAt:   expiresAt,
		})
		m.batchMux.Unlock()
	}

	// Reset speeds for apps that didn't have activity this cycle
	for appName, stat := range m.stats {
		if _, hasActivity := processes[appName]; !hasActivity {
			stat.UploadSpeed = 0
			stat.DownloadSpeed = 0
		}
	}

	m.lastUpdate = now
	return nil
}

// flushBatch writes batched records to database
func (m *Monitor) flushBatch() {
	// Check if saving is enabled
	m.saveMux.RLock()
	if !m.saveEnabled {
		m.saveMux.RUnlock()
		// Clear the batch without saving
		m.batchMux.Lock()
		m.batch = make([]batchRecord, 0)
		m.batchMux.Unlock()
		return
	}
	m.saveMux.RUnlock()

	m.batchMux.Lock()
	if len(m.batch) == 0 {
		m.batchMux.Unlock()
		return
	}

	batch := m.batch
	m.batch = make([]batchRecord, 0)
	m.batchMux.Unlock()

	// Convert batch records to database records
	var totalUpload, totalDownload int64
	records := make([]database.UsageRecord, len(batch))
	appMetadataMap := make(map[string]database.AppMetadata)
	now := time.Now().Unix()

	for i, rec := range batch {
		records[i] = database.UsageRecord{
			AppName:       rec.appName,
			ProcessID:     rec.processID,
			UploadBytes:   rec.upload,
			DownloadBytes: rec.download,
			Timestamp:     rec.timestamp,
			IsTemporary:   rec.isTemporary,
			ExpiresAt:     rec.expiresAt,
		}
		totalUpload += rec.upload
		totalDownload += rec.download

		// Track app metadata (deduplicate by app name)
		if _, exists := appMetadataMap[rec.appName]; !exists {
			appMetadataMap[rec.appName] = database.AppMetadata{
				AppName:        rec.appName,
				ExecutablePath: rec.appName, // Could be enhanced with full path
				FirstSeen:      rec.timestamp,
				LastSeen:       now,
			}
		}
	}

	// Write to database with proper error handling
	if db, ok := m.db.(*database.DB); ok {
		if err := db.BatchInsertUsageRecords(records); err != nil {
			fmt.Printf("Failed to insert batch records: %v\n", err)
			return
		}

		// Update app metadata
		for _, metadata := range appMetadataMap {
			if err := db.UpsertAppMetadata(metadata); err != nil {
				fmt.Printf("Failed to update app metadata for %s: %v\n", metadata.AppName, err)
			}
		}

		// Update daily summary
		today := time.Now().Format("2006-01-02")
		if err := db.UpdateDailySummary(today, totalUpload, totalDownload); err != nil {
			fmt.Printf("Failed to update daily summary: %v\n", err)
		}

		fmt.Printf("Successfully flushed %d records (%s up, %s down)\n",
			len(batch), formatBytes(totalUpload), formatBytes(totalDownload))
	}
}

// formatBytes is a helper function for logging
func formatBytes(bytes int64) string {
	if bytes == 0 {
		return "0 B"
	}
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// cleanupInactive removes inactive processes from tracking
func (m *Monitor) cleanupInactive() {
	now := time.Now()
	m.statsMux.Lock()
	defer m.statsMux.Unlock()

	for appName, stat := range m.stats {
		if now.Sub(stat.LastUpdate) > CLEANUP_THRESHOLD {
			delete(m.stats, appName)
		}
	}
}

// GetStats returns a copy of current statistics
func (m *Monitor) GetStats() map[string]*NetworkStat {
	m.statsMux.RLock()
	defer m.statsMux.RUnlock()

	stats := make(map[string]*NetworkStat)
	for k, v := range m.stats {
		statCopy := *v
		stats[k] = &statCopy
	}
	return stats
}

// GetMonitorStatus returns current monitor status
func (m *Monitor) GetMonitorStatus() MonitorStatus {
	m.pauseMux.RLock()
	paused := m.paused
	m.pauseMux.RUnlock()

	return MonitorStatus{
		Running:        m.ctx != nil,
		Paused:         paused,
		UpdateInterval: int(UPDATE_INTERVAL.Seconds()),
		LastUpdate:     m.lastUpdate,
	}
}

// Pause pauses network monitoring
func (m *Monitor) Pause() {
	m.pauseMux.Lock()
	m.paused = true
	m.pauseMux.Unlock()
}

// Resume resumes network monitoring
func (m *Monitor) Resume() {
	m.pauseMux.Lock()
	m.paused = false
	m.pauseMux.Unlock()
}

// SetSaveEnabled enables or disables saving data to database
func (m *Monitor) SetSaveEnabled(enabled bool) {
	m.saveMux.Lock()
	m.saveEnabled = enabled
	m.saveMux.Unlock()
	if enabled {
		fmt.Println("Data saving enabled")
	} else {
		fmt.Println("Data saving disabled - data will not be stored")
	}
}
