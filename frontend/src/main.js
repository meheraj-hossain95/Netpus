// Netpus Network Monitor - Frontend JavaScript

// Cache for last known stats to prevent UI flickering
let lastStatsCache = {};
let lastUsageCache = [];
let currentPage = 'dashboard';
let monitoringPaused = false;
let currentSortColumn = 'uploadSpeed'; // Default sort by upload speed
let currentSortDirection = 'desc'; // 'asc' or 'desc'
let usageSortColumn = 'totalData'; // Default sort by total data for usage page
let usageSortDirection = 'desc'; // 'asc' or 'desc'

// Initialize application
document.addEventListener('DOMContentLoaded', function () {
    initializeUI();
    loadSettings();
    startDataUpdates();
});

// Initialize UI event handlers
function initializeUI() {
    // Tab navigation
    const tabs = document.querySelectorAll('.nav-tab');
    tabs.forEach(tab => {
        tab.addEventListener('click', () => {
            const page = tab.dataset.page;
            switchPage(page);
        });
    });

    // Dashboard controls
    document.getElementById('pauseBtn')?.addEventListener('click', pauseMonitoring);
    document.getElementById('resumeBtn')?.addEventListener('click', resumeMonitoring);

    // Settings controls - auto-save on change
    document.getElementById('cleanupBtn')?.addEventListener('click', cleanupOldData);
    document.getElementById('autoStartCheck')?.addEventListener('change', autoSaveSettings);
    document.getElementById('themeSelect')?.addEventListener('change', autoSaveSettings);
    document.getElementById('retentionSelect')?.addEventListener('change', autoSaveSettings);

    // Theme selector - apply immediately
    document.getElementById('themeSelect')?.addEventListener('change', (e) => {
        applyTheme(e.target.value);
    });

    // Search functionality
    document.getElementById('searchInput')?.addEventListener('input', (e) => {
        filterUsageTable(e.target.value);
    });

    // Sorting functionality - use event delegation
    document.addEventListener('click', (e) => {
        const sortableHeader = e.target.closest('.sortable');
        if (sortableHeader) {
            const column = sortableHeader.dataset.column;
            const table = sortableHeader.closest('table');
            if (table.id === 'statsTable') {
                handleSort(column);
            } else if (table.id === 'usageTable') {
                handleUsageSort(column);
            }
        }
    });
}

// Switch between pages
function switchPage(page) {
    currentPage = page;

    // Update tabs
    document.querySelectorAll('.nav-tab').forEach(tab => {
        tab.classList.remove('active');
        if (tab.dataset.page === page) {
            tab.classList.add('active');
        }
    });

    // Update pages
    document.querySelectorAll('.page').forEach(p => {
        p.classList.remove('active');
    });
    document.getElementById(`${page}-page`)?.classList.add('active');

    // Load page-specific data
    if (page === 'usage') {
        loadNetworkUsage();
    } else if (page === 'settings') {
        loadDatabaseStats();
    }
}

// Start real-time data updates
function startDataUpdates() {
    updateDashboard();
    setInterval(updateDashboard, 500); // Update every 500ms for responsive real-time display
}

// Update dashboard with real-time stats
async function updateDashboard() {
    if (currentPage !== 'dashboard') return;

    try {
        // Get network stats - REAL-TIME data from Windows API
        // Data is collected via GetExtendedTcpTable, GetExtendedUdpTable, and GetIfTable2
        // Updates every second with live network traffic information
        const stats = await window.go.main.App.GetNetworkStats();
        updateStatsTable(stats || {});

        // Get today's stats - REAL data from SQLite database
        // Aggregated from actual network usage records stored in database
        const todayStats = await window.go.main.App.GetTodayStats();
        updateUsageSummary('today', todayStats);

    } catch (error) {
        console.error('Failed to update dashboard:', error);
        // Use cached data if available
        updateStatsTable(lastStatsCache);
    }
}

// Update stats table
function updateStatsTable(stats) {
    const tbody = document.getElementById('statsBody');
    if (!tbody) return;

    const statsArray = Object.entries(stats || {});

    if (statsArray.length === 0) {
        tbody.innerHTML = '<tr><td colspan="5" class="no-data">No active network connections</td></tr>';
        return;
    }

    // Cache stats
    lastStatsCache = stats;

    // Apply sorting if a column is selected
    if (currentSortColumn) {
        sortStatsArray(statsArray, currentSortColumn, currentSortDirection);
    }

    tbody.innerHTML = statsArray.map(([appName, stat]) => `
        <tr>
            <td>${escapeHtml(appName)}</td>
            <td class="upload-speed">${formatSpeed(stat.UploadSpeed || 0)}</td>
            <td class="download-speed">${formatSpeed(stat.DownloadSpeed || 0)}</td>
            <td class="total-upload">${formatBytes(stat.TotalUpload || 0)}</td>
            <td class="total-download">${formatBytes(stat.TotalDownload || 0)}</td>
        </tr>
    `).join('');

    // Update sort indicators
    updateSortIndicators();
}

