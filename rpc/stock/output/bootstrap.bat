@echo off
setlocal

:: 获取当前脚本所在目录
set "CURDIR=%~dp0"

:: 检查是否传递了参数
if "%~1" neq "" (
    set "RUNTIME_ROOT=%~1"
) else (
    set "RUNTIME_ROOT=%CURDIR%"
)

:: 设置环境变量
set "KITEX_RUNTIME_ROOT=%RUNTIME_ROOT%"
set "KITEX_LOG_DIR=%RUNTIME_ROOT%\log"

:: 创建日志目录
if not exist "%KITEX_LOG_DIR%\app" (
    mkdir "%KITEX_LOG_DIR%\app"
)

if not exist "%KITEX_LOG_DIR%\rpc" (
    mkdir "%KITEX_LOG_DIR%\rpc"
)

:: 执行应用程序
call "%CURDIR%\bin\example.shop.stock.exe"

endlocal