HERE ?= $(shell pwd)
LOCALBIN ?= $(shell pwd)/bin
LIBDIR ?= "/usr/lib64"

BUILDENVVAR=CGO_LDFLAGS="-lhwloc -lstdc++ -L${LIBDIR}"

.PHONY: all

all: build

.PHONY: $(LOCALBIN)
$(LOCALBIN):
	mkdir -p $(LOCALBIN)
	
build: $(LOCALBIN)
	$(BUILDENVVAR) GO111MODULE="on" go build -ldflags '-w' -o $(LOCALBIN)/compspec cmd/compspec/compspec.go

build-arm: $(LOCALBIN)
	$(BUILDENVVAR) GO111MODULE="on" GOARCH=arm64 go build -ldflags '-w' -o $(LOCALBIN)/compspec-arm cmd/compspec/compspec.go

build-ppc: $(LOCALBIN)
	$(BUILDENVVAR) GO111MODULE="on" GOARCH=ppc64le go build -ldflags '-w' -o $(LOCALBIN)/compspec-ppc cmd/compspec/compspec.go