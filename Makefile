PROJECT := BinkNetUI
ROOTDIR := $(shell pwd)
VERSION := $(shell cat VERSION)
COMMIT := $(shell git rev-parse --short HEAD)

ORGPATH := github.com/binkynet
REPONAME := $(PROJECT)
REPOPATH := $(ORGPATH)/$(REPONAME)
BINNAME := bnUI

SOURCES := $(shell find . -name '*.go')

.PHONY: all clean deps bootstrap binaries test

all: binaries

clean:
	rm -Rf $(ROOTDIR)/bin

bootstrap:
	go get github.com/mitchellh/gox

binaries: $(SOURCES)
	gox \
		-osarch="darwin/amd64 windows/amd64" \
		-ldflags="-X main.projectVersion=$(VERSION) -X main.projectBuild=$(COMMIT)" \
		-output="bin/{{.OS}}/{{.Arch}}/$(BINNAME)" \
		-tags="netgo" \
		./...

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
