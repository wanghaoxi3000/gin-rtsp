@echo off
set CGO_ENABLED=1 
set GOOS=windows
set GOARCH=amd64
go build -o ./bin/windows/rtsp-relay.exe main.go