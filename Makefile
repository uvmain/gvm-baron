BINARY_NAME ?= gvm
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

LDFLAGS := -s -w
BUILD_FLAGS := -trimpath -buildvcs=false -ldflags="$(LDFLAGS)"

.PHONY: all build clean

all: build

build:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(BUILD_FLAGS) -o bin/$(BINARY_NAME) .

clean:
	rm -rf bin

dependencies:
	go mod tidy