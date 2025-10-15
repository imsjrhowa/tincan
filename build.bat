@echo off
REM Build script for TinCan (non-embedded version)

REM Get git commit hash
for /f "tokens=*" %%i in ('git rev-parse --short HEAD 2^>nul') do set GIT_COMMIT=%%i
if "%GIT_COMMIT%"=="" set GIT_COMMIT=unknown

REM Get current date/time in ISO format
for /f "tokens=*" %%i in ('powershell -Command "Get-Date -Format 'yyyy-MM-ddTHH:mm:ssZ'"') do set BUILD_DATE=%%i
if "%BUILD_DATE%"=="" set BUILD_DATE=unknown

REM Set version (can be overridden by environment variable)
if "%VERSION%"=="" set VERSION=1.0.0

echo Building TinCan v%VERSION% (commit: %GIT_COMMIT%, date: %BUILD_DATE%)

REM Build the binary with version information
go build -ldflags "-X main.Version=%VERSION% -X main.GitCommit=%GIT_COMMIT% -X main.BuildDate=%BUILD_DATE%" -o tincan.exe ./cmd/tincan

if %ERRORLEVEL% EQU 0 (
    echo.
    echo Build successful! Binary created: tincan.exe
    echo.
    echo To verify, run: tincan.exe version
) else (
    echo.
    echo Build failed with error code %ERRORLEVEL%
    exit /b %ERRORLEVEL%
)
