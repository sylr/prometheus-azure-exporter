GIT_DESCRIBE ?= $(shell git describe --tags --dirty --broken || git rev-parse --short HEAD)

build:
	go build -ldflags "-X main.version=$(GIT_DESCRIBE)"

install:
	go install -ldflags "-X main.version=$(GIT_DESCRIBE)"

go-mod-tidy:
	go mod tidy

.PHONY: build install go-mod-tidy
