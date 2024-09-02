# Makefile for TDX Guest Example

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
GORUN=$(GOCMD) run
GOMOD=$(GOCMD) mod
BINARY_NAME=tdx-guest-example
MAIN_PACKAGE=.

.PHONY: all build clean run deps tidy

all: deps tidy build

build:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PACKAGE)

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f quote.dat

run: build
	sudo ./$(BINARY_NAME)

run-direct:
	sudo $(GORUN) $(MAIN_PACKAGE)

deps:
	$(GOGET) github.com/google/go-tdx-guest@latest
	$(GOGET) github.com/google/go-configfs-tsm@latest
	$(GOGET) google.golang.org/protobuf@latest
	$(GOGET) golang.org/x/crypto@latest
	$(GOGET) github.com/google/logger@latest
	$(GOGET) golang.org/x/sys@latest

tidy:
	$(GOMOD) tidy

.DEFAULT_GOAL := all