#!/bin/bash
go build -ldflags="-s -w" -o dist/dirry main.go


GOOS=js GOARCH=wasm go build -o dist/web/main.wasm pkg/wasm/dirry-web.go
