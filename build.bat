@echo off
setlocal enabledelayedexpansion

echo ========================================
echo Building Octopus Network Monitor
echo ========================================
echo.

REM Parse command line arguments
set "BUILD_MODE=production"
set "SKIP_CLEAN=0"
set "VERBOSE=0"

:parse_args
if "%~1"=="" goto end_parse
if /i "%~1"=="--debug" set "BUILD_MODE=debug"
if /i "%~1"=="--dev" set "BUILD_MODE=dev"
if /i "%~1"=="--skip-clean" set "SKIP_CLEAN=1"
if /i "%~1"=="--verbose" set "VERBOSE=1"
if /i "%~1"=="-v" set "VERBOSE=1"
if /i "%~1"=="--help" goto show_help
if /i "%~1"=="-h" goto show_help
shift
goto parse_args

:show_help
echo Usage: build.bat [OPTIONS]
echo.
echo Options:
echo   --debug         Build with debug symbols (no optimizations)
echo   --dev           Build in development mode
echo   --skip-clean    Skip cleaning before build
echo   --verbose, -v   Show detailed build output
echo   --help, -h      Show this help message
echo.
echo Examples:
echo   build.bat                    Build production release
echo   build.bat --debug            Build with debug symbols
echo   build.bat --dev --verbose    Build dev mode with verbose output
echo.
exit /b 0

:end_parse

REM Check if Go is installed
echo [1/5] Checking Go installation...
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
echo [2/5] Checking Wails CLI...
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
wails version
echo.

REM Clean previous build
if "%SKIP_CLEAN%"=="0" (
    echo [3/5] Cleaning previous build...
    if exist "build\bin" (
        rmdir /s /q "build\bin"
    )
    if exist "frontend\dist" (
        rmdir /s /q "frontend\dist"
    )
    echo Clean complete.
    echo.
) else (
    echo [3/5] Skipping clean step...
    echo.
)

REM Download dependencies
echo [4/5] Downloading dependencies...
go mod download
if %errorlevel% neq 0 (
    echo Failed to download dependencies
    pause
    exit /b 1
)
go mod tidy
echo.

REM Build application
echo [5/5] Building application in %BUILD_MODE% mode...
echo.

set "BUILD_FLAGS="
if "%BUILD_MODE%"=="debug" (
    echo Building with debug symbols...
    set "BUILD_FLAGS=-debug"
)
if "%BUILD_MODE%"=="dev" (
    echo Building in development mode...
    set "BUILD_FLAGS=-devbuild"
)

if "%VERBOSE%"=="1" (
    set "BUILD_FLAGS=%BUILD_FLAGS% -v 2"
)

echo Build command: wails build %BUILD_FLAGS%
echo.

wails build %BUILD_FLAGS%
if %errorlevel% neq 0 (
    echo.
    echo ========================================
    echo Build failed!
    echo ========================================
    pause
    exit /b 1
)

echo.
echo ========================================
echo Build successful!
echo ========================================
echo.
echo Build mode:     %BUILD_MODE%
echo Executable:     build\bin\octopus.exe
echo.

REM Check file size
if exist "build\bin\octopus.exe" (
    for %%A in ("build\bin\octopus.exe") do (
        set size=%%~zA
        set /a sizeMB=!size! / 1048576
        echo File size:      !sizeMB! MB
    )
)
echo.
echo To run the app:
echo   build\bin\octopus.exe
echo.
echo To install shortcuts:
echo   build\bin\octopus.exe --install
echo.
pause
