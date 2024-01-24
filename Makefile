HERE ?= $(shell pwd)
LOCALBIN ?= $(shell pwd)/bin

.PHONY: all build

all: build-extract

.PHONY: $(LOCALBIN)
$(LOCALBIN):
	mkdir -p $(LOCALBIN)
	
build-extract: $(LOCALBIN)
	GO111MODULE="on" go build -o $(LOCALBIN)/compspec cmd/compspec/compspec.go