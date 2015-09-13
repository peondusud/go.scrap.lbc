.PHONY: clean build

GOCMD=go
GOBUILD=$(GOCMD) build
GOFLAGS ?= $(GOFLAGS:)
SOURCE=src/lbc.go
EXECUTABLE=build/lbc.bin
GCFLAGS=""


build:
	@$(GOBUILD) -o $(EXECUTABLE)  $(GOFLAGS) $(SOURCE)

clean:
	@rm -rf $(EXECUTABLE)
