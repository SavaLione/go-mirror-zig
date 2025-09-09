SOURCES := $(wildcard *.go cmd/*/*.go)

VERSION=$(shell git describe --tags --always --dirty  --long)

build: $(SOURCES)
	go build -ldflags "-X main.version=${VERSION}" -o go-mirror-zig ./cmd/main.go

