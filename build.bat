@echo off
setlocal enabledelayedexpansion

REM Set PowerShell path
set "PS=%SystemRoot%\System32\WindowsPowerShell\v1.0\powershell.exe"

echo.
echo  ================================================
echo       NETPUS - One-Click Build ^& Install
echo  ================================================
echo.

REM =============================================
REM Step 1: Check Go Installation
REM =============================================
echo [1/5] Checking Go installation...
go version >nul 2>nul
if %errorlevel% neq 0 (
    echo.
    echo  ERROR: Go is not installed!
    echo.
    echo  Please install Go 1.22+ from: https://go.dev/dl/
    echo  After installing, restart your terminal and run this script again.
    echo.
    pause
    exit /b 1
)
echo       Go found.
echo.

REM =============================================
REM Step 2: Install/Check Wails CLI
REM =============================================
echo [2/5] Checking Wails CLI...
wails version >nul 2>nul
if %errorlevel% neq 0 (
    echo       Wails not found. Installing...
    go install github.com/wailsapp/wails/v2/cmd/wails@latest
    if %errorlevel% neq 0 (
        echo.
        echo  ERROR: Failed to install Wails CLI
        echo  Please check your internet connection and try again.
        echo.
        pause
        exit /b 1
    )
    echo       Wails CLI installed successfully!
) else (
    echo       Wails found.
)
echo.

REM =============================================
REM Step 3: Download Dependencies
REM =============================================
echo [3/5] Downloading dependencies...
go mod download >nul 2>nul
if %errorlevel% neq 0 (
    echo.
    echo  ERROR: Failed to download dependencies
    echo.
    pause
    exit /b 1
)
go mod tidy >nul 2>nul
echo       Dependencies ready.
echo.

REM =============================================
REM Step 4: Build Production Executable
REM =============================================
echo [4/5] Building Netpus.exe...
echo       This may take 1-2 minutes on first build...
echo.

wails build -clean -o Netpus.exe
if %errorlevel% neq 0 (
    echo.
    echo  ERROR: Build failed!
    echo  Please check the error messages above.
    echo.
    pause
    exit /b 1
)
echo       Build complete.
echo.

REM =============================================
REM Step 5: Install to Program Files & Create Shortcuts
REM =============================================
echo [5/5] Installing Netpus...

REM Create installation directory
set "INSTALL_DIR=%LOCALAPPDATA%\Netpus"
if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"

REM Copy executable to install directory
copy /Y "build\bin\Netpus.exe" "%INSTALL_DIR%\Netpus.exe" >nul
if %errorlevel% neq 0 (
    echo       Warning: Could not copy to install directory.
    echo       Using build directory instead.
    set "INSTALL_DIR=%~dp0build\bin"
)

REM Create Start Menu shortcut (makes it searchable in Windows)
set "STARTMENU=%APPDATA%\Microsoft\Windows\Start Menu\Programs"
echo       Creating Start Menu shortcut...

"%PS%" -NoProfile -Command "$WshShell = New-Object -ComObject WScript.Shell; $Shortcut = $WshShell.CreateShortcut('%STARTMENU%\Netpus.lnk'); $Shortcut.TargetPath = '%INSTALL_DIR%\Netpus.exe'; $Shortcut.WorkingDirectory = '%INSTALL_DIR%'; $Shortcut.Description = 'Network Bandwidth Monitor'; $Shortcut.Save()"

if %errorlevel% neq 0 (
    echo       Warning: Could not create Start Menu shortcut.
) else (
    echo       Start Menu shortcut created.
)

REM Create Desktop shortcut
set "DESKTOP=%USERPROFILE%\Desktop"
echo       Creating Desktop shortcut...

"%PS%" -NoProfile -Command "$WshShell = New-Object -ComObject WScript.Shell; $Shortcut = $WshShell.CreateShortcut('%DESKTOP%\Netpus.lnk'); $Shortcut.TargetPath = '%INSTALL_DIR%\Netpus.exe'; $Shortcut.WorkingDirectory = '%INSTALL_DIR%'; $Shortcut.Description = 'Network Bandwidth Monitor'; $Shortcut.Save()"

if %errorlevel% neq 0 (
    echo       Warning: Could not create Desktop shortcut.
) else (
    echo       Desktop shortcut created.
)

REM Register in Windows Apps & Features (Add/Remove Programs)
echo       Registering in Apps ^& Features...
set "UNINSTALL_KEY=HKCU\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\Netpus"
set "REG=%SystemRoot%\System32\reg.exe"
"%REG%" add "%UNINSTALL_KEY%" /v "DisplayName" /t REG_SZ /d "Netpus Network Monitor" /f >nul 2>nul
"%REG%" add "%UNINSTALL_KEY%" /v "DisplayVersion" /t REG_SZ /d "1.0.0" /f >nul 2>nul
"%REG%" add "%UNINSTALL_KEY%" /v "Publisher" /t REG_SZ /d "Netpus" /f >nul 2>nul
"%REG%" add "%UNINSTALL_KEY%" /v "InstallLocation" /t REG_SZ /d "%INSTALL_DIR%" /f >nul 2>nul
"%REG%" add "%UNINSTALL_KEY%" /v "UninstallString" /t REG_SZ /d "\"%INSTALL_DIR%\Netpus.exe\" --uninstall" /f >nul 2>nul
"%REG%" add "%UNINSTALL_KEY%" /v "DisplayIcon" /t REG_SZ /d "%INSTALL_DIR%\Netpus.exe" /f >nul 2>nul
"%REG%" add "%UNINSTALL_KEY%" /v "NoModify" /t REG_DWORD /d 1 /f >nul 2>nul
"%REG%" add "%UNINSTALL_KEY%" /v "NoRepair" /t REG_DWORD /d 1 /f >nul 2>nul
"%REG%" add "%UNINSTALL_KEY%" /v "EstimatedSize" /t REG_DWORD /d 15000 /f >nul 2>nul
echo       Registered in Apps ^& Features.

echo.
echo  ================================================
echo       INSTALLATION COMPLETE!
echo  ================================================
echo.
echo  Netpus has been installed and is ready to use!
echo.
echo  You can now:
echo    - Search "Netpus" in Windows Start Menu
echo    - Use the Desktop shortcut
echo    - Pin it to your Taskbar (right-click shortcut)
echo.
echo  Installed to: %INSTALL_DIR%
echo.
echo  ================================================
echo.

REM Auto-launch the app
echo Starting Netpus...
start "" "%INSTALL_DIR%\Netpus.exe"

echo.
