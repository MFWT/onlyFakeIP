@echo off

go build -o onlyFakeIP.exe main.go

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64

go build -o onlyFakeIP main.go