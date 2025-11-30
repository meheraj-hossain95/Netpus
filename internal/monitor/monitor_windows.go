//go:build windows

package monitor

import (
	"fmt"
	"path/filepath"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	iphlpapi                = syscall.NewLazyDLL("iphlpapi.dll")
	procGetExtendedTcpTable = iphlpapi.NewProc("GetExtendedTcpTable")
	procGetExtendedUdpTable = iphlpapi.NewProc("GetExtendedUdpTable")
	procGetIfTable2         = iphlpapi.NewProc("GetIfTable2")
)

type processData struct {
	processID     int
	uploadBytes   int64
	downloadBytes int64
}

// getNetworkProcesses collects network statistics for all processes on Windows
func getNetworkProcesses() (map[string]processData, error) {
	// Get system-wide network I/O
	totalUpload, totalDownload, err := getSystemNetworkIO()
	if err != nil {
		return nil, fmt.Errorf("failed to get system network I/O: %w", err)
	}

	// Get connection information
	tcpConns, err := getTCPStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get TCP stats: %w", err)
	}

	udpConns, err := getUDPStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get UDP stats: %w", err)
	}

	// Build process connection map with weights
	processWeights := make(map[uint32]float64)
	processNames := make(map[uint32]string)
	var totalWeight float64

	// TCP connections (higher weight)
	for _, conn := range tcpConns {
		weight := 1.0
		if conn.State == 5 { // ESTABLISHED
			weight = 10.0
		}
		processWeights[conn.OwningPid] += weight
		totalWeight += weight

		if _, exists := processNames[conn.OwningPid]; !exists {
			if name := getProcessPath(conn.OwningPid); name != "" {
				processNames[conn.OwningPid] = name
			}
		}
	}

	// UDP connections (lower weight)
	for _, conn := range udpConns {
		weight := 0.5
		processWeights[conn.OwningPid] += weight
		totalWeight += weight

		if _, exists := processNames[conn.OwningPid]; !exists {
			if name := getProcessPath(conn.OwningPid); name != "" {
				processNames[conn.OwningPid] = name
			}
		}
	}

	// Distribute network bytes based on weights
	result := make(map[string]processData)

	if totalWeight > 0 {
		for pid, weight := range processWeights {
			name, exists := processNames[pid]
			if !exists || name == "" {
				continue
			}

			proportion := weight / totalWeight
			upload := int64(float64(totalUpload) * proportion)
			download := int64(float64(totalDownload) * proportion)

			result[name] = processData{
				processID:     int(pid),
				uploadBytes:   upload,
				downloadBytes: download,
			}
		}
	}

	return result, nil
}

// MIB_IF_ROW2 structure (simplified)
type mibIfRow2 struct {
	InterfaceLuid               uint64
	InterfaceIndex              uint32
	InterfaceGuid               [16]byte
	Alias                       [257]uint16
	Description                 [257]uint16
	PhysicalAddressLength       uint32
	PhysicalAddress             [32]byte
	PermanentPhysicalAddress    [32]byte
	Mtu                         uint32
	Type                        uint32
	TunnelType                  uint32
	MediaType                   uint32
	PhysicalMediumType          uint32
	AccessType                  uint32
	DirectionType               uint32
	InterfaceAndOperStatusFlags byte
	OperStatus                  uint32
	AdminStatus                 uint32
	MediaConnectState           uint32
	NetworkGuid                 [16]byte
	ConnectionType              uint32
	TransmitLinkSpeed           uint64
	ReceiveLinkSpeed            uint64
	InOctets                    uint64
	InUcastPkts                 uint64
	InNUcastPkts                uint64
	InDiscards                  uint64
	InErrors                    uint64
	InUnknownProtos             uint64
	InUcastOctets               uint64
	InMulticastOctets           uint64
	InBroadcastOctets           uint64
	OutOctets                   uint64
	OutUcastPkts                uint64
	OutNUcastPkts               uint64
	OutDiscards                 uint64
	OutErrors                   uint64
	OutUcastOctets              uint64
	OutMulticastOctets          uint64
	OutBroadcastOctets          uint64
	OutQLen                     uint64
}

