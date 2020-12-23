# vi: ft=Dockerfile:

ARG GO_VERSION=1.15

FROM golang:$GO_VERSION as builder

RUN apt-get update && apt-get dist-upgrade -y && apt-get install -y build-essential git

WORKDIR $GOPATH/src/github.com/sylr/prometheus-azure-exporter

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Run a git command otherwise git describe in the Makefile could report a dirty git dir
RUN git diff --exit-code

RUN make build

# -----------------------------------------------------------------------------

FROM debian:buster-slim

WORKDIR /usr/local/bin
RUN apt-get update && apt-get dist-upgrade -y && apt-get install -y bash curl

COPY --from=builder "/go/src/github.com/sylr/prometheus-azure-exporter/prometheus-azure-exporter" .

ENTRYPOINT ["/usr/local/bin/prometheus-azure-exporter"]
