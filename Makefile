GIT_DESCRIBE ?= $(shell git describe --tags --dirty --broken || git rev-parse --short HEAD)

# -- build --------------------------------------------------------------------

.PHONY: build build-static debug test install install-static

build:
	go build -ldflags "-w -s -X main.version=$(GIT_DESCRIBE)"

build-static:
	CGO_ENABLED=0 go build -tags netgo -ldflags "-extldflags '-static' -w -s -X main.version=$(GIT_DESCRIBE)"

debug:
	go build -ldflags "-X main.version=$(GIT_DESCRIBE)"

test:
	go test -ldflags "-X main.version=$(GIT_DESCRIBE)" ./...

install:
	go install -ldflags "-w -s -X main.version=$(GIT_DESCRIBE)"

install-static:
	CGO_ENABLED=0 go install -a -tags netgo -ldflags "-extldflags '-static' -w -s -X main.version=$(GIT_DESCRIBE)"

# -- go -----------------------------------------------------------------------

.PHONY: go-mod-download go-mod-tidy go-dep-upgrade-minor go-dep-upgrade-major verify-go-mod

go-mod-download:
	go mod download

go-mod-tidy:
	go mod tidy

go-dep-upgrade-minor:
	go get -u=patch ./...

go-dep-upgrade-major:
	go get -u ./...

verify-go-mod: go-mod-download
	@git diff --quiet go.mod go.sum || { \
	    git diff go.mod go.sum; \
	    git diff --exit-code go.mod go.sum; \
	}
