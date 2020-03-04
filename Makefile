GIT_DESCRIBE ?= $(shell git describe --tags --dirty --broken || git rev-parse --short HEAD)

# -- build --------------------------------------------------------------------

.PHONY: build debug test install install-static

build:
	go build -ldflags "-w -s -X main.version=$(GIT_DESCRIBE)"

debug:
	go build -ldflags "-X main.version=$(GIT_DESCRIBE)"

test:
	go test -ldflags "-X main.version=$(GIT_DESCRIBE)" ./...

install:
	go install -ldflags "-w -s -X main.version=$(GIT_DESCRIBE)"

install-static:
	go install -ldflags "-linkmode external -extldflags -static -w -s -X main.version=$(GIT_DESCRIBE)"

# -- go -----------------------------------------------------------------------

.PHONY: go-mod-download go-mod-tidy go-dep-upgrade-minor go-dep-upgrade-major verify-go-mod

go-mod-download:
	go mod download

go-mod-tidy:
	go mod tidy

verify-go-mod: go-mod-download
	git diff --quiet go.mod go.sum || git diff --exit-code go.mod go.sum
