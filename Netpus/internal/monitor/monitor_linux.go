//go:build linux

package monitor

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type processData struct {
	processID     int
	uploadBytes   int64
	downloadBytes int64
}

// getNetworkProcesses collects network statistics for all processes on Linux
func getNetworkProcesses() (map[string]processData, error) {
	// Parse /proc/net files for connections
	tcpInodes, err := parseNetFile("/proc/net/tcp")
	if err != nil {
		return nil, fmt.Errorf("failed to parse TCP: %w", err)
	}

	tcp6Inodes, err := parseNetFile("/proc/net/tcp6")
	if err != nil {
		return nil, fmt.Errorf("failed to parse TCP6: %w", err)
	}

	udpInodes, err := parseNetFile("/proc/net/udp")
	if err != nil {
		return nil, fmt.Errorf("failed to parse UDP: %w", err)
	}

	udp6Inodes, err := parseNetFile("/proc/net/udp6")
	if err != nil {
		return nil, fmt.Errorf("failed to parse UDP6: %w", err)
	}

	// Combine all inodes
	allInodes := make(map[string]bool)
	for inode := range tcpInodes {
		allInodes[inode] = true
	}
	for inode := range tcp6Inodes {
		allInodes[inode] = true
	}
	for inode := range udpInodes {
		allInodes[inode] = true
	}
	for inode := range udp6Inodes {
		allInodes[inode] = true
	}

	// Map inodes to processes
	inodeToProc := make(map[string]int)
	procDirs, err := filepath.Glob("/proc/[0-9]*")
	if err != nil {
		return nil, err
	}

	for _, procDir := range procDirs {
		pidStr := filepath.Base(procDir)
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			continue
		}

		fdDir := filepath.Join(procDir, "fd")
		fds, err := os.ReadDir(fdDir)
		if err != nil {
			continue
		}

		for _, fd := range fds {
			link, err := os.Readlink(filepath.Join(fdDir, fd.Name()))
			if err != nil {
				continue
			}

			if strings.HasPrefix(link, "socket:[") {
				inode := strings.TrimPrefix(link, "socket:[")
				inode = strings.TrimSuffix(inode, "]")

				if allInodes[inode] {
					inodeToProc[inode] = pid
				}
			}
		}
	}

	// Get network statistics for each process
	result := make(map[string]processData)
	processedPIDs := make(map[int]bool)

	for _, pid := range inodeToProc {
		if processedPIDs[pid] {
			continue
		}
		processedPIDs[pid] = true

		// Get process name
		exePath, err := os.Readlink(fmt.Sprintf("/proc/%d/exe", pid))
		if err != nil {
			continue
		}
		appName := filepath.Base(exePath)

		// Get network stats from /proc/[pid]/net/dev
		upload, download, err := getProcessNetStats(pid)
		if err != nil {
			continue
		}

		if existing, exists := result[appName]; exists {
			existing.uploadBytes += upload
			existing.downloadBytes += download
			result[appName] = existing
		} else {
			result[appName] = processData{
				processID:     pid,
				uploadBytes:   upload,
				downloadBytes: download,
			}
		}
	}

	return result, nil
}

// parseNetFile parses /proc/net/* files and returns socket inodes
func parseNetFile(path string) (map[string]bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	inodes := make(map[string]bool)
	scanner := bufio.NewScanner(file)

	// Skip header
	scanner.Scan()

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 10 {
			continue
		}

		inode := fields[9]
		if inode != "0" {
			inodes[inode] = true
		}
	}

	return inodes, scanner.Err()
}

// getProcessNetStats retrieves network statistics for a specific process
func getProcessNetStats(pid int) (upload, download int64, err error) {
	devPath := fmt.Sprintf("/proc/%d/net/dev", pid)
	file, err := os.Open(devPath)
	if err != nil {
		// Fallback to system-wide stats
		return getSystemNetStats()
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Skip first two header lines
	scanner.Scan()
	scanner.Scan()

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}

		// Skip loopback interface
		if strings.HasPrefix(fields[0], "lo:") {
			continue
		}

		// Receive bytes (download)
		if rx, err := strconv.ParseInt(fields[1], 10, 64); err == nil {
			download += rx
		}

		// Transmit bytes (upload)
		if tx, err := strconv.ParseInt(fields[9], 10, 64); err == nil {
			upload += tx
		}
	}

	return upload, download, scanner.Err()
}

// getSystemNetStats gets system-wide network statistics
func getSystemNetStats() (upload, download int64, err error) {
	file, err := os.Open("/proc/net/dev")
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Skip first two header lines
	scanner.Scan()
	scanner.Scan()

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}

		// Skip loopback interface
		if strings.HasPrefix(fields[0], "lo:") {
			continue
		}

		// Receive bytes (download)
		if rx, err := strconv.ParseInt(fields[1], 10, 64); err == nil {
			download += rx
		}

		// Transmit bytes (upload)
		if tx, err := strconv.ParseInt(fields[9], 10, 64); err == nil {
			upload += tx
		}
	}

	return upload, download, scanner.Err()
}
