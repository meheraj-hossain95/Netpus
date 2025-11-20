@echo off
echo ========================================
echo Running Octopus Network Monitor (Dev)
echo ========================================
echo.

REM Check if Go is installed
go version >nul 2>nul
if %errorlevel% neq 0 (
    echo Error: Go is not installed or not in PATH
    echo Please install Go 1.22 or higher from https://go.dev/dl/
    pause
    exit /b 1
)

REM Check if Wails is installed
wails version >nul 2>nul
if %errorlevel% neq 0 (
    echo Error: Wails CLI is not installed
    echo Installing Wails CLI...
    go install github.com/wailsapp/wails/v2/cmd/wails@latest
    if %errorlevel% neq 0 (
        echo Failed to install Wails CLI
        pause
        exit /b 1
    )
    echo Wails CLI installed successfully
    echo.
)

echo Downloading dependencies...
go mod download
if %errorlevel% neq 0 (
    echo Failed to download dependencies
    pause
    exit /b 1
)
echo.

echo Starting application in development mode...
echo Press Ctrl+C to stop
echo.
wails dev
if %errorlevel% neq 0 (
    echo.
    echo Failed to start application!
    pause
    exit /b 1
)
