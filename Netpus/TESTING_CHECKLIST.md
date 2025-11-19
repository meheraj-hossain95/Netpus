# Octopus Network Monitor - Testing Checklist

**Version:** 1.0.0  
**Platform:** Windows & Linux  
**Last Updated:** October 21, 2025

---

## 📋 Table of Contents

1. [Pre-Test Setup](#pre-test-setup)
2. [Installation Testing](#installation-testing)
3. [Core Functionality Testing](#core-functionality-testing)
4. [User Interface Testing](#user-interface-testing)
5. [Performance Testing](#performance-testing)
6. [Data Management Testing](#data-management-testing)
7. [System Integration Testing](#system-integration-testing)
8. [Security & Privacy Testing](#security--privacy-testing)
9. [Uninstallation Testing](#uninstallation-testing)
10. [Edge Cases & Error Handling](#edge-cases--error-handling)
11. [Cross-Platform Testing](#cross-platform-testing)

---

## Pre-Test Setup

### Environment Preparation

- [ ] **Test Machine 1:** Clean Windows 10/11 installation
- [ ] **Test Machine 2:** Clean Windows 10/11 with existing applications
- [ ] **Test Machine 3:** Linux (Ubuntu/Debian)
- [ ] **Test Machine 4:** Linux (Fedora/Arch)
- [ ] **Network Configuration:** Active internet connection
- [ ] **Disk Space:** At least 500MB free space
- [ ] **Permissions:** Standard user account (non-admin)
- [ ] **Backup:** System restore point created

### Required Tools

- [ ] Task Manager (Windows) / System Monitor (Linux)
- [ ] Process Explorer / htop
- [ ] Network traffic generator (browser, download manager, etc.)
- [ ] Database viewer (DB Browser for SQLite)
- [ ] Performance monitoring tools

---

## 1. Installation Testing

### 1.1 Portable Mode (No Installation)

#### Windows
- [ ] **TC-001:** Download `octopus.exe` to Downloads folder
- [ ] **TC-002:** Double-click to run - application launches successfully
- [ ] **TC-003:** Verify window appears with title "Octopus Network Monitor"
- [ ] **TC-004:** Close and rerun - no errors
- [ ] **TC-005:** Run from Desktop location - works correctly
- [ ] **TC-006:** Run from custom folder (C:\Tools\) - works correctly
- [ ] **TC-007:** Run from network drive - verify behavior
- [ ] **TC-008:** Run from path with spaces (C:\Program Files\Octopus\) - works correctly
- [ ] **TC-009:** Run from path with special characters - works correctly

#### Linux
- [ ] **TC-010:** Download `octopus` binary to Downloads folder
- [ ] **TC-011:** Make executable: `chmod +x octopus`
- [ ] **TC-012:** Run with `./octopus` - application launches successfully
- [ ] **TC-013:** Run from /opt/ directory - works correctly
- [ ] **TC-014:** Run from home directory - works correctly
- [ ] **TC-015:** Verify file permissions are correct (not requiring root)

### 1.2 Installation with Shortcuts

#### Windows
- [ ] **TC-016:** Run `octopus.exe --install` from command prompt
- [ ] **TC-017:** Verify success message appears
- [ ] **TC-018:** Check Desktop shortcut exists
- [ ] **TC-019:** Check Start Menu shortcut exists in correct location
- [ ] **TC-020:** Desktop shortcut launches application correctly
- [ ] **TC-021:** Start Menu shortcut launches application correctly
- [ ] **TC-022:** Shortcuts have correct icon
- [ ] **TC-023:** Right-click shortcut properties show correct path
- [ ] **TC-024:** Install from path with spaces - shortcuts work correctly

#### Linux
- [ ] **TC-025:** Run `./octopus --install` from terminal
- [ ] **TC-026:** Verify success message appears
- [ ] **TC-027:** Check application menu entry exists (~/.local/share/applications/)
- [ ] **TC-028:** Application appears in system launcher/menu
- [ ] **TC-029:** Launch from application menu works correctly
- [ ] **TC-030:** Desktop file has correct permissions
- [ ] **TC-031:** Icon displays correctly in launcher

### 1.3 Version Check
- [ ] **TC-032:** Run `octopus.exe --version` (Windows)
- [ ] **TC-033:** Run `./octopus --version` (Linux)
- [ ] **TC-034:** Verify version number displays correctly (1.0.0)
- [ ] **TC-035:** Verify platform information displays correctly

---

## 2. Core Functionality Testing

### 2.1 First Launch

- [ ] **TC-036:** Application launches without errors
- [ ] **TC-037:** Database is created in correct location:
  - Windows: `%APPDATA%\octopus\octopus.db`
  - Linux: `~/.local/share/octopus/octopus.db`
- [ ] **TC-038:** Database file is created with correct permissions
- [ ] **TC-039:** Config file is created
- [ ] **TC-040:** Application window appears with default size (1440x768)
- [ ] **TC-041:** Dashboard page loads by default
- [ ] **TC-042:** System tray icon appears
- [ ] **TC-043:** No console errors in developer tools

### 2.2 Network Monitoring

- [ ] **TC-044:** Dashboard shows "No active applications" initially
- [ ] **TC-045:** Open browser - browser appears in monitoring list within 5 seconds
- [ ] **TC-046:** Download a file - download speed is tracked correctly
- [ ] **TC-047:** Upload a file - upload speed is tracked correctly
- [ ] **TC-048:** Process name is displayed correctly
- [ ] **TC-049:** Process icon/image is displayed (if available)
- [ ] **TC-050:** Upload speed shows with "↑" indicator
- [ ] **TC-051:** Download speed shows with "↓" indicator
- [ ] **TC-052:** Total data transferred displays correctly
- [ ] **TC-053:** Speed updates every 1 second
- [ ] **TC-054:** Multiple applications monitored simultaneously
- [ ] **TC-055:** Process disappears from list 30 seconds after network activity stops
- [ ] **TC-056:** Speed values are formatted correctly (KB/s, MB/s, GB/s)

### 2.3 Data Recording

- [ ] **TC-057:** Network activity is saved to database every 10 seconds
- [ ] **TC-058:** Open database file - verify data is being written
- [ ] **TC-059:** Check table structure is correct
- [ ] **TC-060:** Verify timestamps are accurate
- [ ] **TC-061:** Process names are stored correctly
- [ ] **TC-062:** Upload/download values are accurate
- [ ] **TC-063:** Long running process data accumulates correctly
- [ ] **TC-064:** Database grows appropriately with usage

### 2.4 Real-Time Display

- [ ] **TC-065:** Start large download - verify real-time speed display
- [ ] **TC-066:** Pause download - speed drops to zero
- [ ] **TC-067:** Resume download - speed resumes correctly
- [ ] **TC-068:** Multiple concurrent downloads - all tracked separately
- [ ] **TC-069:** Background uploads (cloud sync) are detected
- [ ] **TC-070:** System processes are monitored
- [ ] **TC-071:** Browser tabs in same process aggregate correctly
- [ ] **TC-072:** Speeds are human-readable (not showing bytes)

---

## 3. User Interface Testing

### 3.1 Dashboard Page

- [ ] **TC-073:** Dashboard loads within 2 seconds
- [ ] **TC-074:** Table headers are visible and properly labeled:
  - Application
  - Upload Speed
  - Download Speed
  - Total Uploaded
  - Total Downloaded
- [ ] **TC-075:** Click "Application" header - sorts alphabetically
- [ ] **TC-076:** Click "Upload Speed" header - sorts by upload speed
- [ ] **TC-077:** Click "Download Speed" header - sorts by download speed
- [ ] **TC-078:** Click header again - reverses sort order
- [ ] **TC-079:** Ascending/descending indicator shows correctly
- [ ] **TC-080:** Table rows alternate colors for readability
- [ ] **TC-081:** Hover over row - highlights correctly
- [ ] **TC-082:** Long application names are handled properly (no overflow)
- [ ] **TC-083:** Large numbers are formatted with commas
- [ ] **TC-084:** Empty state shows "No active applications" message
- [ ] **TC-085:** Scrollbar appears when many applications are listed
- [ ] **TC-086:** Smooth scrolling works correctly

### 3.2 Usage History Page

- [ ] **TC-087:** Click "Usage History" tab - page loads
- [ ] **TC-088:** Search box is visible and functional
- [ ] **TC-089:** Type application name - filters results immediately
- [ ] **TC-090:** Search is case-insensitive
- [ ] **TC-091:** Clear search - all results return
- [ ] **TC-092:** Historical data is displayed in table format
- [ ] **TC-093:** Table shows: Application, Date/Time, Upload, Download
- [ ] **TC-094:** Data is sorted by date (newest first) by default
- [ ] **TC-095:** Click column headers to sort
- [ ] **TC-096:** Pagination works if many records exist
- [ ] **TC-097:** Date format is readable (YYYY-MM-DD HH:MM:SS)
- [ ] **TC-098:** Data matches what was monitored in real-time
- [ ] **TC-099:** Large datasets load without freezing UI
- [ ] **TC-100:** Export functionality (if available) works correctly

### 3.3 Settings Page

- [ ] **TC-101:** Click "Settings" tab - page loads
- [ ] **TC-102:** All settings sections are visible
- [ ] **TC-103:** Auto-start toggle is present
- [ ] **TC-104:** Theme selector is present (Auto/Light/Dark)
- [ ] **TC-105:** Data retention dropdown is present
- [ ] **TC-106:** Database size is displayed
- [ ] **TC-107:** "Clear Old Data" button is present
- [ ] **TC-108:** Settings are organized logically
- [ ] **TC-109:** Tooltips/help text explain each setting
- [ ] **TC-110:** Settings page is responsive

### 3.4 Navigation

- [ ] **TC-111:** Tab navigation between Dashboard/History/Settings works smoothly
- [ ] **TC-112:** Active tab is highlighted
- [ ] **TC-113:** Browser-style back/forward (if applicable) works
- [ ] **TC-114:** Keyboard navigation (Tab key) works
- [ ] **TC-115:** Enter key activates buttons
- [ ] **TC-116:** Escape key closes dialogs/modals

### 3.5 Theme & Appearance

- [ ] **TC-117:** Default theme matches system theme (if set to Auto)
- [ ] **TC-118:** Switch to Light theme - UI updates immediately
- [ ] **TC-119:** Switch to Dark theme - UI updates immediately
- [ ] **TC-120:** All text is readable in Light theme
- [ ] **TC-121:** All text is readable in Dark theme
- [ ] **TC-122:** Theme preference is saved after restart
- [ ] **TC-123:** No flash of wrong theme on startup
- [ ] **TC-124:** Colors are consistent across all pages
- [ ] **TC-125:** System theme changes are detected (if set to Auto)

### 3.6 Window Management

- [ ] **TC-126:** Minimize window - application stays in tray
- [ ] **TC-127:** Restore from tray - window reappears
- [ ] **TC-128:** Maximize window - UI scales correctly
- [ ] **TC-129:** Resize window - content adjusts responsively
- [ ] **TC-130:** Window position is remembered between sessions
- [ ] **TC-131:** Window size is remembered between sessions
- [ ] **TC-132:** Close button minimizes to tray (not exit)
- [ ] **TC-133:** Multi-monitor setup - window appears on correct screen
- [ ] **TC-134:** Window dragging works smoothly
- [ ] **TC-135:** Window borders and title bar display correctly

---

## 4. Performance Testing

### 4.1 Resource Usage

- [ ] **TC-136:** Check RAM usage at idle - should be ~50MB or less
- [ ] **TC-137:** Check RAM usage with 10 active processes - remains reasonable
- [ ] **TC-138:** Check RAM usage with 50+ active processes - no memory leak
- [ ] **TC-139:** CPU usage at idle - should be <1%
- [ ] **TC-140:** CPU usage during monitoring - should be <5%
- [ ] **TC-141:** Run for 1 hour - no performance degradation
- [ ] **TC-142:** Run for 24 hours - remains stable
- [ ] **TC-143:** Run for 1 week - remains stable
- [ ] **TC-144:** Database file size grows predictably
- [ ] **TC-145:** No memory leaks detected over extended usage

### 4.2 Update Speed

- [ ] **TC-146:** Dashboard updates every 1 second exactly
- [ ] **TC-147:** No lag when updating with many processes
- [ ] **TC-148:** Sorting doesn't interrupt real-time updates
- [ ] **TC-149:** Switching tabs doesn't freeze UI
- [ ] **TC-150:** Background database writes don't block UI

### 4.3 Database Performance

- [ ] **TC-151:** Initial database creation is fast (<1 second)
- [ ] **TC-152:** Queries execute quickly (<100ms)
- [ ] **TC-153:** Database writes complete in background
- [ ] **TC-154:** No database locks causing errors
- [ ] **TC-155:** WAL mode is enabled for concurrency
- [ ] **TC-156:** Database vacuum completes successfully
- [ ] **TC-157:** Large database (1GB+) doesn't slow queries

### 4.4 Stress Testing

- [ ] **TC-158:** Monitor 100+ concurrent network processes
- [ ] **TC-159:** Generate 1GB/s network traffic - application remains responsive
- [ ] **TC-160:** Rapidly start/stop 20 applications - no crashes
- [ ] **TC-161:** Fill database to 5GB - application still functional
- [ ] **TC-162:** Run on machine with 100+ processes running
- [ ] **TC-163:** Low memory environment (2GB RAM) - graceful degradation

---

## 5. Data Management Testing

### 5.1 Data Retention

- [ ] **TC-164:** Set retention to 7 days - verify old data is deleted
- [ ] **TC-165:** Set retention to 30 days - verify behavior
- [ ] **TC-166:** Set retention to 90 days - verify behavior
- [ ] **TC-167:** Set retention to Forever - no data is deleted
- [ ] **TC-168:** Change retention setting - takes effect immediately
- [ ] **TC-169:** Manual cleanup runs after setting change
- [ ] **TC-170:** Automatic cleanup runs daily

### 5.2 Manual Data Clearing

- [ ] **TC-171:** Click "Clear Old Data" button
- [ ] **TC-172:** Confirmation dialog appears
- [ ] **TC-173:** Confirm - data older than retention period is deleted
- [ ] **TC-174:** Cancel - no data is deleted
- [ ] **TC-175:** Database size reduces after clearing
- [ ] **TC-176:** Current monitoring data is preserved
- [ ] **TC-177:** Success message appears after clearing

### 5.3 Database Integrity

- [ ] **TC-178:** Force-close application - database remains intact
- [ ] **TC-179:** System crash - database recovers on restart
- [ ] **TC-180:** Run database integrity check - no corruption
- [ ] **TC-181:** Backup database file - can be restored
- [ ] **TC-182:** Delete database file - recreated on next launch
- [ ] **TC-183:** Corrupt database file - application handles gracefully
- [ ] **TC-184:** Multiple rapid writes - no data loss

---

## 6. System Integration Testing

### 6.1 System Tray

#### Windows
- [ ] **TC-185:** System tray icon appears after launch
- [ ] **TC-186:** Icon has correct appearance
- [ ] **TC-187:** Icon updates with live speed (if implemented)
- [ ] **TC-188:** Hover over icon - tooltip shows current speeds
- [ ] **TC-189:** Right-click icon - context menu appears
- [ ] **TC-190:** Context menu shows: Show/Hide, Pause, Resume, Settings, Quit
- [ ] **TC-191:** Click "Show" - window appears
- [ ] **TC-192:** Click "Hide" - window minimizes to tray
- [ ] **TC-193:** Click "Pause" - monitoring pauses
- [ ] **TC-194:** Click "Resume" - monitoring resumes
- [ ] **TC-195:** Click "Settings" - Settings page opens
- [ ] **TC-196:** Click "Quit" - application closes completely
- [ ] **TC-197:** Left-click icon - toggles window visibility
- [ ] **TC-198:** Icon persists when window is closed
- [ ] **TC-199:** Icon is removed when app quits

#### Linux
- [ ] **TC-200:** System tray icon appears (GNOME/KDE/XFCE)
- [ ] **TC-201:** AppIndicator works on Ubuntu
- [ ] **TC-202:** Icon displays correctly in different themes
- [ ] **TC-203:** All tray menu functions work
- [ ] **TC-204:** Icon updates correctly (if animated)

### 6.2 Auto-Start

#### Windows
- [ ] **TC-205:** Enable auto-start in Settings
- [ ] **TC-206:** Restart computer - application starts automatically
- [ ] **TC-207:** Application starts minimized to tray
- [ ] **TC-208:** Disable auto-start - verify removed from startup
- [ ] **TC-209:** Check registry/startup folder for correct entry
- [ ] **TC-210:** Auto-start works with User Account Control (UAC)

#### Linux
- [ ] **TC-211:** Enable auto-start in Settings
- [ ] **TC-212:** Restart/re-login - application starts automatically
- [ ] **TC-213:** Check ~/.config/autostart/ for .desktop file
- [ ] **TC-214:** Disable auto-start - file is removed
- [ ] **TC-215:** Works across different desktop environments

### 6.3 Multiple Instance Protection

- [ ] **TC-216:** Launch first instance - starts normally
- [ ] **TC-217:** Try to launch second instance - error message appears
- [ ] **TC-218:** Error message is clear: "Only one instance can run"
- [ ] **TC-219:** Lock file is created in temp directory
- [ ] **TC-220:** First instance quits - lock file is deleted
- [ ] **TC-221:** After first instance quits - second instance can launch
- [ ] **TC-222:** Kill first instance (Task Manager) - second instance can launch after lock cleanup
- [ ] **TC-223:** Lock file contains PID of running process

### 6.4 File Associations

- [ ] **TC-224:** Database file opens in correct application
- [ ] **TC-225:** Config file is valid JSON
- [ ] **TC-226:** Log files (if any) are readable

---

## 7. Security & Privacy Testing

### 7.1 Network Privacy

- [ ] **TC-227:** Use Wireshark - verify NO outgoing connections
- [ ] **TC-228:** Check firewall logs - no external connections
- [ ] **TC-229:** Block internet - application works normally
- [ ] **TC-230:** No telemetry or analytics endpoints contacted
- [ ] **TC-231:** No DNS lookups for tracking domains
- [ ] **TC-232:** All data stays local

### 7.2 Data Security

- [ ] **TC-233:** Database file has correct permissions (user read/write only)
- [ ] **TC-234:** Config file has correct permissions
- [ ] **TC-235:** No sensitive data logged to system logs
- [ ] **TC-236:** Process names are stored securely
- [ ] **TC-237:** No password or credential capture
- [ ] **TC-238:** Database is not encrypted (documented as local storage)

### 7.3 Application Security

- [ ] **TC-239:** No SQL injection vulnerabilities in search
- [ ] **TC-240:** No XSS vulnerabilities in UI
- [ ] **TC-241:** File paths are validated
- [ ] **TC-242:** No arbitrary code execution vectors
- [ ] **TC-243:** WebView is sandboxed properly
- [ ] **TC-244:** Application runs without admin privileges

### 7.4 System Impact

- [ ] **TC-245:** Application doesn't modify system files
- [ ] **TC-246:** Application doesn't modify registry (except auto-start)
- [ ] **TC-247:** Application doesn't intercept network traffic
- [ ] **TC-248:** Application only reads network statistics
- [ ] **TC-249:** Uninstall removes all traces (if requested)

---

## 8. Uninstallation Testing

### 8.1 Uninstall Command

#### Windows
- [ ] **TC-250:** Run `octopus.exe --uninstall`
- [ ] **TC-251:** Success message appears
- [ ] **TC-252:** Desktop shortcut is removed
- [ ] **TC-253:** Start Menu shortcut is removed
- [ ] **TC-254:** Auto-start entry is removed
- [ ] **TC-255:** Application executable remains (as documented)
- [ ] **TC-256:** Database remains (as documented)
- [ ] **TC-257:** Can run application again after uninstall

#### Linux
- [ ] **TC-258:** Run `./octopus --uninstall`
- [ ] **TC-259:** Success message appears
- [ ] **TC-260:** .desktop file is removed from applications folder
- [ ] **TC-261:** Auto-start .desktop file is removed
- [ ] **TC-262:** Application binary remains

### 8.2 Manual Uninstall

- [ ] **TC-263:** Close application completely
- [ ] **TC-264:** Delete executable file
- [ ] **TC-265:** Delete database folder
  - Windows: %APPDATA%\octopus
  - Linux: ~/.local/share/octopus
- [ ] **TC-266:** Delete config folder (if separate)
- [ ] **TC-267:** Check for remaining files - none found
- [ ] **TC-268:** No residual registry entries (Windows)

### 8.3 Clean Uninstall

- [ ] **TC-269:** Follow complete uninstall guide from README
- [ ] **TC-270:** All components removed successfully
- [ ] **TC-271:** No orphaned shortcuts
- [ ] **TC-272:** No orphaned processes
- [ ] **TC-273:** System is in clean state

---

## 9. Edge Cases & Error Handling

### 9.1 System Conditions

- [ ] **TC-274:** Low disk space (<100MB) - graceful warning
- [ ] **TC-275:** Full disk - error handling works
- [ ] **TC-276:** Read-only filesystem - appropriate error message
- [ ] **TC-277:** Network adapter disabled - application handles gracefully
- [ ] **TC-278:** No network adapters - application shows appropriate message
- [ ] **TC-279:** VPN connection - monitoring works correctly
- [ ] **TC-280:** Virtual network adapters - handled correctly
- [ ] **TC-281:** Multiple network adapters - all monitored

### 9.2 Data Scenarios

- [ ] **TC-282:** Database file locked by another process - error message
- [ ] **TC-283:** Database file deleted while running - recreated
- [ ] **TC-284:** Database file corrupted - recovery attempted
- [ ] **TC-285:** Config file corrupted - defaults loaded
- [ ] **TC-286:** Config file missing - recreated with defaults
- [ ] **TC-287:** Empty database - application functions normally
- [ ] **TC-288:** Very large database (10GB+) - performance impact documented

### 9.3 User Actions

- [ ] **TC-289:** Spam-click buttons rapidly - no crashes
- [ ] **TC-290:** Open 100+ browser tabs - all tracked or limit handled
- [ ] **TC-291:** Force close during database write - data integrity maintained
- [ ] **TC-292:** Change all settings rapidly - no conflicts
- [ ] **TC-293:** Switch themes rapidly - no visual glitches
- [ ] **TC-294:** Resize window to minimum - UI remains usable
- [ ] **TC-295:** Resize window to maximum - UI scales properly

### 9.4 Application State

- [ ] **TC-296:** Pause monitoring - data collection stops
- [ ] **TC-297:** Resume monitoring - data collection resumes
- [ ] **TC-298:** System sleep/hibernate - application resumes correctly
- [ ] **TC-299:** System wake - monitoring resumes automatically
- [ ] **TC-300:** User logout/login - application state preserved (if auto-start)
- [ ] **TC-301:** Fast user switching - application continues running
- [ ] **TC-302:** Remote desktop connection - monitoring works

### 9.5 Unusual Scenarios

- [ ] **TC-303:** System clock changes - timestamps remain accurate
- [ ] **TC-304:** Timezone change - data recorded correctly
- [ ] **TC-305:** Daylight saving time - no issues
- [ ] **TC-306:** Very long application names (200+ chars) - handled correctly
- [ ] **TC-307:** Application names with special characters - displayed correctly
- [ ] **TC-308:** Unicode characters in process names - rendered correctly
- [ ] **TC-309:** Emojis in process names - handled appropriately

---

## 10. Cross-Platform Testing

### 10.1 Windows Specific

- [ ] **TC-310:** Windows 10 (21H2) - fully functional
- [ ] **TC-311:** Windows 10 (22H2) - fully functional
- [ ] **TC-312:** Windows 11 (21H2) - fully functional
- [ ] **TC-313:** Windows 11 (22H2) - fully functional
- [ ] **TC-314:** Windows 11 (23H2) - fully functional
- [ ] **TC-315:** WebView2 is available (or bundled)
- [ ] **TC-316:** High DPI displays - UI scales correctly
- [ ] **TC-317:** Multiple monitors - no issues
- [ ] **TC-318:** Touchscreen input - works if applicable

### 10.2 Linux Specific

#### Desktop Environments
- [ ] **TC-319:** GNOME - fully functional
- [ ] **TC-320:** KDE Plasma - fully functional
- [ ] **TC-321:** XFCE - fully functional
- [ ] **TC-322:** Cinnamon - fully functional
- [ ] **TC-323:** MATE - fully functional
- [ ] **TC-324:** i3/tiling WM - works with limitations documented

#### Distributions
- [ ] **TC-325:** Ubuntu 20.04 LTS - works
- [ ] **TC-326:** Ubuntu 22.04 LTS - works
- [ ] **TC-327:** Ubuntu 24.04 LTS - works
- [ ] **TC-328:** Debian 11 - works
- [ ] **TC-329:** Debian 12 - works
- [ ] **TC-330:** Fedora (latest) - works
- [ ] **TC-331:** Arch Linux - works
- [ ] **TC-332:** Linux Mint - works

#### Linux Capabilities
- [ ] **TC-333:** Network monitoring works without root (if applicable)
- [ ] **TC-334:** Or requires capabilities (documented)
- [ ] **TC-335:** Wayland session - works correctly
- [ ] **TC-336:** X11 session - works correctly
- [ ] **TC-337:** HiDPI scaling - renders correctly

---

## 11. Regression Testing

After any code changes, run this subset:

### Critical Path Tests
- [ ] **RT-001:** Application launches successfully
- [ ] **RT-002:** Network monitoring works in real-time
- [ ] **RT-003:** Data is saved to database
- [ ] **RT-004:** System tray functions correctly
- [ ] **RT-005:** Settings are saved and loaded
- [ ] **RT-006:** Theme switching works
- [ ] **RT-007:** Auto-start toggle works
- [ ] **RT-008:** Data retention works
- [ ] **RT-009:** Multiple instance protection works
- [ ] **RT-010:** Application closes cleanly

### Smoke Tests (Quick Validation)
- [ ] **ST-001:** Build completes without errors
- [ ] **ST-002:** Application launches
- [ ] **ST-003:** Main UI loads
- [ ] **ST-004:** Network monitoring starts
- [ ] **ST-005:** Application closes without crash

---

## 12. Acceptance Criteria

Before releasing version 1.0.0, ALL of the following must be TRUE:

### Functionality
- [ ] All core features work as documented
- [ ] No critical bugs present
- [ ] No data loss scenarios
- [ ] Performance meets specifications (<50MB RAM, <1% CPU idle)

### Quality
- [ ] No crashes during normal use
- [ ] UI is responsive and intuitive
- [ ] Error messages are helpful
- [ ] Documentation is complete and accurate

### Compatibility
- [ ] Works on Windows 10/11
- [ ] Works on Ubuntu/Debian Linux
- [ ] Works on Fedora/Arch Linux
- [ ] Tested on multiple desktop environments (Linux)

### Privacy & Security
- [ ] No external network connections
- [ ] Data stored locally only
- [ ] Appropriate file permissions
- [ ] No unnecessary privileges required

### User Experience
- [ ] Installation is simple (or unnecessary)
- [ ] First-run experience is smooth
- [ ] Settings are discoverable
- [ ] Uninstallation is complete

---

## Test Execution Summary

**Tester Name:** _____________________  
**Test Date:** _____________________  
**Version Tested:** _____________________  
**Platform:** _____________________

### Results

| Category | Total Tests | Passed | Failed | Skipped | Pass Rate |
|----------|-------------|--------|--------|---------|-----------|
| Installation | 35 | | | | |
| Core Functionality | 28 | | | | |
| User Interface | 65 | | | | |
| Performance | 22 | | | | |
| Data Management | 14 | | | | |
| System Integration | 39 | | | | |
| Security & Privacy | 23 | | | | |
| Uninstallation | 24 | | | | |
| Edge Cases | 37 | | | | |
| Cross-Platform | 27 | | | | |
| **TOTAL** | **314** | | | | |

### Critical Issues Found

1. _____________________
2. _____________________
3. _____________________

### Recommendations

_____________________

### Sign-Off

**Tester Signature:** _____________________  
**Date:** _____________________

---

## Automated Testing Notes

For future versions, consider implementing:

- Unit tests for Go backend functions
- Integration tests for database operations
- UI automation tests (Selenium/Playwright)
- Performance benchmarks
- Continuous integration testing

---

## Additional Resources

- **Bug Report Template:** Include steps to reproduce, expected vs actual behavior, system info
- **Performance Baseline:** Document acceptable metrics for comparison
- **Test Data Generator:** Create scripts to simulate network traffic
- **Screenshot Comparison:** For UI regression testing

---

**Document Version:** 1.0  
**Last Updated:** October 21, 2025  
**Maintained By:** QA Team
