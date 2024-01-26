HERE ?= $(shell pwd)
LOCALBIN ?= $(shell pwd)/bin

.PHONY: all

all: build

.PHONY: $(LOCALBIN)
$(LOCALBIN):
	mkdir -p $(LOCALBIN)
	
build: $(LOCALBIN)
	GO111MODULE="on" go build -o $(LOCALBIN)/compspec cmd/compspec/compspec.go

build-arm: $(LOCALBIN)
	GO111MODULE="on" GOARCH=arm64 go build -o $(LOCALBIN)/compspec-arm cmd/compspec/compspec.go

build-ppc: $(LOCALBIN)
	GO111MODULE="on" GOARCH=ppc64le go build -o $(LOCALBIN)/compspec-ppc cmd/compspec/compspec.go