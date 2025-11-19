//go:build darwin

package monitor

import (
	"fmt"
)

type processData struct {
	processID     int
	uploadBytes   int64
	downloadBytes int64
}

// getNetworkProcesses collects network statistics for all processes on macOS (stub)
func getNetworkProcesses() (map[string]processData, error) {
	// macOS implementation not yet available
	// This stub prevents build errors but returns no data
	return make(map[string]processData), fmt.Errorf("macOS network monitoring not yet implemented")
}
