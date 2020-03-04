FROM golang:1.14-alpine as builder

RUN apk update && apk upgrade && apk add --no-cache alpine-sdk

ADD . $GOPATH/src/github.com/sylr/prometheus-azure-exporter
WORKDIR $GOPATH/src/github.com/sylr/prometheus-azure-exporter

RUN uname -a && go version
RUN git update-index --refresh; make install-static

# -----------------------------------------------------------------------------

FROM scratch

WORKDIR /usr/local/bin
COPY --from=builder "/go/bin/prometheus-azure-exporter" .

ENTRYPOINT ["/usr/local/bin/prometheus-azure-exporter"]
