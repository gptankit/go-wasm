GOOS=js
GOARCH=wasm
BINPATH=site/main.wasm

all: build

build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags="-s -w" -trimpath -o $(BINPATH) main.go
	@echo "done"
