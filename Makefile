PROJECT := BinkNetUI
ROOTDIR := $(shell pwd)
VERSION := $(shell cat VERSION)
COMMIT := $(shell git rev-parse --short HEAD)

BINNAME := bnUI

SOURCES := $(shell find . -name '*.go')

.PHONY: all clean deps bootstrap binaries package test

all: binaries

clean:
	rm -Rf $(ROOTDIR)/bin

bootstrap:
	GO111MODULE=on go get github.com/lucor/fyne-cross/v2/cmd/fyne-cross

binaries: $(SOURCES)
	mkdir -p bin
	CGO_ENABLED=1 go build \
		-o ./bin/$(BINNAME) \
		-ldflags="-X main.projectVersion=$(VERSION) -X main.projectBuild=$(COMMIT)" \
		github.com/binkynet/BinkyNetUI

package:
	fyne-cross darwin

test:
	go test ./...

.PHONY: update-modules
update-modules:
	rm -f go.mod go.sum 
	go mod init github.com/binkynet/BinkyNetUI
	go mod edit \
		-replace github.com/coreos/go-systemd=github.com/coreos/go-systemd@e64a0ec8b42a61e2a9801dc1d0abe539dea79197
	go get -u \
		github.com/binkynet/BinkyNet@0.1.1
	go mod tidy
