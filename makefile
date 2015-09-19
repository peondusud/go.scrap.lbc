.PHONY: clean build install format get

GOROOT := /usr/lib/go
GOPATH :=  $(shell pwd)
GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOFORMAT=$(GOCMD) fmt
GOGET=$(GOCMD) get
GOFLAGS ?= $(GOFLAGS:)
SOURCE=src/lbc/*.go
EXECUTABLE=build/lbc.bin
GCFLAGS=""


build:
	@$(GOBUILD) -o $(EXECUTABLE)  $(GOFLAGS) $(SOURCE)

install:
	export GOPATH=$(GOPATH)
	export GOROOT=$(GOROOT)
	$(GOINSTALL) $(SOURCE)

format:
	@$(GOFORMAT) $(SOURCE)

clean:
	@rm -rf $(EXECUTABLE)

get:
	export GOPATH=$(GOPATH)
	export GOROOT=$(GOROOT)
	@$(GOGET) -v golang.org/x/net/html
	@$(GOGET) -v golang.org/x/text/encoding/charmap
	@$(GOGET) -v launchpad.net/xmlpath
