GIT_DESCRIBE ?= $(shell git describe --tags --dirty --broken || git rev-parse --short HEAD)
GO111MODULE  ?= on
GOPROXY      ?= https://proxy.golang.org/

export GO111MODULE
export GOPROXY

build:
	go build -ldflags "-X main.version=$(GIT_DESCRIBE)"

install:
	go install -ldflags "-X main.version=$(GIT_DESCRIBE)"

vendor:
	GO111MODULE=on go mod vendor
	git add vendor && git diff --cached --exit-code > /dev/null || git commit -s -m "Update vendored libs"

.PHONY: build install vendor
