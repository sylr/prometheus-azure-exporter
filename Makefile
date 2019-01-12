GIT_DESCRIBE ?= $(shell git describe --tags --dirty --broken || git rev-parse --short HEAD)

build:
	go build -ldflags "-X main.Version=$(GIT_DESCRIBE)"

install:
	go install -ldflags "-X main.Version=$(GIT_DESCRIBE)"

vendor:
	GO111MODULE=on go mod vendor
	git add vendor && git commit -s -m "dependencies: update vendored libs"

.PHONY: build install vendor
