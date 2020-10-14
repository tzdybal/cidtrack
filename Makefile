GOCC ?= go
GOFLAGS ?=

# If set, override the install location for plugins
IPFS_PATH ?= $(HOME)/.ipfs

# If set, override the IPFS version to build against. This _modifies_ the local
# go.mod/go.sum files and permanently sets this version.
IPFS_VERSION ?= $(lastword $(shell $(GOCC) list -m github.com/ipfs/go-ipfs))

# make reproducible
ifneq ($(findstring /,$(IPFS_VERSION)),)
# Locally built go-ipfs
GOFLAGS += -asmflags=all=-trimpath="$(GOPATH)" -gcflags=all=-trimpath="$(GOPATH)"
else
# Remote version of go-ipfs (e.g. via `go get -trimpath` or official distribution)
GOFLAGS += -trimpath
endif

.PHONY: build install

all: ~/.ipfs/plugins/cidtrack.so

clean:
	rm -rf $(IPFS_PATH)/plugins/cidtrack.so ../cidtrack.so
	go clean

build: cidtrack.so

install: $(IPFS_PATH)/plugins/cidtrack.so

FORCE:

go.mod: FORCE
	./set-target.sh $(IPFS_VERSION) 

test: *.go
	go test -race -cover ./...

$(IPFS_PATH)/plugins/cidtrack.so: cidtrack.so
	mkdir -p $(IPFS_PATH)/plugins
	cp cidtrack.so $(IPFS_PATH)/plugins/cidtrack.so

cidtrack.so: *.go go.mod
	$(GOCC) build $(GOFLAGS) -buildmode=plugin -o cidtrack.so
	chmod +x cidtrack.so
