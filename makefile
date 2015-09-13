.PHONY: clean build

GOCMD=go
GOBUILD=$(GOCMD) build
GOGET=$(GOCMD) get
GOFLAGS ?= $(GOFLAGS:)
SOURCE=src/lbc.go
EXECUTABLE=build/lbc.bin
GCFLAGS=""


build:
	@$(GOBUILD) -o $(EXECUTABLE)  $(GOFLAGS) $(SOURCE)

clean:
	@rm -rf $(EXECUTABLE)

get:
	@$(GOGET) golang.org/x/net/html
	@$(GOGET) golang.org/x/text/transform
	@$(GOGET) golang.org/x/text/encoding/charmap
	@$(GOGET) launchpad.net/xmlpath
