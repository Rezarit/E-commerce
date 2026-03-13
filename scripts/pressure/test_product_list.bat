@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

echo.
echo ===============================================
echo           商品列表接口压测工具
echo ===============================================
echo.

:: 设置变量
set "BASE_URL=http://localhost:8080"
set "LOG_FILE=product_list_test_%date:~0,4%%date:~5,2%%date:~8,2%.log"

:: 检查服务是否运行
echo [INFO] 检查服务状态...
powershell -Command "try { Invoke-WebRequest -Uri '%BASE_URL%/product/list' -TimeoutSec 3 -UseBasicParsing > $null; Write-Output '服务运行正常' } catch { Write-Output '服务未运行'; exit 1 }"
if !errorlevel! neq 0 (
    echo [ERROR] 服务未运行
    echo 请先启动: go run main.go
    pause
    exit /b 1
)

:: 创建日志文件
echo ===== 商品列表接口压测开始 ===== > "%LOG_FILE%"
echo 目标URL: %BASE_URL%/product/list >> "%LOG_FILE%"
echo 测试时间: 30秒 >> "%LOG_FILE%"
echo. >> "%LOG_FILE%"

echo.
echo [INFO] 开始压测商品列表接口（30秒）...
echo 日志文件: %LOG_FILE%
echo.

:: 压测商品列表接口
echo [1] 商品列表接口压测 - %date% %time% >> "%LOG_FILE%"
hey.exe -c 50 -z 30s %BASE_URL%/product/list >> "%LOG_FILE%"
echo. >> "%LOG_FILE%"

echo [SUCCESS] 压测完成

:: 显示结果
echo.
echo ===============================================
echo             压测结果
echo ===============================================
echo.

echo 关键指标:
type "%LOG_FILE%" | findstr /i "requests/sec\|average\|fastest\|slowest\|errors"

echo.
echo 完整结果已保存到: %LOG_FILE%
echo.

:: 询问是否查看日志
echo 是否打开日志文件查看详细结果? (y/n)
set /p choice=
if /i "!choice!"=="y" (
    start notepad "%LOG_FILE%"
)

echo.
echo 测试完成!
pause