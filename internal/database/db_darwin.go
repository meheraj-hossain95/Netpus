//go:build darwin

package database

// getAvailableDiskSpace returns available disk space in bytes for macOS
func getAvailableDiskSpace(path string) (uint64, error) {
	// For macOS, we could use syscall.Statfs but for now return a safe default
	// This allows compilation on macOS for testing purposes
	return 1024 * 1024 * 1024, nil // Return 1GB as available
}
