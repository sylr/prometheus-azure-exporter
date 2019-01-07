FROM golang:1.11-alpine3.8 as builder

ADD . $GOPATH/src/github.com/sylr/prometheus-azure-exporter
WORKDIR $GOPATH/src/github.com/sylr/prometheus-azure-exporter

RUN apk update && apk upgrade && apk add --no-cache git

RUN uname -a && go version && git update-index --refresh && (git describe --tags --dirty --broken || git rev-parse --short HEAD)
RUN git update-index --refresh && go install -ldflags "-X main.Version=$(git describe --tags --dirty --broken || git rev-parse --short HEAD)"

# -----------------------------------------------------------------------------

FROM alpine:3.8

WORKDIR /usr/local/bin
RUN apk --no-cache add ca-certificates
RUN apk update && apk upgrade && apk add --no-cache bash curl
COPY --from=builder "/go/bin/prometheus-azure-exporter" .

ENTRYPOINT ["/usr/local/bin/prometheus-azure-exporter"]
