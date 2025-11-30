package database

import (
	"database/sql"
	"fmt"
	"os"
	"runtime"
	"time"

	_ "modernc.org/sqlite"
)

// DB represents the database connection
type DB struct {
	conn *sql.DB
	path string
}

// UsageRecord represents a network usage record
type UsageRecord struct {
	ID            int64
	AppName       string
	ProcessID     int
	UploadBytes   int64
	DownloadBytes int64
	Timestamp     int64
	IsTemporary   bool
	ExpiresAt     int64
}

// DailySummary represents daily aggregated statistics
type DailySummary struct {
	ID            int64
	Date          string
	TotalUpload   int64
	TotalDownload int64
}

// AppMetadata represents application metadata
type AppMetadata struct {
	AppName        string
	ExecutablePath string
	FirstSeen      int64
	LastSeen       int64
}

// AppUsageStat represents aggregated app usage statistics
type AppUsageStat struct {
	AppName       string
	TotalUpload   int64
	TotalDownload int64
	LastSeen      int64
}

// New creates a new database connection
func New(dbPath string) (*DB, error) {
	// Ensure directory exists
	dir := dbPath[:len(dbPath)-len("/netpus.db")]
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Check for existing database corruption before opening
	if _, err := os.Stat(dbPath); err == nil {
		// Database exists, try to check integrity
		if err := checkDatabaseIntegrity(dbPath); err != nil {
			// Database is corrupted, backup and recreate
			backupPath := dbPath + ".corrupted." + time.Now().Format("20060102_150405")
			fmt.Printf("⚠ Database corruption detected! Backing up to: %s\n", backupPath)
			if copyErr := os.Rename(dbPath, backupPath); copyErr != nil {
				fmt.Printf("Warning: Could not backup corrupted database: %v\n", copyErr)
			}
			// Remove WAL and SHM files if they exist
			os.Remove(dbPath + "-wal")
			os.Remove(dbPath + "-shm")
		}
	}

	// Open database
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	conn.SetMaxOpenConns(5) // Increased for better concurrency
	conn.SetMaxIdleConns(3)
	conn.SetConnMaxLifetime(time.Hour)

	// Enable WAL mode and other pragmas
	pragmas := []string{
		"PRAGMA journal_mode=WAL",
		"PRAGMA synchronous=NORMAL",
		"PRAGMA busy_timeout=10000", // Increased to 10 seconds for better concurrency
		"PRAGMA foreign_keys=ON",
		"PRAGMA cache_size=-64000", // 64MB cache
	}

	for _, pragma := range pragmas {
		if _, err := conn.Exec(pragma); err != nil {
			conn.Close()
			return nil, fmt.Errorf("failed to set pragma: %w", err)
		}
	}

	db := &DB{
		conn: conn,
		path: dbPath,
	}

	// Initialize schema
	if err := db.initSchema(); err != nil {
		conn.Close()
		return nil, err
	}

	return db, nil
}

// checkDatabaseIntegrity checks if database file is corrupted
func checkDatabaseIntegrity(dbPath string) error {
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("cannot open database: %w", err)
	}
	defer conn.Close()

	// Run integrity check
	var result string
	err = conn.QueryRow("PRAGMA integrity_check").Scan(&result)
	if err != nil {
		return fmt.Errorf("integrity check failed: %w", err)
	}

	if result != "ok" {
		return fmt.Errorf("database integrity check failed: %s", result)
	}

	return nil
}

// initSchema creates the database schema
func (db *DB) initSchema() error {
	// Create basic schema without optional columns
	schema := `
	CREATE TABLE IF NOT EXISTS usage_records (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		app_name TEXT NOT NULL,
		process_id INTEGER,
		upload_bytes INTEGER NOT NULL,
		download_bytes INTEGER NOT NULL,
		timestamp INTEGER NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_usage_timestamp ON usage_records(timestamp);
	CREATE INDEX IF NOT EXISTS idx_usage_app ON usage_records(app_name);

	CREATE TABLE IF NOT EXISTS daily_summaries (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT UNIQUE NOT NULL,
		total_upload INTEGER NOT NULL,
		total_download INTEGER NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_daily_date ON daily_summaries(date);

	CREATE TABLE IF NOT EXISTS app_metadata (
		app_name TEXT PRIMARY KEY,
		executable_path TEXT,
		first_seen INTEGER NOT NULL,
		last_seen INTEGER NOT NULL
	);

	CREATE TABLE IF NOT EXISTS settings (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL
	);
	`

	_, err := db.conn.Exec(schema)
	if err != nil {
		return err
	}

	// Perform schema migration to add missing columns
	if err := db.migrateSchema(); err != nil {
		return err
	}

	// Create indexes for new columns after migration
	indexSchema := `
	CREATE INDEX IF NOT EXISTS idx_usage_expires ON usage_records(expires_at);
	CREATE INDEX IF NOT EXISTS idx_usage_temporary ON usage_records(is_temporary);
	`
	_, err = db.conn.Exec(indexSchema)
	return err
}

