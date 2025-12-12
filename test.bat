@echo off
setlocal enabledelayedexpansion

chcp 65001 >nul

set GOCOVERDIR=coverage

if exist %GOCOVERDIR% rmdir /s /q %GOCOVERDIR%

mkdir %GOCOVERDIR%

go test -coverprofile=%GOCOVERDIR%\coverage.out ./slowsql
if !ERRORLEVEL! neq 0 (
    echo ❌ Error: there are test failures
    exit /b 1
)
go tool cover -html=%GOCOVERDIR%\coverage.out -o %GOCOVERDIR%\coverage.html

echo.
echo ✅ Test coverage report generated: %GOCOVERDIR%\coverage.html
echo.

if "-html" == "%1" (
    cd %GOCOVERDIR%
    start "" "coverage.html"
)
