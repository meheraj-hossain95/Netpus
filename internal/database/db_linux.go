//go:build linux

package database

import (
	"fmt"
	"syscall"
)

// getAvailableDiskSpace returns available disk space in bytes for Linux
func getAvailableDiskSpace(path string) (uint64, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return 0, fmt.Errorf("failed to get filesystem stats: %w", err)
	}

	// Available blocks * block size
	return stat.Bavail * uint64(stat.Bsize), nil
}
