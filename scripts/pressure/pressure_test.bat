@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

echo.
echo ===============================================
echo           E-commerce Pressure Test Tool
echo ===============================================
echo.

:: 设置变量
set "BASE_URL=http://localhost:8080"
set "TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE3NzM0Mzk0MTYsImlzcyI6ImUtY29tbWVyY2UifQ.yK3RIBXqkceYq9I9DAr99_I9ARAfgcC6FTb3bhQC4lg"

:: 检查服务是否运行
echo [INFO] Checking if service is running...
powershell -Command "try { Invoke-WebRequest -Uri '%BASE_URL%/product/list' -TimeoutSec 3 -UseBasicParsing > $null; Write-Output 'Service is running' } catch { Write-Output 'Service not running'; exit 1 }"
if !errorlevel! neq 0 (
    echo [ERROR] Service is not running
    echo Please start: go run main.go
    pause
    exit /b 1
)

echo [SUCCESS] Service is running

:: 创建日志文件
set "LOG_FILE=pressure_results.log"
echo ===== Pressure Test Started ===== > "%LOG_FILE%"
echo Base URL: %BASE_URL% >> "%LOG_FILE%"
echo. >> "%LOG_FILE%"

echo.
echo [INFO] Starting pressure test...
echo Log file: %LOG_FILE%
echo.

:: 1. 购物车添加测试
echo [1] Testing cart add...
echo [1] Cart Add Test >> "%LOG_FILE%"
hey.exe -c 20 -z 30s -m POST -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"quantity\":2}" %BASE_URL%/cart/add/2 >> "%LOG_FILE%"
echo. >> "%LOG_FILE%"
echo [SUCCESS] Cart add test completed

:: 2. 订单创建测试
echo [2] Testing order create...
echo [2] Order Create Test >> "%LOG_FILE%"
hey.exe -c 30 -z 30s -m POST -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"address\":\"Test Address\"}" %BASE_URL%/order/create >> "%LOG_FILE%"
echo. >> "%LOG_FILE%"
echo [SUCCESS] Order create test completed

:: 3. 订单列表测试
echo [3] Testing order list...
echo [3] Order List Test >> "%LOG_FILE%"
hey.exe -c 50 -z 30s -H "Authorization: Bearer %TOKEN%" %BASE_URL%/order/list >> "%LOG_FILE%"
echo. >> "%LOG_FILE%"
echo [SUCCESS] Order list test completed

:: 4. 商品列表测试
echo [4] Testing product list...
echo [4] Product List Test >> "%LOG_FILE%"
hey.exe -c 100 -z 30s %BASE_URL%/product/list >> "%LOG_FILE%"
echo. >> "%LOG_FILE%"
echo [SUCCESS] Product list test completed

:: 显示结果
echo.
echo ===============================================
echo              Test Results
echo ===============================================
echo.

echo Key metrics from log file:
type "%LOG_FILE%" | findstr /i "requests/sec\|average\|fastest\|slowest\|errors"

echo.
echo Full results saved to: %LOG_FILE%
echo.

:: 询问是否查看日志
echo Open log file? (y/n)
set /p choice=
if /i "!choice!"=="y" (
    start notepad "%LOG_FILE%"
)

echo.
echo Test completed!
pause