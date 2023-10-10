.PHONY: all windows mac linux wasm

BINARY_NAME=dirry
CLI_PATH=pkg/cli/main.go
WASM_PATH=pkg/wasm/main.go
DIST=dist
LDFLAGS=-ldflags="-s -w"

all: windows mac linux wasm

windows:
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(DIST)/windows/$(BINARY_NAME).exe $(CLI_PATH)

mac:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(DIST)/mac/$(BINARY_NAME) $(CLI_PATH)

linux:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(DIST)/linux/$(BINARY_NAME) $(CLI_PATH)

wasm:
	GOOS=js GOARCH=wasm go build $(LDFLAGS) -o $(DIST)/web/main.wasm $(WASM_PATH)
	cp "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" $(DIST)/web/
	cp pkg/wasm/web/* $(DIST)/web/

run:
	go run $(CLI_PATH) $(ARGS)
