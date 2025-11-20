@echo off
echo ========================================
echo Cleaning Build Artifacts
echo ========================================
echo.

REM Remove build directory
if exist "build\bin" (
    echo Removing build\bin...
    rmdir /s /q "build\bin"
)

REM Remove frontend dist directory
if exist "frontend\dist" (
    echo Removing frontend\dist...
    rmdir /s /q "frontend\dist"
)

REM Remove frontend node_modules if exists
if exist "frontend\node_modules" (
    echo Removing frontend\node_modules...
    rmdir /s /q "frontend\node_modules"
)

REM Clean Go cache
echo Cleaning Go cache...
go clean -cache -modcache -i -r

echo.
echo ========================================
echo Clean complete!
echo ========================================
echo.
pause
