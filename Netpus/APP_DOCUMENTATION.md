# Octopus Network Monitor - Application Documentation

## Table of Contents
1. [Overview](#overview)
2. [System Architecture](#system-architecture)
3. [Database Structure](#database-structure)
4. [Data Collection Process](#data-collection-process)
5. [Background Services](#background-services)
6. [Platform-Specific Details](#platform-specific-details)
7. [Configuration & Storage](#configuration--storage)
8. [Performance Characteristics](#performance-characteristics)

---

## Overview

**Octopus** is a real-time network bandwidth monitoring application that tracks per-application network usage on Windows and Linux systems. Built with Go and Wails v2, it provides a native desktop experience with minimal resource footprint.

**Key Features:**
- Real-time network monitoring (1-second update interval)
- Per-application bandwidth tracking
- Historical usage data with configurable retention
- System tray integration with background monitoring
- Auto-start capability
- Dark/Light theme support
- Completely offline - no telemetry or external connections

**Technology Stack:**
- **Backend:** Go 1.22+
- **Framework:** Wails v2 (Go + HTML/CSS/JavaScript)
- **Database:** SQLite with WAL mode
- **UI:** Vanilla JavaScript (no frameworks)
- **Supported Platforms:** Windows 10+, Linux (X11/Wayland)

---

## System Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────┐
│               Frontend (HTML/CSS/JavaScript)             │
│  • Dashboard (real-time stats, 1-second updates)        │
│  • Usage History (historical data analysis)             │
│  • Settings (configuration management)                  │
└────────────────┬────────────────────────────────────────┘
                 │
                 │ Wails Bridge (IPC over WebSocket)
                 │
┌────────────────▼────────────────────────────────────────┐
│                   Application Core (Go)                  │
│  • app.go - Business logic & API endpoints              │
│  • main.go - Entry point & lifecycle management         │
└──┬──────────┬──────────┬──────────┬──────────┬─────────┘
   │          │          │          │          │
   ▼          ▼          ▼          ▼          ▼
┌──────┐ ┌─────────┐ ┌──────┐ ┌───────┐ ┌──────────┐
│Monitor│ │Database │ │Config│ │ Tray  │ │Installer │
│ Module│ │ Module  │ │Module│ │Module │ │  Module  │
└───┬───┘ └────┬────┘ └───┬──┘ └───┬───┘ └─────┬────┘
    │          │          │        │           │
┌───▼──────────▼──────────▼────────▼───────────▼─────┐
│            Operating System APIs                     │
│  • Network Interface Stats (Windows API / /proc)     │
│  • Process Information (PID, executable paths)       │
│  • System Tray (Windows Shell / Linux systray)       │
│  • File System (SQLite database, config storage)     │
└──────────────────────────────────────────────────────┘
```

### Module Responsibilities

#### 1. **Monitor Module** (`internal/monitor/`)
- Collects network statistics every 1 second
- Maps network traffic to specific processes
- Calculates upload/download speeds
- Batches data writes to database (every 10 seconds)
- Manages process lifecycle (cleanup inactive processes after 30 seconds)

#### 2. **Database Module** (`internal/database/`)
- SQLite database with WAL mode for concurrent access
- Stores raw usage records, daily summaries, and app metadata
- Handles data retention policies
- Provides query interfaces for historical data
- Automatic cleanup and vacuum operations

#### 3. **Configuration Module** (`internal/utils/`)
- Loads and persists user settings
- Manages database path resolution
- Provides utility functions (formatting, path handling)

#### 4. **Tray Module** (`internal/tray/`)
- System tray icon integration
- Real-time tooltip updates (shows current speeds)
- Context menu for quick actions
- Window show/hide functionality

#### 5. **Installer Module** (`internal/installer/`)
- Creates desktop and start menu shortcuts
- Handles installation/uninstallation via CLI flags

#### 6. **Autostart Module** (`internal/autostart/`)
- Manages system startup integration
- Windows: Registry entries
- Linux: XDG autostart desktop files

---

## Database Structure

### SQLite Database Schema

**Location:**
- Windows: `%APPDATA%\octopus\octopus.db`
- Linux: `~/.local/share/octopus/octopus.db`

**Connection Settings:**
```go
Journal Mode: WAL (Write-Ahead Logging)
Synchronous: NORMAL
Busy Timeout: 5000ms
Max Open Connections: 3
Max Idle Connections: 2
Connection Max Lifetime: 1 hour
```

### Tables

#### 1. `usage_records`
Stores raw network usage data collected every second.

```sql
CREATE TABLE usage_records (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    app_name TEXT NOT NULL,
    process_id INTEGER,
    upload_bytes INTEGER NOT NULL,
    download_bytes INTEGER NOT NULL,
    timestamp INTEGER NOT NULL
);

CREATE INDEX idx_usage_timestamp ON usage_records(timestamp);
CREATE INDEX idx_usage_app ON usage_records(app_name);
```

**Fields:**
- `id`: Unique record identifier
- `app_name`: Executable name (e.g., "firefox.exe", "chrome")
- `process_id`: Operating system process ID
- `upload_bytes`: Bytes uploaded since last collection
- `download_bytes`: Bytes downloaded since last collection
- `timestamp`: Unix timestamp (seconds since epoch)

**Retention:** Records are automatically deleted based on user settings:
- 7/30/90 days retention (configurable)
- Forever (retention = -1)
- Testing mode: 1 minute (retention = 0)

#### 2. `daily_summaries`
Aggregated daily totals for quick dashboard queries.

```sql
CREATE TABLE daily_summaries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date TEXT UNIQUE NOT NULL,
    total_upload INTEGER NOT NULL,
    total_download INTEGER NOT NULL
);

CREATE INDEX idx_daily_date ON daily_summaries(date);
```

**Fields:**
- `date`: Date string in format "YYYY-MM-DD"
- `total_upload`: Total bytes uploaded on this date
- `total_download`: Total bytes downloaded on this date

**Purpose:** Provides fast "Today's Usage" statistics without scanning all records.

#### 3. `app_metadata`
Tracks application information and activity.

```sql
CREATE TABLE app_metadata (
    app_name TEXT PRIMARY KEY,
    executable_path TEXT,
    first_seen INTEGER NOT NULL,
    last_seen INTEGER NOT NULL
);
```

**Fields:**
- `app_name`: Executable name (primary key)
- `executable_path`: Full path to executable
- `first_seen`: Unix timestamp when first detected
- `last_seen`: Unix timestamp of most recent activity

**Purpose:** Provides "Last Seen" information in Usage History page.

#### 4. `settings`
User configuration stored as key-value pairs.

```sql
CREATE TABLE settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL
);
```

**Stored Settings:**
- `auto_start`: "true" or "false"
- `theme`: "auto", "light", or "dark"
- `data_retention`: "-1" (forever), "0" (testing), "7", "30", or "90" (days)
- `network_interface`: Reserved for future use

---

## Data Collection Process

### Real-Time Monitoring Flow

```
┌─────────────────────────────────────────────────────────────┐
│                     Monitoring Loop                         │
│                  (Every 1 second)                           │
└──────┬──────────────────────────────────────────────────────┘
       │
       ▼
┌─────────────────────────────────────────┐
│  1. Collect Network Statistics          │
│     • Windows: GetExtendedTcpTable()    │
│                GetExtendedUdpTable()    │
│                GetIfTable2()            │
│     • Linux: Parse /proc/net/*          │
│              Parse /proc/[pid]/net/dev  │
└──────┬──────────────────────────────────┘
       │
       ▼
┌─────────────────────────────────────────┐
│  2. Map Traffic to Processes            │
│     • Windows: Weight by connections    │
│     • Linux: Direct from /proc          │
└──────┬──────────────────────────────────┘
       │
       ▼
┌─────────────────────────────────────────┐
│  3. Calculate Speeds                    │
│     Speed = (Current - Previous) / 1s   │
└──────┬──────────────────────────────────┘
       │
       ▼
┌─────────────────────────────────────────┐
│  4. Update In-Memory Statistics         │
│     • Per-process stats map             │
│     • Total upload/download             │
└──────┬──────────────────────────────────┘
       │
       ▼
┌─────────────────────────────────────────┐
│  5. Add to Batch Buffer                 │
│     Batched records queued              │
└──────┬──────────────────────────────────┘
       │
       │ Every 10 seconds
       ▼
┌─────────────────────────────────────────┐
│  6. Flush to Database                   │
│     • Batch insert usage_records        │
│     • Update daily_summaries            │
│     • Update app_metadata               │
└─────────────────────────────────────────┘
```

### Initialization Process

1. **Baseline Collection:**
   - Collect network stats twice with 1-second delay
   - Establishes baseline for speed calculation
   - First collection stores initial byte counts

2. **Continuous Monitoring:**
   - Every 1 second: collect, calculate, update
   - Every 10 seconds: flush batch to database
   - Every 30 seconds: cleanup inactive processes

### Speed Calculation

```go
uploadSpeed = (currentUploadBytes - previousUploadBytes) / timeDelta
downloadSpeed = (currentDownloadBytes - previousDownloadBytes) / timeDelta
```

Where `timeDelta` is typically 1 second.

---

## Background Services

The application runs several background goroutines for maintenance:

### 1. **Tray Tooltip Updater**
**Frequency:** Every 1 second  
**Purpose:** Updates system tray tooltip with current network speeds

```
Octopus
↑ 1.2 MB/s ↓ 3.4 MB/s
```

### 2. **Daily Cleanup**
**Frequency:** Every 24 hours  
**Actions:**
- Delete expired records (older than retention period)
- Vacuum database to reclaim disk space
- Runs immediately on startup, then every 24 hours

### 3. **Hourly Cleanup**
**Frequency:** Every 30 minutes  
**Actions:**
- Enforce data retention policy
- Delete records older than configured days

### 4. **Minute Cleanup**
**Frequency:** Every 1 minute  
**Actions:**
- Delete expired records (automatic 24-hour cleanup)
- Handle testing mode (1-minute retention)
- More immediate enforcement of retention policies

---

## Platform-Specific Details

### Windows Implementation

#### Network Monitoring (`monitor_windows.go`)

**Windows APIs Used:**
- `GetExtendedTcpTable()` - TCP connection table with process IDs
- `GetExtendedUdpTable()` - UDP connection table with process IDs
- `GetIfTable2()` - Network interface statistics
- `OpenProcess()` - Access process information
- `QueryFullProcessImageNameW()` - Get process executable path

**Algorithm:**
1. Get system-wide network I/O from all interfaces
2. Enumerate all TCP and UDP connections
3. Weight connections by state:
   - TCP ESTABLISHED: High weight (active data transfer)
   - TCP other states: Medium weight
   - UDP: Lower weight (connectionless)
4. Distribute total network bytes based on connection weights
5. Return per-process byte counts

**Advantages:**
- Accurate for active connections
- Low overhead (no per-process tracking)

**Limitations:**
- Estimation-based (not exact per-process)
- May be less accurate for many concurrent connections

#### System Tray (`tray_windows.go`)
- Uses `getlantern/systray` library
- Icon loaded from file or fallback to default
- Context menu with Show/Hide/Pause/Resume/Quit

#### Autostart (`autostart_windows.go`)
- Registry path: `HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Run`
- Key name: "Octopus"
- Value: Full path to executable (quoted)

#### Installer (`install_windows.go`)
- Creates shortcuts via COM automation (WScript.Shell)
- Desktop shortcut: `%USERPROFILE%\Desktop\Octopus.lnk`
- Start Menu shortcut: `%APPDATA%\Microsoft\Windows\Start Menu\Programs\Octopus.lnk`

### Linux Implementation

#### Network Monitoring (`monitor_linux.go`)

**Data Sources:**
- `/proc/net/tcp`, `/proc/net/tcp6` - TCP connection tables
- `/proc/net/udp`, `/proc/net/udp6` - UDP connection tables
- `/proc/[pid]/fd/` - Process file descriptors (socket inodes)
- `/proc/[pid]/net/dev` - Per-process network statistics
- `/proc/[pid]/exe` - Process executable path

**Algorithm:**
1. Parse `/proc/net/tcp*` and `/proc/net/udp*` for socket inodes
2. Map inodes to processes by scanning `/proc/[pid]/fd/*`
3. Read network statistics from `/proc/[pid]/net/dev`
4. Return per-process byte counts

**Advantages:**
- Direct per-process statistics from kernel
- More accurate than Windows estimation

**Limitations:**
- Requires read access to `/proc` (standard for user processes)
- More file I/O overhead

#### System Tray
- Uses `getlantern/systray` library
- Works with both X11 and Wayland desktop environments
- Requires system tray support (most DEs have this)

#### Autostart (`autostart_linux.go`)
- XDG autostart specification
- File: `~/.config/autostart/octopus.desktop`
- Standard Desktop Entry format

#### Installer (`install_linux.go`)
- Application menu entry
- File: `~/.local/share/applications/octopus.desktop`
- Appears in application launcher

---

## Configuration & Storage

### Settings Management

**Default Configuration:**
```go
AutoStart: false
Theme: "auto"
DataRetention: 30 (days)
NetworkInterface: "" (reserved)
```

**Theme Options:**
- `auto`: Follow system theme (light/dark)
- `light`: Always light mode
- `dark`: Always dark mode

**Data Retention Options:**
- `-1`: Forever (never delete old records)
- `0`: 1 minute (testing mode - immediate cleanup)
- `7`, `30`, `90`: Days to keep records

### File Locations

#### Windows
- **Database:** `C:\Users\[Username]\AppData\Roaming\octopus\octopus.db`
- **Lock File:** `C:\Users\[Username]\AppData\Local\Temp\octopus.lock`
- **Autostart:** Registry `HKCU\Software\Microsoft\Windows\CurrentVersion\Run\Octopus`

#### Linux
- **Database:** `~/.local/share/octopus/octopus.db`
- **Lock File:** `/tmp/octopus.lock`
- **Autostart:** `~/.config/autostart/octopus.desktop`
- **App Entry:** `~/.local/share/applications/octopus.desktop`

---

## Performance Characteristics

### Resource Usage

**CPU:**
- Idle: <1% (1-second timer, efficient API calls)
- Active monitoring: 1-3% (depends on number of active processes)
- Database writes: Brief spikes every 10 seconds

**Memory:**
- Base: 20-50 MB (Wails webview + Go runtime)
- Growth: Linear with number of unique applications
- Typical: 50-100 MB with 20-30 tracked applications

**Disk I/O:**
- Reads: On-demand (settings load, historical queries)
- Writes: Batched every 10 seconds
- Database growth: 1-5 MB per day (typical usage)

### Database Performance

**Record Count vs. Size:**
- 1 day: ~86,400 records per app (1-second intervals)
- 30 days: ~2.5M records per app
- Average record size: ~50 bytes
- With compression (SQLite): ~10-20 bytes effective

**Query Performance:**
- Real-time stats: In-memory (instant)
- Today's usage: Single query on daily_summaries (< 1ms)
- Usage history: Aggregation query on usage_records (10-100ms)
- Historical data: Index-optimized (< 50ms for 30 days)

### Optimization Techniques

1. **Batch Writes:** Reduce database I/O by 10x
2. **WAL Mode:** Concurrent reads don't block, writes don't block reads
3. **Indexed Queries:** Fast lookups on timestamp and app_name
4. **Daily Summaries:** Pre-aggregated data for dashboard
5. **In-Memory Cache:** Frontend caches stats to prevent UI flicker
6. **Inactive Cleanup:** Remove processes with no activity for 30 seconds

---

## API Reference (Frontend to Backend)

### Real-Time Data

**GetNetworkStats()** → `map[string]*NetworkStat`
- Returns current network statistics for all active applications
- Called every 1 second by frontend

**GetMonitorStatus()** → `MonitorStatus`
- Returns monitoring state (interval, last update, paused)

**GetTodayStats()** → `{upload: int64, download: int64}`
- Returns today's total usage from daily_summaries table

**Get24HourUsage()** → `{upload: int64, download: int64}`
- Returns rolling 24-hour totals

### Historical Data

**GetNetworkUsageStats()** → `[]AppUsageStat`
- Returns aggregated usage per app based on retention period

**GetHistoricalData(days int)** → `[]DailySummary`
- Returns daily summaries for last N days

### Configuration

**GetSettings()** → `Config`
- Returns current configuration

**UpdateSettings(Config)** → `error`
- Saves new configuration, applies autostart changes

### Control

**PauseMonitoring()** / **ResumeMonitoring()**
- Control monitoring state

**ClearOldData()** → `error`
- Manually delete all old records and vacuum database

**GetDatabaseSize()** → `int64`
- Returns database file size in bytes

**GetDatabaseStats()** → `{size, records, oldest}`
- Returns comprehensive database statistics

---

## Security & Privacy

### Data Privacy
- **No telemetry:** Application is completely offline
- **Local storage:** All data stored on user's machine
- **No network access:** Application only monitors, never transmits

### Security Considerations
- **Process access:** Read-only access to network statistics
- **Database:** Standard file permissions (user-owned)
- **Settings validation:** Input validation on all settings
- **SQL injection:** Parameterized queries throughout
- **No elevated privileges:** Runs with standard user permissions

---

## Troubleshooting

### Common Issues

**"Only one instance can run at a time"**
- Another instance is already running
- Check for `octopus.lock` in temp directory
- Kill existing process or delete lock file

**High CPU usage**
- Check number of active network processes
- Consider pausing monitoring temporarily
- Restart application to clear any stuck processes

**Database growing too large**
- Check data retention setting
- Manually clear old data in Settings
- Consider reducing retention period

**Missing system tray icon (Linux)**
- Ensure desktop environment supports system tray
- Install `libappindicator` or equivalent
- Try different desktop environment

---

## Build Information

**Build Requirements:**
- Go 1.22 or higher
- Wails CLI v2
- Node.js (for frontend build)

**Build Commands:**
```bash
# Windows
build.bat

# Linux
chmod +x build.sh
./build.sh
```

**Output:**
- Windows: `build\bin\octopus.exe`
- Linux: `build/bin/octopus`

---

## License

Copyright © 2025 Octopus  
Licensed under MIT License