type mibIfTable2 struct {
	NumEntries uint32
	Table      [1]mibIfRow2
}

// getSystemNetworkIO gets total network I/O from all interfaces
func getSystemNetworkIO() (upload, download int64, err error) {
	var table *mibIfTable2
	ret, _, _ := procGetIfTable2.Call(uintptr(unsafe.Pointer(&table)))
	if ret != 0 {
		return 0, 0, fmt.Errorf("GetIfTable2 failed with code %d", ret)
	}
	if table == nil {
		return 0, 0, nil
	}

	numEntries := int(table.NumEntries)
	if numEntries > 0 {
		entries := unsafe.Slice(&table.Table[0], numEntries)
		for _, entry := range entries {
			upload += int64(entry.OutOctets)
			download += int64(entry.InOctets)
		}
	}

	return upload, download, nil
}

type tcpRow struct {
	State      uint32
	LocalAddr  uint32
	LocalPort  uint32
	RemoteAddr uint32
	RemotePort uint32
	OwningPid  uint32
}

type tcpTable struct {
	NumEntries uint32
	Table      [1]tcpRow
}

// getTCPStats retrieves TCP connection information
func getTCPStats() ([]tcpRow, error) {
	var size uint32
	family := uint32(windows.AF_INET)
	class := uint32(5) // TCP_TABLE_OWNER_PID_ALL

	// Get required size
	procGetExtendedTcpTable.Call(
		0,
		uintptr(unsafe.Pointer(&size)),
		0,
		uintptr(family),
		uintptr(class),
		0,
	)

	if size == 0 {
		return nil, fmt.Errorf("failed to get TCP table size")
	}

	// Allocate buffer
	buf := make([]byte, size)
	ret, _, _ := procGetExtendedTcpTable.Call(
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&size)),
		0,
		uintptr(family),
		uintptr(class),
		0,
	)

	if ret != 0 {
		return nil, fmt.Errorf("GetExtendedTcpTable failed with code %d", ret)
	}

	table := (*tcpTable)(unsafe.Pointer(&buf[0]))
	numEntries := int(table.NumEntries)
	entries := unsafe.Slice(&table.Table[0], numEntries)

	return entries, nil
}

type udpRow struct {
	LocalAddr uint32
	LocalPort uint32
	OwningPid uint32
}

type udpTable struct {
	NumEntries uint32
	Table      [1]udpRow
}

// getUDPStats retrieves UDP connection information
func getUDPStats() ([]udpRow, error) {
	var size uint32
	family := uint32(windows.AF_INET)
	class := uint32(1) // UDP_TABLE_OWNER_PID

	// Get required size
	procGetExtendedUdpTable.Call(
		0,
		uintptr(unsafe.Pointer(&size)),
		0,
		uintptr(family),
		uintptr(class),
		0,
	)

	if size == 0 {
		return nil, fmt.Errorf("failed to get UDP table size")
	}

	// Allocate buffer
	buf := make([]byte, size)
	ret, _, _ := procGetExtendedUdpTable.Call(
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&size)),
		0,
		uintptr(family),
		uintptr(class),
		0,
	)

	if ret != 0 {
		return nil, fmt.Errorf("GetExtendedUdpTable failed with code %d", ret)
	}

	table := (*udpTable)(unsafe.Pointer(&buf[0]))
	numEntries := int(table.NumEntries)
	entries := unsafe.Slice(&table.Table[0], numEntries)

	return entries, nil
}

// getProcessPath retrieves the executable path for a process ID
func getProcessPath(pid uint32) string {
	const PROCESS_QUERY_LIMITED_INFORMATION = 0x1000

	handle, err := windows.OpenProcess(PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if err != nil {
		return ""
	}
	defer windows.CloseHandle(handle)

	var buf [windows.MAX_PATH]uint16
	size := uint32(len(buf))

	err = windows.QueryFullProcessImageName(handle, 0, &buf[0], &size)
	if err != nil {
		return ""
	}

	path := syscall.UTF16ToString(buf[:size])
	return filepath.Base(path)
}
