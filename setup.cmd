rem @echo off

set ZH_HOME=%1
set SCRIPTDIR=%~dp0

if "%ZH_HOME%" == "" (
	echo setup.cmd targetfolder\
	exit /B 1
)

cd %SCRIPTDIR%
git pull

set FILE=zhsign.exe
set SOURCE=%SCRIPTDIR%%FILE%
set TARGET=%ZH_HOME%%FILE%

go build -o %SOURCE% cmd/zhsign/zhsign.go

if exist %TARGET%.old del /F %TARGET%.old
if exist %TARGET% ren %TARGET% %FILE%.old
copy %SOURCE% %TARGET%

set FILE=zhget.exe
set SOURCE=%SCRIPTDIR%%FILE%
set TARGET=%ZH_HOME%%FILE%

go build -o %SOURCE% cmd/zhget/zhget.go

if exist %TARGET%.old del /F %TARGET%.old
if exist %TARGET% ren %TARGET% %FILE%.old
copy %SOURCE% %TARGET%

exit /B 0