// migrateSchema handles database migrations
func (db *DB) migrateSchema() error {
	// Get list of existing columns in usage_records table
	rows, err := db.conn.Query("PRAGMA table_info(usage_records)")
	if err != nil {
		return fmt.Errorf("failed to get table info: %w", err)
	}
	defer rows.Close()

	existingColumns := make(map[string]bool)
	for rows.Next() {
		var cid int
		var name, colType string
		var notNull, pk int
		var dfltValue sql.NullString
		if err := rows.Scan(&cid, &name, &colType, &notNull, &dfltValue, &pk); err != nil {
			return err
		}
		existingColumns[name] = true
	}

	// Add expires_at column if it doesn't exist
	if !existingColumns["expires_at"] {
		_, err := db.conn.Exec("ALTER TABLE usage_records ADD COLUMN expires_at INTEGER")
		if err != nil {
			return fmt.Errorf("failed to add expires_at column: %w", err)
		}
		fmt.Println("✓ Database migrated: added expires_at column")
	}

	// Add is_temporary column if it doesn't exist
	if !existingColumns["is_temporary"] {
		_, err := db.conn.Exec("ALTER TABLE usage_records ADD COLUMN is_temporary INTEGER DEFAULT 0")
		if err != nil {
			return fmt.Errorf("failed to add is_temporary column: %w", err)
		}
		fmt.Println("✓ Database migrated: added is_temporary column")
	}

	return nil
}

