@echo off
set MICROSERVICES=cmd/config-seed/config-seed cmd/export-client/export-client cmd/export-distro/export-distro cmd/core-metadata/core-metadata cmd/core-data/core-data cmd/core-command/core-command cmd/support-logging/support-logging cmd/support-notifications/support-notifications cmd/sys-mgmt-executor/sys-mgmt-executor cmd/sys-mgmt-agent/sys-mgmt-agent cmd/support-scheduler/support-scheduler
set DOCKER_TAG=%VERSION%-dev
set FILE=cmd.sh
for /f "delims=" %%i in ('git rev-parse HEAD') do set GIT_SHA=%%
echo #!/usr/bin/env bash > %FILE%

Setlocal ENABLEDELAYEDEXPANSION
set GO=CGO_ENABLED=0 GO111MODULE=on go
set GOCGO=CGO_ENABLED=1 GO111MODULE=on go
set /P VERSION=<VERSION
set GOFLAGS=-ldflags "-X github.com/edgexfoundry/edgex-go.Version=%VERSION%"
set ch1=/
for %%i in (%MICROSERVICES%) do (
    call :findChar %%i %ch1%
)

::注意，这里是区分大小写的！
:findChar

set result=./
set OBJECT=./%1.exe
set module=
set str=%1
set ch=%2
set count=0
set num=0
::复制字符串，用来截短，而不影响源字符串
:next
if not "%str%"=="" (
  set /a num+=1
  if "!str:~0,1!"=="%ch%" (
    if 0==%count% (set /a count+=1) else goto last
  )
  if 1==%count% (set "module=!module!!str:~0,1!")
  ::比较首字符是否为要求的字符，如果是则跳出循环
  set "result=!result!!str:~0,1!"
  set "str=%str:~1%"
  goto next
)
set /a num=0
::没有找到字符时，将num置零
:last
:: 根据模块名设置go环境变量 core-data export-distro要使用GOCGO编译
::
if !module! == core-data (
    set "GOCMD=!GOCGO! build !GOFLAGS! -o !OBJECT! !result!"
    echo !GOCMD! >> %FILE%
) else (
    if !module! == export-distro (
    set "GOCMD=!GOCGO! build !GOFLAGS! -o !OBJECT! !result!"
    echo !GOCMD! >> %FILE%
    ) else (
      if not !result! == ./ (
        set "GOCMD=!GO! build !GOFLAGS! -o !OBJECT! !result!"
        echo !GOCMD! >> %FILE%
      )
    )
)

goto :eof
:: CGO_ENABLED=0 GO111MODULE=on go build -ldflags "-X github.com/edgexfoundry/edgex-go.Version=1.0.0" -o ./cmd/core-metadata/core-metadata.exe ./cmd/core-metadata

