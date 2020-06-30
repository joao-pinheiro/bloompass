# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

all: build

build:
	$(GOBUILD) -tags static -o bin/bloompass cmd/bloompass/main.go

clean:
	$(GOCLEAN)
	rm -f bin/*
