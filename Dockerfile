# vi: ft=Dockerfile:

ARG GO_VERSION=1.15

FROM golang:$GO_VERSION as builder

RUN apt-get update && apt-get dist-upgrade -y && apt-get install -y build-essential git

WORKDIR $GOPATH/src/github.com/sylr/prometheus-azure-exporter

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN uname -a && go version
RUN git update-index --refresh || true
RUN make install

# -----------------------------------------------------------------------------

FROM debian:buster-slim

WORKDIR /usr/local/bin
RUN apt-get update && apt-get dist-upgrade -y && apt-get install -y bash curl
COPY --from=builder "/go/bin/prometheus-azure-exporter" .

ENTRYPOINT ["/usr/local/bin/prometheus-azure-exporter"]
