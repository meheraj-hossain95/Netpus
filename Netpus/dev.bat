@echo off
echo ========================================
echo Starting Octopus in Development Mode
echo ========================================
echo.
echo This will start the app with:
echo - Hot reload enabled
echo - Debug console
echo - Development tools
echo.

wails dev

if %errorlevel% neq 0 (
    echo.
    echo Failed to start development mode!
    pause
    exit /b 1
)
