GO           ?= go
GOENV_GOOS   := $(shell go env GOOS)
GOENV_GOARCH := $(shell go env GOARCH)
GOENV_GOARM  := $(shell go env GOARM)
GOOS         ?= $(GOENV_GOOS)
GOARCH       ?= $(GOENV_GOARCH)
GOARM        ?= $(GOENV_GOARM)
GIT_REVISION ?= $(shell git rev-parse HEAD)
GIT_VERSION  ?= $(shell git describe --always --tags --dirty --broken 2>/dev/null || echo dev)

DOCKER_BUILD_IMAGE      ?= ghcr.io/sylr/prometheus-azure-exporter
DOCKER_BUILD_VERSION    ?= $(GIT_VERSION)
DOCKER_BUILD_GO_VERSION ?= 1.15
DOCKER_BUILD_LABELS      = --label org.opencontainers.image.title=prometheus-azure-exporter
DOCKER_BUILD_LABELS     += --label org.opencontainers.image.description="Azure metrics exporter for prometheus"
DOCKER_BUILD_LABELS     += --label org.opencontainers.image.url="https://github.com/sylr/prometheus-azure-exporter"
DOCKER_BUILD_LABELS     += --label org.opencontainers.image.revision=$(GIT_REVISION)
DOCKER_BUILD_LABELS     += --label org.opencontainers.image.version=$(GIT_VERSION)
DOCKER_BUILD_LABELS     += --label org.opencontainers.image.created=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
DOCKER_BUILD_BUILD_ARGS ?= --build-arg=GO_VERSION=$(DOCKER_BUILD_GO_VERSION)
DOCKER_BUILDX_PLATFORMS ?= "linux/amd64,linux/arm64,linux/arm/v7,linux/arm/v6"
DOCKER_BUILDX_CACHE     ?= /tmp/.buildx-cache

ifeq ($(GOOS)/$(GOARCH),$(GOENV_GOOS)/$(GOENV_GOARCH))
GO_BUILD_TARGET          := prometheus-azure-exporter
GO_BUILD_VERSION_TARGET  := prometheus-azure-exporter-$(GIT_VERSION)
else
ifeq ($(GOARCH),arm)
GO_BUILD_TARGET          := prometheus-azure-exporter-$(GOOS)-$(GOARCH)
GO_BUILD_VERSION_TARGET  := prometheus-azure-exporter-$(GIT_VERSION)-$(GOOS)-$(GOARCH)
else
GO_BUILD_TARGET          := prometheus-azure-exporter-$(GOOS)-$(GOARCH)$(GOARM)
GO_BUILD_VERSION_TARGET  := prometheus-azure-exporter-$(GIT_VERSION)-$(GOOS)-$(GOARCH)$(GOARM)
endif # ifeq ($(GOARCH),arm)
endif # ifeq ($(GOOS)/$(GOARCH),$(GOENV_GOOS)/$(GOENV_GOARCH))

# -- build ---------------------------------------------------------------------

.PHONY: build build-static debug test install install-static

build:
	 GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) $(GO) build -ldflags "-w -s -X main.version=$(GIT_VERSION)" -o $(GO_BUILD_TARGET)

build-static:
	GIT_VERSION=$(GIT_VERSION) GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) CGO_ENABLED=0 $(GO) build -tags netgo -ldflags "-extldflags '-static' -w -s -X main.version=$(GIT_VERSION)" -o $(GO_BUILD_TARGET)

debug:
	GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) $(GO) build -ldflags "-X main.version=$(GIT_VERSION)"

test:
	GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) $(GO) test -ldflags "-X main.version=$(GIT_VERSION)" ./...

install:
	GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) $(GO) install -ldflags "-w -s -X main.version=$(GIT_VERSION)"

install-static:
	GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) CGO_ENABLED=0 $(GO) install -a -tags netgo -ldflags "-extldflags '-static' -w -s -X main.version=$(GIT_VERSION)"

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

# -- docker --------------------------------------------------------------------

.PHONY: docker-build docker-push docker-buildx-build docker-buildx-push docker-buildx-inspect

docker-build:
	@docker build . -f Dockerfile \
		-t $(DOCKER_BUILD_IMAGE):$(DOCKER_BUILD_VERSION)-debian \
		$(DOCKER_BUILD_BUILD_ARGS) \
		$(DOCKER_BUILD_LABELS)

docker-push:
	@docker push $(DOCKER_BUILD_IMAGE):$(DOCKER_BUILD_VERSION)-debian

docker-buildx-build:
	@docker buildx build . -f Dockerfilex \
		-t $(DOCKER_BUILD_IMAGE):$(DOCKER_BUILD_VERSION) \
		--cache-to=type=local,dest=$(DOCKER_BUILDX_CACHE) \
		--platform=$(DOCKER_BUILDX_PLATFORMS) \
		$(DOCKER_BUILD_BUILD_ARGS) \
		$(DOCKER_BUILD_LABELS)

docker-buildx-push:
	@docker buildx build . -f Dockerfilex \
		-t $(DOCKER_BUILD_IMAGE):$(DOCKER_BUILD_VERSION) \
		--cache-from=type=local,src=$(DOCKER_BUILDX_CACHE) \
		--platform=$(DOCKER_BUILDX_PLATFORMS) \
		$(DOCKER_BUILD_BUILD_ARGS) \
		$(DOCKER_BUILD_LABELS) \
		--push

docker-buildx-inspect:
	@docker buildx imagetools inspect $(DOCKER_BUILD_IMAGE):$(DOCKER_BUILD_VERSION)
