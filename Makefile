GO           ?= go
GIT_DESCRIBE ?= $(shell git describe --always --tags --dirty --broken 2>/dev/null || git rev-parse --short HEAD)

# -- build ---------------------------------------------------------------------

.PHONY: build build-static debug test install install-static

build:
	$(GO) build -ldflags "-w -s -X main.version=$(GIT_DESCRIBE)"

build-static:
	CGO_ENABLED=0 $(GO) build -tags netgo -ldflags "-extldflags '-static' -w -s -X main.version=$(GIT_DESCRIBE)"

debug:
	$(GO) build -ldflags "-X main.version=$(GIT_DESCRIBE)"

test:
	$(GO) test -ldflags "-X main.version=$(GIT_DESCRIBE)" ./...

install:
	$(GO) install -ldflags "-w -s -X main.version=$(GIT_DESCRIBE)"

install-static:
	CGO_ENABLED=0 $(GO) install -a -tags netgo -ldflags "-extldflags '-static' -w -s -X main.version=$(GIT_DESCRIBE)"

# -- go ------------------------------------------------------------------------

.PHONY: go-mod-download go-mod-tidy go-dep-upgrade-minor go-dep-upgrade-major verify-go-mod

go-mod-download:
	$(GO) mod download

go-mod-tidy:
	$(GO) mod tidy

go-dep-upgrade-minor:
	$(GO) get -u=patch ./...

go-dep-upgrade-major:
	$(GO) get -u ./...

verify-go-mod: go-mod-download
	@git diff --quiet go.mod go.sum || { \
	    git diff go.mod go.sum; \
	    git diff --exit-code go.mod go.sum; \
	}
