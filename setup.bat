@echo off
echo ========================================
echo Setting Up Octopus Development Environment
echo ========================================
echo.

REM Check if Go is installed
echo [1/4] Checking Go installation...
go version >nul 2>nul
if %errorlevel% neq 0 (
    echo Error: Go is not installed or not in PATH
    echo Please install Go 1.22 or higher from https://go.dev/dl/
    pause
    exit /b 1
)
go version
echo.

REM Check if Wails is installed
echo [2/4] Checking Wails CLI...
wails version >nul 2>nul
if %errorlevel% neq 0 (
    echo Wails CLI not found. Installing...
    go install github.com/wailsapp/wails/v2/cmd/wails@latest
    if %errorlevel% neq 0 (
        echo Failed to install Wails CLI
        pause
        exit /b 1
    )
    echo.
    echo Wails CLI installed successfully!
    echo.
    echo IMPORTANT: Please restart your terminal or add the following to your PATH:
    echo %USERPROFILE%\go\bin
    echo.
    pause
) else (
    wails version
)
echo.

REM Download Go dependencies
echo [3/4] Downloading Go dependencies...
go mod download
if %errorlevel% neq 0 (
    echo Failed to download dependencies
    pause
    exit /b 1
)
echo.

REM Verify Wails doctor
echo [4/4] Running Wails doctor to verify setup...
wails doctor
echo.

echo ========================================
echo Setup Complete!
echo ========================================
echo.
echo You can now:
echo - Run in dev mode:    run.bat or dev.bat
echo - Build for release:  build.bat
echo - Clean artifacts:    clean.bat
echo.
pause
