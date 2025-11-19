# M-Net (Octopus) Network Monitor# Octopus Network Monitor



**Real-time network bandwidth monitoring application for Windows and Linux****Real-time network bandwidth monitoring for Windows and Linux**



A lightweight desktop application that tracks network usage per application in real-time, built with Wails v2 (Go backend + JavaScript frontend).A lightweight desktop application that tracks network usage per application in real-time, built with Wails v2 (Go + JavaScript frontend).



------



## ✨ Features## ✨ Features



- **Real-Time Monitoring** – Track upload and download speeds for all applications with 1-second precision- **Real-Time Monitoring** – Track upload and download speeds for all applications with 1-second precision

- **Historical Data** – Analyze past network usage with configurable data retention (7/30/90 days or forever)- **Historical Data** – Analyze past network usage with configurable data retention (7/30/90 days or forever)

- **System Tray Integration** – Minimize to tray and monitor in the background with live speed tooltips- **System Tray Integration** – Minimize to tray and monitor in the background with live speed tooltips

- **Single Instance Protection** – Only one copy can run at a time to prevent conflicts- **Cross-Platform** – Native builds for Windows and Linux

- **Cross-Platform** – Native builds for Windows and Linux- **Privacy First** – Completely offline, no telemetry, all data stored locally

- **Privacy First** – Completely offline, no telemetry, all data stored locally- **Customizable** – Dark/Light/Auto themes, auto-start on boot, and flexible settings

- **Customizable** – Dark/Light/Auto themes, auto-start on boot, and flexible settings- **Lightweight** – Minimal resource usage (~50MB RAM, <1% CPU)

- **Lightweight** – Minimal resource usage (~50MB RAM, <1% CPU)

---

---

## 📸 Screenshots

## 📦 Installation

> **Note:** Screenshots will be added here in future releases.

### Windows Installation

---

#### Option 1: Quick Start (Portable - No Installation Required)

## 🚀 Quick Start

1. **Download** the latest `octopus.exe` from the [Releases](../../releases) page

2. **Run** the executable by double-clicking `octopus.exe`### Option 1: Download Pre-built Binary (Recommended)

3. **Done!** The app runs without any installation

#### Windows

**This is the simplest way to use M-Net.** The application works immediately without installing anything.1. **Download** the latest `octopus.exe` from [Releases](../../releases)

2. **Run** the executable by double-clicking it

#### Option 2: Install with Shortcuts (Optional)3. *Optional:* Run from Command Prompt to install shortcuts:

   ```cmd

If you want shortcuts on your Desktop and Start Menu:   octopus.exe --install

   ```

