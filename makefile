.PHONY: clean build

GOCMD=go
GOBUILD=$(GOCMD) build
GOFLAGS ?= $(GOFLAGS:)
SOURCE=src/lbc.go
EXECUTABLE=build/lbc.exe
GCFLAGS=""


build:
	@$(GOBUILD) -o $(EXECUTABLE) -i $(GOFLAGS) $(SOURCE)

clean:
	@rm -rf $(EXECUTABLE)
