//go:build windows

package database

import (
	"fmt"
	"syscall"
	"unsafe"
)

// getAvailableDiskSpace returns available disk space in bytes for Windows
func getAvailableDiskSpace(path string) (uint64, error) {
	// Convert path to UTF-16
	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return 0, err
	}

	// Call GetDiskFreeSpaceEx
	var freeBytesAvailable uint64
	var totalBytes uint64
	var totalFreeBytes uint64

	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getDiskFreeSpaceEx := kernel32.NewProc("GetDiskFreeSpaceExW")

	ret, _, callErr := getDiskFreeSpaceEx.Call(
		uintptr(unsafe.Pointer(pathPtr)),
		uintptr(unsafe.Pointer(&freeBytesAvailable)),
		uintptr(unsafe.Pointer(&totalBytes)),
		uintptr(unsafe.Pointer(&totalFreeBytes)),
	)

	if ret == 0 {
		return 0, fmt.Errorf("GetDiskFreeSpaceEx failed: %v", callErr)
	}

	return freeBytesAvailable, nil
}