1. **Download** `octopus.exe` and place it in a permanent folder (e.g., `C:\Program Files\M-Net\`)

2. **Open Command Prompt** in that folder:#### Linux

   - Press `Win + R`1. **Download** the latest `octopus` binary from [Releases](../../releases)

   - Type `cmd` and press Enter2. **Make executable:**

   - Navigate to the folder: `cd "C:\Program Files\M-Net"`   ```bash

3. **Run the install command:**   chmod +x octopus

   ```cmd   ```

   octopus.exe --install3. **Run** the application:

   ```   ```bash

4. **Done!** You'll now have:   ./octopus

   - Desktop shortcut   ```

   - Start Menu entry4. *Optional:* Install to application menu:

   - Easy access from anywhere   ```bash

   ./octopus --install

**What --install does:**   ```

- Creates a shortcut on your Desktop

- Adds M-Net to your Start Menu---

- Does NOT modify system files or registry (except for shortcuts)

- The app still runs as a portable application### Option 2: Build from Source



---#### Prerequisites



### Linux Installation- **Go** 1.22 or higher ([Download](https://go.dev/dl/))

- **Wails CLI** v2 (auto-installed by build scripts if missing)

#### Option 1: Quick Start (Portable - No Installation Required)- **Node.js** (optional, only needed for frontend development)



1. **Download** the latest `octopus` binary from the [Releases](../../releases) page#### Build Instructions

2. **Make it executable:**

   ```bash##### Windows

   chmod +x octopus

   ```1. **Open Command Prompt or PowerShell** in the project directory:

3. **Run** the application:   ```cmd

   ```bash   cd C:\path\to\Octopus

   ./octopus   ```

   ```

4. **Done!** The app runs without any installation2. **Run the build script:**

   ```cmd

**This is the simplest way to use M-Net.** The application works immediately without installing anything.   build.bat

   ```

#### Option 2: Install to Application Menu (Optional)   or in PowerShell:

   ```powershell

If you want M-Net in your application launcher:   .\build.bat

   ```

1. **Move the binary** to a permanent location:

   ```bash3. **Wait for build completion** (takes 3-5 seconds)

   sudo mkdir -p /opt/m-net

   sudo mv octopus /opt/m-net/4. **Run the app:**

   ```   ```cmd

2. **Run the install command:**   .\build\bin\octopus.exe

   ```bash   ```

   /opt/m-net/octopus --install

   ```##### Linux

3. **Done!** You'll now have:

   - Application menu entry (in your app launcher)1. **Open Terminal** in the project directory:

   - Easy access from your desktop environment   ```bash

   cd /path/to/Octopus

**What --install does:**   ```

- Creates a `.desktop` file in `~/.local/share/applications/`

- Adds M-Net to your application menu2. **Make the build script executable** (first time only):

- Does NOT require root access   ```bash

- The app still runs as a portable application   chmod +x build.sh

   ```

---

3. **Run the build script:**

### Building from Source   ```bash

   ./build.sh

If you want to build M-Net yourself:   ```



#### Prerequisites4. **Run the app:**

   ```bash

- **Go** 1.22 or higher ([Download](https://go.dev/dl/))   ./build/bin/octopus

- **Wails CLI** v2 (auto-installed by build scripts if missing)   ```

- **Node.js** (optional, only for frontend development)

#### Rebuilding After Code Changes

#### Windows Build

Simply re-run the build script:

1. **Open Command Prompt** in the project folder:- **Windows:** `build.bat`

   ```cmd- **Linux:** `./build.sh`

   cd C:\path\to\M-Net

   ```The scripts automatically detect changes and rebuild everything.

2. **Run the build script:**

   ```cmd**Build output location:**

   build.bat- Windows: `build\bin\octopus.exe`

   ```- Linux: `build/bin/octopus`

3. **Find the executable:**

   ```---

   build\bin\octopus.exe

   ```## 📖 Usage



#### Linux Build### Application Interface



1. **Open Terminal** in the project folder:The application has three main pages:

   ```bash

   cd /path/to/M-Net1. **Dashboard** – Real-time network statistics with sortable table of active applications

   ```2. **Usage History** – Historical data analysis with search and filtering

2. **Make the script executable** (first time only):3. **Settings** – Configure auto-start, theme, and data retention policies

   ```bash

   chmod +x build.sh### First Run

   ```

3. **Run the build script:**On the first run, the application will:

   ```bash- Create a database in your user data folder

   ./build.sh- Start monitoring network traffic

   ```- Minimize to the system tray (look for the icon in your taskbar)

4. **Find the executable:**

   ```### Command Line Options

   build/bin/octopus

   ```#### Windows

```cmd

---octopus.exe                # Run the application (GUI mode)

octopus.exe --install      # Install shortcuts (Desktop + Start Menu)

## 🚀 Usageoctopus.exe --uninstall    # Remove all shortcuts and installations

octopus.exe --version      # Show version information

### First Run```



When you run M-Net for the first time:#### Linux

```bash

1. The application window will open./octopus                  # Run the application (GUI mode)

2. Network monitoring starts automatically./octopus --install        # Install to application menu

3. Data is saved to a local database in your user data folder./octopus --uninstall      # Remove application menu entry

4. The app minimizes to system tray when you close the window./octopus --version        # Show version information

```

**Important:** Only **ONE instance** of M-Net can run at a time. If you try to run it twice, you'll see an error message.

**Note:** The `--install` flag creates shortcuts for easy access but is optional. The application will work without installation.

### Command Line Options

### Configuration Settings

```cmd

# Windows- **Auto-start on boot** – Launch Octopus automatically when your system starts

octopus.exe              # Start the application- **Theme** – Choose between Auto (system), Light, or Dark mode

octopus.exe --install    # Install shortcuts (Desktop + Start Menu)- **Data Retention** – Keep data for 7, 30, 90 days, or forever

octopus.exe --uninstall  # Remove shortcuts- **Database Management** – View database size and manually clear old data

octopus.exe --version    # Show version information

---

# Linux

./octopus                # Start the application## 🗑️ How to Uninstall

./octopus --install      # Install to application menu

./octopus --uninstall    # Remove from application menu### Complete Uninstallation

./octopus --version      # Show version information

```#### Windows



### Application Interface**Method 1: Using the built-in uninstaller (Recommended)**

1. **Navigate to the application folder:**

**Dashboard** – Real-time network statistics with live upload/download speeds   ```cmd

   cd C:\path\to\Octopus\build\bin

**Usage History** – Historical data analysis with sortable tables   ```



**Settings** – Configure:2. **Run the uninstall command:**

- Auto-start on boot   ```cmd

- Theme (Auto/Light/Dark)   octopus.exe --uninstall

- Data retention policy (7/30/90 days or forever)   ```

- Database management   This removes all shortcuts from Desktop, Start Menu, and disables auto-start.



### System Tray3. **Delete the application folder:**

   - Go to the folder where you installed Octopus

When minimized, M-Net stays in your system tray with:   - Delete the entire `Octopus` folder

- Live upload/download speed tooltip

- Quick access menu4. **Delete app data** (optional - removes all network usage history):

- Pause/Resume monitoring   - Press `Win + R`

- Show/Hide window   - Type: `%APPDATA%\octopus`

- Quit application   - Delete the entire `octopus` folder



---**Method 2: Manual removal**

1. **Stop any running instances:**

## 🗑️ Uninstallation Guide   - Open Task Manager (`Ctrl+Shift+Esc`)

   - Find and end any `octopus.exe` processes

### Complete Removal (Windows)

2. **Delete shortcuts:**

#### Step 1: Remove Shortcuts (If Installed)   - Desktop: Delete `Octopus.lnk` if present

   - Start Menu: Navigate to `%APPDATA%\Microsoft\Windows\Start Menu\Programs` and delete `Octopus.lnk`

If you ran `--install`, remove the shortcuts first:

3. **Disable auto-start:**

```cmd   - Press `Win + R` and type `shell:startup`

# Open Command Prompt in the M-Net folder   - Delete any `Octopus` shortcuts

cd "C:\Program Files\M-Net"

4. **Delete application folder and data:**

# Run uninstall command   - Delete the application folder

octopus.exe --uninstall   - Delete `%APPDATA%\octopus` folder

```

#### Linux

**This removes:**

- Desktop shortcut**Method 1: Using the built-in uninstaller (Recommended)**

- Start Menu entry1. **Navigate to the application folder:**

- Auto-start configuration   ```bash

   cd /path/to/Octopus/build/bin

#### Step 2: Stop the Application   ```



1. **Right-click** the M-Net icon in your system tray2. **Run the uninstall command:**

2. **Click** "Quit"   ```bash

   ./octopus --uninstall

OR   ```

   This removes the application menu entry and autostart configuration.

1. **Open Task Manager** (Ctrl+Shift+Esc)

2. **Find** `octopus.exe`3. **Delete the application folder:**

3. **End Task**   ```bash

   cd .. && cd ..

#### Step 3: Delete the Application   rm -rf Octopus

   ```

Simply delete the folder where you placed `octopus.exe`:

4. **Delete app data** (optional - removes all network usage history):

```cmd   ```bash

# Example if you placed it in C:\Program Files\M-Net   rm -rf ~/.local/share/octopus

rmdir /s "C:\Program Files\M-Net"   ```

```

**Method 2: Manual removal**

OR just delete the folder in File Explorer.1. **Stop any running instances:**

   ```bash

#### Step 4: Delete User Data (Optional)   pkill -f octopus

   ```

If you want to remove all your network usage history:

2. **Delete application menu entry:**

1. Press `Win + R`   ```bash

2. Type: `%APPDATA%\octopus`   rm ~/.local/share/applications/octopus.desktop

3. Delete the entire `octopus` folder   ```



**Full path:** `C:\Users\YourName\AppData\Roaming\octopus\`3. **Disable auto-start:**

   ```bash

**This contains:**   rm ~/.config/autostart/octopus.desktop

- `octopus.db` – Your network usage database   ```

- `config.json` – Application settings

4. **Delete application folder and data:**

**Done!** M-Net is completely removed from your system.   ```bash

   rm -rf /path/to/Octopus

---   rm -rf ~/.local/share/octopus

   ```

### Complete Removal (Linux)

### What Gets Removed

#### Step 1: Remove Application Menu Entry (If Installed)

- ✅ **Desktop shortcuts** (Windows only)

If you ran `--install`, remove it first:- ✅ **Start Menu entries** (Windows) or Application Menu (Linux)

- ✅ **Auto-start configuration**

```bash- ✅ **Application executable** (manual deletion)

# Navigate to where you installed M-Net- ✅ **Database and settings** (optional, manual deletion)

cd /opt/m-net

**Note:** The `--uninstall` command removes shortcuts and auto-start but does NOT delete the application executable or your data. You must manually delete those if desired.

# Run uninstall command

./octopus --uninstall---

```

## 💾 Data Storage

**This removes:**

- Application menu entry (`.desktop` file)All data is stored locally on your machine:

- Auto-start configuration

### Windows

#### Step 2: Stop the Application- **Database:** `%APPDATA%\octopus\octopus.db`

- **Full path example:** `C:\Users\YourName\AppData\Roaming\octopus\octopus.db`

```bash

# Find the process### Linux

ps aux | grep octopus- **Database:** `~/.local/share/octopus/octopus.db`

- **Full path example:** `/home/username/.local/share/octopus/octopus.db`

# Kill the process (replace XXXX with the PID)

kill XXXXNo data is ever transmitted to external servers. Everything stays on your computer.

```

---

OR right-click the tray icon and select "Quit".

## 🏗️ Architecture

#### Step 3: Delete the Application

M-Net is built with:

```bash

# If you installed to /opt/m-net- **Backend:** Go 1.22+ with Wails v2 framework

sudo rm -rf /opt/m-net- **Frontend:** HTML, CSS, and Vanilla JavaScript (no framework dependencies)

- **Database:** SQLite with WAL mode for concurrent access

# If you kept it elsewhere, delete that folder- **Platform APIs:**

rm -rf /path/to/m-net  - Windows: Windows API (GetExtendedTcpTable, GetIfTable2, etc.)

```  - Linux: `/proc` filesystem parsing



#### Step 4: Delete User Data (Optional)For detailed technical documentation, see [APP_DOCUMENTATION.md](APP_DOCUMENTATION.md).



If you want to remove all your network usage history:### Project Structure



```bash```

rm -rf ~/.local/share/octopusOctopus/

```├── main.go                 # Application entry point

├── app.go                  # Application logic layer

**This contains:**├── go.mod                  # Go module definition

- `octopus.db` – Your network usage database├── wails.json              # Wails configuration

- `config.json` – Application settings├── build.bat               # Windows build script

├── build.sh                # Linux build script

**Done!** M-Net is completely removed from your system.├── APP_DOCUMENTATION.md    # Detailed architecture documentation

├── frontend/               # Frontend assets

---│   ├── index.html

│   ├── src/

### Quick Uninstall Commands│   │   ├── main.js

│   │   └── style.css

#### Windows (One Command)│   └── wailsjs/           # Generated Wails bindings

```cmd├── internal/              # Internal packages

# Run from M-Net folder│   ├── autostart/         # System autostart management

octopus.exe --uninstall && cd .. && rmdir /s M-Net && rmdir /s "%APPDATA%\octopus"│   ├── database/          # SQLite database operations

```│   ├── installer/         # Installation utilities

│   ├── monitor/           # Network monitoring

#### Linux (One Command)│   ├── tray/              # System tray integration

```bash│   └── utils/             # Utility functions

# If installed to /opt/m-net└── build/

./octopus --uninstall && sudo rm -rf /opt/m-net && rm -rf ~/.local/share/octopus    └── bin/               # Compiled executables

``````



------



## 💾 Data Storage Locations## 📋 Requirements



All data is stored locally on your computer:### Windows

- Windows 10 or later

### Windows- WebView2 Runtime (usually pre-installed)

- **Database:** `%APPDATA%\octopus\octopus.db`

- **Full path:** `C:\Users\YourName\AppData\Roaming\octopus\octopus.db`### Linux

- X11 or Wayland desktop environment

### Linux- System tray support (most DEs include this)

- **Database:** `~/.local/share/octopus/octopus.db`- GTK 3 (usually pre-installed)

- **Full path:** `/home/username/.local/share/octopus/octopus.db`

---

**No data is ever sent to any external servers.** Everything stays on your machine.

## 🐛 Troubleshooting

---

### "Only one instance can run at a time"

## 🔧 TroubleshootingOctopus prevents multiple instances. If you see this error but no window appears:

- Check for the process in Task Manager (Windows) or `ps aux | grep octopus` (Linux)

### "Octopus is already running" Error- Kill the process if stuck:

  - Windows: `taskkill /F /IM octopus.exe`

**Problem:** You see this message when trying to start M-Net, but no window appears.  - Linux: `pkill -f octopus`

- Delete the lock file: 

**Cause:** A previous instance is still running or didn't close properly.  - Windows: `%TEMP%\octopus.lock`

  - Linux: `/tmp/octopus.lock`

**Solution:**

### Database growing too large

#### Windows- Adjust data retention in Settings

1. **Check Task Manager:**- Use "Clear Old Data" button to manually clean up

   - Press `Ctrl+Shift+Esc`- Consider reducing retention period

   - Find `octopus.exe`

   - Click "End Task"### Missing system tray icon (Linux)

- Ensure your desktop environment supports system tray

2. **Delete the lock file:**- Try installing `libappindicator3-1` (Debian/Ubuntu) or equivalent

   ```cmd

   del "%TEMP%\octopus.lock"### Build fails - Go not found

   ```- Ensure Go 1.22+ is installed and in your PATH

- Windows: Add Go to PATH in System Environment Variables

3. **Try running again**- Linux: Add `export PATH=$PATH:/usr/local/go/bin` to `~/.bashrc`



#### Linux---

1. **Kill the process:**

   ```bash## 🔒 Privacy & Security

   pkill -f octopus

   ```- **Completely Offline** – No network connections, no telemetry, no tracking

- **Local Data Only** – All data stored on your machine with standard file permissions

2. **Delete the lock file:**- **Read-Only Access** – Only reads network statistics, does not modify or block traffic

   ```bash- **Open Source** – Full source code available for inspection

   rm /tmp/octopus.lock

   ```---



3. **Try running again**## ⚡ Performance



---- **Fixed Update Interval**: 1 second for real-time monitoring

- **Batch Writes**: Database writes every 10 seconds

### Database Growing Too Large- **Memory Efficient**: Inactive processes cleaned up after 30 seconds

- **Database Optimization**: Automatic vacuum and cleanup

**Problem:** The database file is taking up too much space.

---

**Solution:**

## 🤝 Contributing

1. Open M-Net Settings

2. Adjust **Data Retention** to a shorter period (e.g., 30 days instead of 90)Contributions are welcome! Whether it's bug reports, feature requests, or code contributions:

3. Click **"Clear Old Data"** to manually remove old records

4. The database will automatically vacuum and reclaim space1. **Open an issue** to discuss your idea or report a bug

2. **Fork the repository** and create a pull request

---3. **Submit improvements** to documentation or code



### System Tray Icon Missing (Linux)For detailed technical information, see [APP_DOCUMENTATION.md](APP_DOCUMENTATION.md).



**Problem:** The tray icon doesn't appear in your system tray.---



**Cause:** Your desktop environment might not support system tray icons.## 📄 License



**Solution:**MIT License



- **Ubuntu/GNOME:** Install AppIndicator extensionCopyright (c) 2025 Octopus Network Monitor

  ```bash

  sudo apt install gnome-shell-extension-appindicatorPermission is hereby granted, free of charge, to any person obtaining a copy

  ```of this software and associated documentation files (the "Software"), to deal

- **KDE Plasma:** System tray is built-in, check if it's enabled in System Settingsin the Software without restriction, including without limitation the rights

- **XFCE:** System tray is built-in by defaultto use, copy, modify, merge, publish, distribute, sublicense, and/or sell

copies of the Software, and to permit persons to whom the Software is

---furnished to do so, subject to the following conditions:



### Cannot Build from SourceThe above copyright notice and this permission notice shall be included in all

copies or substantial portions of the Software.

**Problem:** Build fails with "Go not found" or similar error.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR

**Solution:**IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,

FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE

1. **Install Go 1.22 or higher:**AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER

   - Windows: Download from [go.dev](https://go.dev/dl/)LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,

   - Linux: `sudo apt install golang-go` or download from [go.dev](https://go.dev/dl/)OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE

SOFTWARE.

2. **Verify Go is in PATH:**

   ```cmd---

   # Windows

   go version## 🙏 Acknowledgments



   # LinuxBuilt with:

   go version- [Wails](https://wails.io/) – Go + Web framework for desktop apps

   ```- [systray](https://github.com/getlantern/systray) – Cross-platform system tray library

- [SQLite](https://www.sqlite.org/) – Embedded database engine

3. **Try building again**

---

---

**Made with ❤️ for network monitoring enthusiasts**

## 🔒 Privacy & Security

- ✅ **Completely Offline** – No network connections, no telemetry, no tracking
- ✅ **Local Data Only** – All data stored on your machine
- ✅ **Read-Only Access** – Only reads network statistics, never modifies traffic
- ✅ **Open Source** – Full source code available for inspection
- ✅ **No Permissions Required** – Runs in user space (except on Linux where it may need capabilities for network monitoring)

---

## 📋 System Requirements

### Windows
- Windows 10 or later (64-bit)
- WebView2 Runtime (usually pre-installed on Windows 10/11)
- ~50MB disk space
- ~50MB RAM during operation

### Linux
- Any modern Linux distribution (64-bit)
- X11 or Wayland display server
- System tray support (GNOME, KDE, XFCE, etc.)
- GTK 3+ (usually pre-installed)
- ~50MB disk space
- ~50MB RAM during operation

---

## 🏗️ Project Structure

```
M-Net/
├── main.go                 # Application entry point
├── app.go                  # Application logic and API
├── go.mod                  # Go dependencies
├── wails.json              # Wails configuration
├── build.bat               # Windows build script
├── build.sh                # Linux build script
├── README.md               # This file
├── APP_DOCUMENTATION.md    # Technical documentation
├── frontend/               # Web-based UI
│   ├── index.html
│   ├── src/
│   │   ├── main.js
│   │   └── style.css
│   └── wailsjs/           # Generated Wails bindings
├── internal/              # Internal Go packages
│   ├── autostart/         # System autostart management
│   ├── database/          # SQLite database operations
│   ├── installer/         # Installation utilities
│   ├── monitor/           # Network monitoring engine
│   ├── tray/              # System tray integration
│   └── utils/             # Helper functions
└── build/
    └── bin/               # Compiled executables
        └── octopus.exe    # (Windows) or octopus (Linux)
```

---

## 🤝 Contributing

Contributions are welcome! Whether it's:

- 🐛 Bug reports
- ✨ Feature requests
- 📝 Documentation improvements
- 💻 Code contributions

**How to contribute:**

1. Open an issue to discuss your idea
2. Fork the repository
3. Make your changes
4. Submit a pull request

For technical details, see [APP_DOCUMENTATION.md](APP_DOCUMENTATION.md).

---

## 📄 License

MIT License - See [LICENSE](LICENSE) file for details.

Copyright (c) 2025 M-Net Contributors

---

## 🙏 Acknowledgments

Built with:
- [Wails v2](https://wails.io/) – Go + Web framework for desktop apps
- [systray](https://github.com/getlantern/systray) – Cross-platform system tray library
- [SQLite](https://www.sqlite.org/) – Embedded database engine

---

**Made with ❤️ for privacy-conscious network monitoring**

Need help? Open an issue on GitHub!
