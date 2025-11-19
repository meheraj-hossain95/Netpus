# Octopus Network Monitor - Build Guide

## 📋 Table of Contents
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Build Scripts](#build-scripts)
- [Build Modes](#build-modes)
- [Development Workflow](#development-workflow)
- [Troubleshooting](#troubleshooting)

---

## 🔧 Prerequisites

### Required Software
1. **Go 1.22+** - [Download](https://go.dev/dl/)
2. **Wails CLI v2** - Installed automatically via setup script
3. **Git** (optional) - For version control

### Windows Requirements
- Windows 10/11 (64-bit)
- WebView2 Runtime (usually pre-installed on Windows 11)
- Command Prompt or PowerShell

---

## 🚀 Quick Start

### First-Time Setup
```batch
# Run the setup script to install dependencies
setup.bat
```

This will:
- Verify Go installation
- Install Wails CLI
- Download Go dependencies
- Run system verification

### Build for Production
```batch
# Build optimized production executable
build.bat
```

Output: `build\bin\octopus.exe`

### Run in Development Mode
```batch
# Start with hot-reload
run.bat
```

---

## 📜 Build Scripts

### `setup.bat`
**Purpose:** Initial environment setup and dependency installation

**Usage:**
```batch
setup.bat
```

**What it does:**
- ✅ Checks Go installation
- ✅ Installs Wails CLI
- ✅ Downloads Go modules
- ✅ Runs `wails doctor` for verification

---

### `build.bat`
**Purpose:** Build production-ready executable

**Usage:**
```batch
build.bat [OPTIONS]
```

**Options:**
| Option | Description |
|--------|-------------|
| `--debug` | Build with debug symbols (no optimizations) |
| `--dev` | Build in development mode |
| `--skip-clean` | Skip cleaning before build |
| `--verbose`, `-v` | Show detailed build output |
| `--help`, `-h` | Show help message |

**Examples:**
```batch
# Standard production build
build.bat

# Debug build with symbols
build.bat --debug

# Development build with verbose output
build.bat --dev --verbose

# Quick rebuild without cleaning
build.bat --skip-clean
```

**Build Steps:**
1. ✅ Verify Go installation
2. ✅ Verify Wails CLI
3. ✅ Clean previous builds (unless skipped)
4. ✅ Download dependencies
5. ✅ Build application
6. ✅ Display build information

---

### `run.bat`
**Purpose:** Run application in development mode

**Usage:**
```batch
run.bat
```

**Features:**
- 🔄 Hot reload on code changes
- 🐛 Debug console enabled
- 🛠️ Development tools available
- 📝 Live error reporting

**What it does:**
- Downloads dependencies
- Starts `wails dev` mode
- Opens application window
- Watches for file changes

---

### `dev.bat`
**Purpose:** Quick start for development (assumes setup is done)

**Usage:**
```batch
dev.bat
```

**Difference from run.bat:**
- Skips dependency checks
- Faster startup
- Assumes environment is ready

---

### `clean.bat`
**Purpose:** Remove all build artifacts and caches

**Usage:**
```batch
clean.bat
```

**What it removes:**
- `build\bin\` - Built executables
- `frontend\dist\` - Frontend build output
- `frontend\node_modules\` - Node dependencies (if exists)
- Go build cache
- Go module cache

**When to use:**
- Build issues or corruption
- Disk space cleanup
- Fresh build needed
- Switching branches

---

## 🏗️ Build Modes

### Production Mode (Default)
```batch
build.bat
```
- ✅ Full optimizations
- ✅ Smallest file size
- ✅ No debug symbols
- ✅ Release configuration
- **Use for:** Final releases, distribution

### Debug Mode
```batch
build.bat --debug
```
- ✅ Debug symbols included
- ✅ No optimizations
- ✅ Larger file size
- ✅ Better error messages
- **Use for:** Debugging issues, profiling

### Development Mode
```batch
build.bat --dev
```
- ✅ Fast compilation
- ✅ Basic optimizations
- ✅ Quick testing
- **Use for:** Testing builds without dev server

---

## 💻 Development Workflow

### Recommended Workflow

#### 1. Initial Setup (First Time Only)
```batch
setup.bat
```

#### 2. Daily Development
```batch
# Start development server
run.bat

# Make code changes (auto-reloads)
# Test in the app window
```

#### 3. Testing a Build
```batch
# Build in dev mode for quick testing
build.bat --dev

# Run the built executable
build\bin\octopus.exe
```

#### 4. Pre-Release
```batch
# Clean everything
clean.bat

# Build production version
build.bat

# Test the production build
build\bin\octopus.exe
```

#### 5. Release
```batch
# Final production build
build.bat --verbose

# Verify the executable
build\bin\octopus.exe --version

# Distribute octopus.exe
```

---

## 🐛 Troubleshooting

### Common Issues

#### 1. "Go is not installed or not in PATH"
**Solution:**
- Download Go from https://go.dev/dl/
- Install and restart terminal
- Verify: `go version`

#### 2. "Wails CLI is not installed"
**Solution:**
- Run `setup.bat` to auto-install
- OR manually: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- Add `%USERPROFILE%\go\bin` to PATH
- Restart terminal

#### 3. "Build failed" with compilation errors
**Solution:**
```batch
# Clean and rebuild
clean.bat
setup.bat
build.bat
```

#### 4. WebView2 Missing
**Solution:**
- Download from: https://developer.microsoft.com/microsoft-edge/webview2/
- Install WebView2 Runtime
- Restart application

#### 5. Port Already in Use (Dev Mode)
**Solution:**
- Close other Wails apps
- Kill process: `taskkill /F /IM octopus.exe`
- Try again: `run.bat`

#### 6. Frontend Build Issues
**Solution:**
```batch
# Clean frontend
rmdir /s /q frontend\dist

# Rebuild
build.bat
```

#### 7. Module Download Errors
**Solution:**
```batch
# Clear Go module cache
go clean -modcache

# Re-download
go mod download
go mod tidy
```

---

## 📊 Build Output

### File Structure After Build
```
build/
├── bin/
│   └── octopus.exe          # Main executable
└── windows/
    ├── info.json             # Build metadata
    └── wails.exe.manifest    # Windows manifest
```

### Typical File Sizes
- **Production:** ~15-25 MB
- **Debug:** ~30-40 MB
- **Development:** ~20-30 MB

---

## 🔍 Verification Commands

### Check Installation
```batch
# Go version
go version

# Wails version
wails version

# Project structure
wails doctor
```

### Build Information
```batch
# After building, check executable
build\bin\octopus.exe --version

# File size
dir build\bin\octopus.exe
```

---

## 📝 Build Configuration

### wails.json
The main configuration file:
```json
{
  "name": "octopus",
  "outputfilename": "octopus",
  "frontend:build": "node build.js",
  "frontend:dev:serverUrl": "auto",
  "outputType": "desktop"
}
```

**Key Settings:**
- `outputfilename` - Name of executable (octopus.exe)
- `frontend:build` - Frontend build command
- `outputType` - Application type (desktop)

---

## 🎯 Best Practices

### DO:
✅ Run `setup.bat` before first build
✅ Use `run.bat` during development
✅ Use `clean.bat` when switching branches
✅ Test production builds before release
✅ Keep Go and Wails updated

### DON'T:
❌ Edit `wailsjs` folder manually (auto-generated)
❌ Commit `build/bin` to version control
❌ Commit `frontend/dist` to version control
❌ Run multiple dev instances simultaneously

---

## 🚀 Performance Tips

### Faster Builds
1. Use `--skip-clean` for incremental builds
2. Close unnecessary applications
3. Exclude project from antivirus scanning
4. Use SSD for project directory

### Smaller Executables
1. Use production mode (default)
2. Strip debug symbols (automatic in production)
3. Use UPX compression (optional): `upx --best build\bin\octopus.exe`

---

## 📞 Support

### Resources
- [Wails Documentation](https://wails.io/docs/introduction)
- [Go Documentation](https://go.dev/doc/)
- Project Issues: See repository

### Debug Information
When reporting issues, include:
```batch
go version
wails version
wails doctor
```

---

## 📄 License
This build system is part of Octopus Network Monitor
Copyright © 2025
