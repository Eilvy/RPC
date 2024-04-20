@echo off
setlocal

rem 定义变量
set "RUN_NAME=example.shop.item.exe"

rem 创建 output 目录
mkdir output
mkdir output\bin

rem 复制项目启动脚本
copy script\* output\

rem 设置启动脚本的执行权限（在 Windows 上不需要）
rem 该行代码在 Windows 平台不适用，可以直接删除
rem chmod +x output\bootstrap.sh

rem 判断环境变量并编译
if "%IS_SYSTEM_TEST_ENV%" NEQ "1" (
    go build -o output\bin\%RUN_NAME%
) else (
    go test -c -covermode=set -o output\bin\%RUN_NAME% -coverpkg=./...
)

echo 编译完成

endlocal