// InsertUsageRecord inserts a single usage record with retry logic
func (db *DB) InsertUsageRecord(record UsageRecord) error {
	query := `INSERT INTO usage_records (app_name, process_id, upload_bytes, download_bytes, timestamp, is_temporary, expires_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?)`

	isTemp := 0
	if record.IsTemporary {
		isTemp = 1
	}

	// Retry with exponential backoff for database lock errors
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		_, err := db.conn.Exec(query, record.AppName, record.ProcessID,
			record.UploadBytes, record.DownloadBytes, record.Timestamp, isTemp, record.ExpiresAt)
		if err == nil {
			return nil
		}

		// Check if it's a lock error
		errStr := err.Error()
		if i < maxRetries-1 && (contains(errStr, "database is locked") || contains(errStr, "SQLITE_BUSY")) {
			// Exponential backoff: 50ms, 100ms, 200ms, 400ms
			backoff := time.Duration(50*(1<<uint(i))) * time.Millisecond
			time.Sleep(backoff)
			continue
		}
		return err
	}
	return fmt.Errorf("failed to insert usage record after retries")
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			len(s) > len(substr) && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// BatchInsertUsageRecords inserts multiple usage records in a transaction
func (db *DB) BatchInsertUsageRecords(records []UsageRecord) error {
	if len(records) == 0 {
		return nil
	}

	// Check available disk space before large writes
	if len(records) > 100 {
		if err := db.checkDiskSpace(); err != nil {
			return fmt.Errorf("insufficient disk space: %w", err)
		}
	}

	tx, err := db.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`INSERT INTO usage_records
		(app_name, process_id, upload_bytes, download_bytes, timestamp, is_temporary, expires_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, record := range records {
		isTemp := 0
		if record.IsTemporary {
			isTemp = 1
		}
		_, err := stmt.Exec(record.AppName, record.ProcessID,
			record.UploadBytes, record.DownloadBytes, record.Timestamp, isTemp, record.ExpiresAt)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// UpdateDailySummary updates or inserts daily summary
func (db *DB) UpdateDailySummary(date string, upload, download int64) error {
	query := `INSERT INTO daily_summaries (date, total_upload, total_download)
	          VALUES (?, ?, ?)
	          ON CONFLICT(date) DO UPDATE SET
	          total_upload = total_upload + excluded.total_upload,
	          total_download = total_download + excluded.total_download`

	_, err := db.conn.Exec(query, date, upload, download)
	return err
}

// UpsertAppMetadata updates or inserts app metadata
func (db *DB) UpsertAppMetadata(metadata AppMetadata) error {
	query := `INSERT INTO app_metadata (app_name, executable_path, first_seen, last_seen)
	          VALUES (?, ?, ?, ?)
	          ON CONFLICT(app_name) DO UPDATE SET
	          executable_path = excluded.executable_path,
	          last_seen = excluded.last_seen`

	_, err := db.conn.Exec(query, metadata.AppName, metadata.ExecutablePath,
		metadata.FirstSeen, metadata.LastSeen)
	return err
}

// GetDailySummary retrieves a daily summary for a specific date
func (db *DB) GetDailySummary(date string) (*DailySummary, error) {
	query := `SELECT id, date, total_upload, total_download
	          FROM daily_summaries WHERE date = ?`

	var summary DailySummary
	err := db.conn.QueryRow(query, date).Scan(
		&summary.ID, &summary.Date, &summary.TotalUpload, &summary.TotalDownload)
	if err == sql.ErrNoRows {
		return &DailySummary{Date: date, TotalUpload: 0, TotalDownload: 0}, nil
	}
	if err != nil {
		return nil, err
	}
	return &summary, nil
}

// GetRecentSummaries retrieves the last N days of summaries
func (db *DB) GetRecentSummaries(days int) ([]DailySummary, error) {
	query := `SELECT id, date, total_upload, total_download
	          FROM daily_summaries
	          ORDER BY date DESC
	          LIMIT ?`

	rows, err := db.conn.Query(query, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []DailySummary
	for rows.Next() {
		var s DailySummary
		if err := rows.Scan(&s.ID, &s.Date, &s.TotalUpload, &s.TotalDownload); err != nil {
			return nil, err
		}
		summaries = append(summaries, s)
	}
	return summaries, rows.Err()
}

// GetAppMetadata retrieves metadata for a specific app
func (db *DB) GetAppMetadata(appName string) (*AppMetadata, error) {
	query := `SELECT app_name, executable_path, first_seen, last_seen
	          FROM app_metadata WHERE app_name = ?`

	var meta AppMetadata
	err := db.conn.QueryRow(query, appName).Scan(
		&meta.AppName, &meta.ExecutablePath, &meta.FirstSeen, &meta.LastSeen)
	if err != nil {
		return nil, err
	}
	return &meta, nil
}

// GetAppUsageStats retrieves aggregated usage statistics for all apps
func (db *DB) GetAppUsageStats(startTime, endTime int64) ([]AppUsageStat, error) {
	query := `SELECT app_name,
	          SUM(upload_bytes) as total_upload,
	          SUM(download_bytes) as total_download,
	          MAX(timestamp) as last_seen
	          FROM usage_records
	          WHERE timestamp BETWEEN ? AND ?
	          GROUP BY app_name
	          ORDER BY (total_upload + total_download) DESC`

	rows, err := db.conn.Query(query, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []AppUsageStat
	for rows.Next() {
		var s AppUsageStat
		if err := rows.Scan(&s.AppName, &s.TotalUpload, &s.TotalDownload, &s.LastSeen); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, rows.Err()
}

// GetAppUsageWithRetention retrieves app usage stats based on retention period
func (db *DB) GetAppUsageWithRetention(days int) ([]AppUsageStat, error) {
	var startTime int64
	if days == 0 {
		startTime = 0 // All time
	} else {
		startTime = time.Now().AddDate(0, 0, -days).Unix()
	}
	endTime := time.Now().Unix()
	return db.GetAppUsageStats(startTime, endTime)
}

// Get24HourUsage retrieves total usage for the last 24 hours
func (db *DB) Get24HourUsage() (map[string]int64, error) {
	startTime := time.Now().Add(-24 * time.Hour).Unix()
	query := `SELECT SUM(upload_bytes), SUM(download_bytes)
	          FROM usage_records
	          WHERE timestamp >= ?`

	var upload, download sql.NullInt64
	err := db.conn.QueryRow(query, startTime).Scan(&upload, &download)
	if err != nil {
		return nil, err
	}

	return map[string]int64{
		"upload":   upload.Int64,
		"download": download.Int64,
	}, nil
}

// GetUsageByTimeRange retrieves records within a specific time range
func (db *DB) GetUsageByTimeRange(startTime, endTime int64) ([]UsageRecord, error) {
	query := `SELECT id, app_name, process_id, upload_bytes, download_bytes, timestamp, is_temporary, COALESCE(expires_at, 0)
	          FROM usage_records
	          WHERE timestamp BETWEEN ? AND ?
	          ORDER BY timestamp DESC`

	rows, err := db.conn.Query(query, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []UsageRecord
	for rows.Next() {
		var r UsageRecord
		var isTemp int
		if err := rows.Scan(&r.ID, &r.AppName, &r.ProcessID,
			&r.UploadBytes, &r.DownloadBytes, &r.Timestamp, &isTemp, &r.ExpiresAt); err != nil {
			return nil, err
		}
		r.IsTemporary = isTemp == 1
		records = append(records, r)
	}
	return records, rows.Err()
}

// DeleteOldRecords deletes records older than the specified timestamp
func (db *DB) DeleteOldRecords(beforeTimestamp int64) error {
	// Delete old usage records
	query := `DELETE FROM usage_records WHERE timestamp < ? AND is_temporary = 0`
	_, err := db.conn.Exec(query, beforeTimestamp)
	if err != nil {
		return fmt.Errorf("failed to delete old usage records: %w", err)
	}

	// Also delete old daily summaries based on the same cutoff date
	cutoffDate := time.Unix(beforeTimestamp, 0).Format("2006-01-02")
	summaryQuery := `DELETE FROM daily_summaries WHERE date < ?`
	_, err = db.conn.Exec(summaryQuery, cutoffDate)
	if err != nil {
		return fmt.Errorf("failed to delete old daily summaries: %w", err)
	}

	return nil
}

// ClearAllData clears all usage records and daily summaries from the database
func (db *DB) ClearAllData() error {
	// Clear all usage records
	if _, err := db.conn.Exec("DELETE FROM usage_records"); err != nil {
		return fmt.Errorf("failed to clear usage records: %w", err)
	}

	// Clear all daily summaries
	if _, err := db.conn.Exec("DELETE FROM daily_summaries"); err != nil {
		return fmt.Errorf("failed to clear daily summaries: %w", err)
	}

	return nil
}

// DeleteExpiredRecords deletes records that have passed their expiration time
func (db *DB) DeleteExpiredRecords() error {
	now := time.Now().Unix()
	query := `DELETE FROM usage_records WHERE expires_at > 0 AND expires_at < ?`
	result, err := db.conn.Exec(query, now)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		fmt.Printf("Deleted %d expired records\n", rowsAffected)
	}

	return nil
}

// DeleteTemporaryRecords deletes all temporary records (24-hour data)
func (db *DB) DeleteTemporaryRecords() error {
	query := `DELETE FROM usage_records WHERE is_temporary = 1`
	_, err := db.conn.Exec(query)
	return err
}

// Vacuum performs database vacuum to reclaim space
func (db *DB) Vacuum() error {
	_, err := db.conn.Exec("VACUUM")
	return err
}

// GetSize returns the database file size in bytes
func (db *DB) GetSize() (int64, error) {
	info, err := os.Stat(db.path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// GetRecordCount returns the total number of usage records
func (db *DB) GetRecordCount() (int64, error) {
	var count int64
	err := db.conn.QueryRow("SELECT COUNT(*) FROM usage_records").Scan(&count)
	return count, err
}

// GetOldestRecord returns the timestamp of the oldest record
func (db *DB) GetOldestRecord() (int64, error) {
	var timestamp sql.NullInt64
	err := db.conn.QueryRow("SELECT MIN(timestamp) FROM usage_records").Scan(&timestamp)
	if err != nil || !timestamp.Valid {
		return 0, err
	}
	return timestamp.Int64, nil
}

// GetSetting retrieves a setting value
func (db *DB) GetSetting(key string) (string, error) {
	var value string
	err := db.conn.QueryRow("SELECT value FROM settings WHERE key = ?", key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return value, err
}

// SetSetting stores a setting value
func (db *DB) SetSetting(key, value string) error {
	query := `INSERT INTO settings (key, value) VALUES (?, ?)
	          ON CONFLICT(key) DO UPDATE SET value = excluded.value`
	_, err := db.conn.Exec(query, key, value)
	return err
}

// GetAllSettings retrieves all settings as a map
func (db *DB) GetAllSettings() (map[string]string, error) {
	rows, err := db.conn.Query("SELECT key, value FROM settings")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	settings := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		settings[key] = value
	}
	return settings, rows.Err()
}

// Close closes the database connection
func (db *DB) Close() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}

// checkDiskSpace checks if there's enough disk space for database operations
func (db *DB) checkDiskSpace() error {
	// Get database directory
	dir := db.path[:len(db.path)-len("/netpus.db")]

	// Check if database file size is reasonable (< 10GB as a safety limit)
	if info, err := os.Stat(db.path); err == nil {
		if info.Size() > 10*1024*1024*1024 {
			return fmt.Errorf("database file exceeds 10GB safety limit")
		}
	}

	// Check actual available disk space on Windows
	if runtime.GOOS == "windows" {
		freeBytesAvailable, err := getAvailableDiskSpace(dir)
		if err != nil {
			// If we can't check disk space, log warning but don't fail
			fmt.Printf("Warning: Could not check disk space: %v\n", err)
			return nil
		}

		// Require at least 100MB free space
		const minFreeSpace = 100 * 1024 * 1024
		if freeBytesAvailable < minFreeSpace {
			return fmt.Errorf("insufficient disk space: %d bytes available, need at least %d bytes",
				freeBytesAvailable, minFreeSpace)
		}
	}

	return nil
}
