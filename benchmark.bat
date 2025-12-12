@echo off
setlocal enabledelayedexpansion

chcp 65001 >nul

echo.
echo 🔍 Running benchmarks...
echo.

go test -bench=BenchmarkParseSlowLog -benchtime=5s ./slowsql