// Handle sorting
function handleSort(column) {
    if (currentSortColumn === column) {
        // Toggle direction if same column
        currentSortDirection = currentSortDirection === 'asc' ? 'desc' : 'asc';
    } else {
        // New column - default to descending
        currentSortColumn = column;
        currentSortDirection = 'desc';
    }

    // Re-render with current stats
    updateStatsTable(lastStatsCache);
}

// Sort stats array based on column and direction
function sortStatsArray(statsArray, column, direction) {
    statsArray.sort((a, b) => {
        let valueA, valueB;

        switch (column) {
            case 'uploadSpeed':
                valueA = a[1].UploadSpeed || 0;
                valueB = b[1].UploadSpeed || 0;
                break;
            case 'downloadSpeed':
                valueA = a[1].DownloadSpeed || 0;
                valueB = b[1].DownloadSpeed || 0;
                break;
            case 'totalUpload':
                valueA = a[1].TotalUpload || 0;
                valueB = b[1].TotalUpload || 0;
                break;
            case 'totalDownload':
                valueA = a[1].TotalDownload || 0;
                valueB = b[1].TotalDownload || 0;
                break;
            default:
                return 0;
        }

        if (direction === 'asc') {
            return valueA - valueB;
        } else {
            return valueB - valueA;
        }
    });
}

// Update sort indicators in table headers
function updateSortIndicators() {
    const table = document.getElementById('statsTable');
    if (!table) return;

    const headers = table.querySelectorAll('.sortable');
    headers.forEach(header => {
        const indicator = header.querySelector('.sort-indicator');
        if (header.dataset.column === currentSortColumn) {
            indicator.textContent = currentSortDirection === 'asc' ? '▲' : '▼';
            header.classList.add('active-sort');
        } else {
            indicator.textContent = '';
            header.classList.remove('active-sort');
        }
    });
}

// Update usage summary cards
function updateUsageSummary(prefix, data) {
    const uploadEl = document.getElementById(`${prefix}Upload`);
    const downloadEl = document.getElementById(`${prefix}Download`);

    if (uploadEl) uploadEl.textContent = formatBytes(data?.upload || 0);
    if (downloadEl) downloadEl.textContent = formatBytes(data?.download || 0);
}

// Load network usage history
async function loadNetworkUsage() {
    try {
        const usageStats = await window.go.main.App.GetNetworkUsageStats();
        lastUsageCache = usageStats || [];
        updateUsageTable(lastUsageCache);
    } catch (error) {
        console.error('Failed to load network usage:', error);
    }
}

// Update usage table
function updateUsageTable(stats) {
    const tbody = document.getElementById('usageBody');
    if (!tbody) return;

    if (stats.length === 0) {
        tbody.innerHTML = '<tr><td colspan="5" class="no-data">No usage data available</td></tr>';
        return;
    }

    // Create a copy to sort
    const sortedStats = [...stats];

    // Apply sorting
    sortUsageArray(sortedStats, usageSortColumn, usageSortDirection);

    tbody.innerHTML = sortedStats.map(stat => `
        <tr>
            <td>${escapeHtml(stat.AppName)}</td>
            <td class="total-upload">${formatBytes(stat.TotalUpload || 0)}</td>
            <td class="total-download">${formatBytes(stat.TotalDownload || 0)}</td>
            <td>${formatBytes((stat.TotalUpload || 0) + (stat.TotalDownload || 0))}</td>
            <td>${formatTimestamp(stat.LastSeen)}</td>
        </tr>
    `).join('');

    // Update sort indicators
    updateUsageSortIndicators();
}

// Handle usage table sorting
function handleUsageSort(column) {
    if (usageSortColumn === column) {
        // Toggle direction if same column
        usageSortDirection = usageSortDirection === 'asc' ? 'desc' : 'asc';
    } else {
        // New column - default to descending
        usageSortColumn = column;
        usageSortDirection = 'desc';
    }

    // Re-render with current data
    updateUsageTable(lastUsageCache);
}

// Sort usage array based on column and direction
function sortUsageArray(statsArray, column, direction) {
    statsArray.sort((a, b) => {
        let valueA, valueB;

        switch (column) {
            case 'totalUpload':
                valueA = a.TotalUpload || 0;
                valueB = b.TotalUpload || 0;
                break;
            case 'totalDownload':
                valueA = a.TotalDownload || 0;
                valueB = b.TotalDownload || 0;
                break;
            case 'totalData':
                valueA = (a.TotalUpload || 0) + (a.TotalDownload || 0);
                valueB = (b.TotalUpload || 0) + (b.TotalDownload || 0);
                break;
            case 'lastSeen':
                valueA = a.LastSeen || 0;
                valueB = b.LastSeen || 0;
                break;
            default:
                return 0;
        }

        if (direction === 'asc') {
            return valueA - valueB;
        } else {
            return valueB - valueA;
        }
    });
}

// Update sort indicators in usage table headers
function updateUsageSortIndicators() {
    const table = document.getElementById('usageTable');
    if (!table) return;

    const headers = table.querySelectorAll('.sortable');
    headers.forEach(header => {
        const indicator = header.querySelector('.sort-indicator');
        if (header.dataset.column === usageSortColumn) {
            indicator.textContent = usageSortDirection === 'asc' ? '▲' : '▼';
            header.classList.add('active-sort');
        } else {
            indicator.textContent = '';
            header.classList.remove('active-sort');
        }
    });
}

// Filter usage table by search term
function filterUsageTable(searchTerm) {
    const tbody = document.getElementById('usageBody');
    if (!tbody) return;

    const rows = tbody.querySelectorAll('tr');
    const term = searchTerm.toLowerCase();

    rows.forEach(row => {
        const appName = row.cells[0]?.textContent.toLowerCase() || '';
        row.style.display = appName.includes(term) ? '' : 'none';
    });
}

// Pause monitoring
async function pauseMonitoring() {
    try {
        await window.go.main.App.PauseMonitoring();
        document.getElementById('pauseBtn').style.display = 'none';
        document.getElementById('resumeBtn').style.display = 'block';
        monitoringPaused = true;
    } catch (error) {
        console.error('Failed to pause monitoring:', error);
    }
}

// Resume monitoring
async function resumeMonitoring() {
    try {
        await window.go.main.App.ResumeMonitoring();
        document.getElementById('resumeBtn').style.display = 'none';
        document.getElementById('pauseBtn').style.display = 'block';
        monitoringPaused = false;
    } catch (error) {
        console.error('Failed to resume monitoring:', error);
    }
}

// Load settings
async function loadSettings() {
    try {
        const settings = await window.go.main.App.GetSettings();

        document.getElementById('autoStartCheck').checked = settings.AutoStart || false;
        document.getElementById('themeSelect').value = settings.Theme || 'auto';
        document.getElementById('retentionSelect').value = String(settings.DataRetention || 30);

        applyTheme(settings.Theme || 'auto');
    } catch (error) {
        console.error('Failed to load settings:', error);
    }
}

// Save settings
async function saveSettings() {
    try {
        const settings = {
            AutoStart: document.getElementById('autoStartCheck').checked,
            Theme: document.getElementById('themeSelect').value,
            DataRetention: parseInt(document.getElementById('retentionSelect').value),
            NetworkInterface: ''
        };

        await window.go.main.App.UpdateSettings(settings);

        applyTheme(settings.Theme);
    } catch (error) {
        console.error('Failed to save settings:', error);
    }
}

// Auto-save settings on change
async function autoSaveSettings() {
    try {
        const settings = {
            AutoStart: document.getElementById('autoStartCheck').checked,
            Theme: document.getElementById('themeSelect').value,
            DataRetention: parseInt(document.getElementById('retentionSelect').value),
            NetworkInterface: ''
        };

        await window.go.main.App.UpdateSettings(settings);

        applyTheme(settings.Theme);
    } catch (error) {
        console.error('Failed to save settings:', error);
        // Revert checkbox if autostart failed
        if (error.message && error.message.includes('autostart')) {
            const checkbox = document.getElementById('autoStartCheck');
            checkbox.checked = !checkbox.checked;
            alert('Failed to change autostart setting. Please try running as administrator.');
        }
    }
}

// Apply theme
function applyTheme(theme) {
    if (theme === 'light') {
        document.body.classList.add('light-theme');
    } else if (theme === 'dark') {
        document.body.classList.remove('light-theme');
    } else {
        // Auto theme - detect system preference
        const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
        if (prefersDark) {
            document.body.classList.remove('light-theme');
        } else {
            document.body.classList.add('light-theme');
        }
    }
}

// Load database statistics
async function loadDatabaseStats() {
    try {
        const stats = await window.go.main.App.GetDatabaseStats();

        document.getElementById('dbSize').textContent = formatBytes(stats.size || 0);
        document.getElementById('dbRecords').textContent = (stats.records || 0).toLocaleString();
        document.getElementById('dbOldest').textContent = stats.oldest ? formatTimestamp(stats.oldest) : 'Never';
    } catch (error) {
        console.error('Failed to load database stats:', error);
    }
}

// Cleanup old data
async function cleanupOldData() {
    try {
        await window.go.main.App.ClearOldData();
        // Refresh all relevant UI components
        loadDatabaseStats();
        loadDashboardData(); // Reset Today's Usage display
        loadUsageHistory(); // Clear the usage history table
    } catch (error) {
        console.error('Failed to cleanup data:', error);
    }
}

// Format bytes to human-readable format
function formatBytes(bytes) {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
}

// Format speed (bytes per second)
function formatSpeed(bytesPerSecond) {
    return formatBytes(bytesPerSecond) + '/s';
}

// Format Unix timestamp
function formatTimestamp(timestamp) {
    if (!timestamp) return 'Never';
    const date = new Date(timestamp * 1000);
    return date.toLocaleString();
}

// Escape HTML to prevent XSS
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}
