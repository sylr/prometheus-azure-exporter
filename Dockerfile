FROM golang:1.14 as builder

RUN apt-get update && apt-get dist-upgrade -y && apt-get install -y build-essential git

ADD . $GOPATH/src/github.com/sylr/prometheus-azure-exporter
WORKDIR $GOPATH/src/github.com/sylr/prometheus-azure-exporter

RUN uname -a && go version
RUN git update-index --refresh; make install

# -----------------------------------------------------------------------------

FROM alpine:3.11

WORKDIR /usr/local/bin
RUN apk --no-cache add ca-certificates
RUN apk update && apk upgrade && apk add --no-cache bash curl
COPY --from=builder "/go/bin/prometheus-azure-exporter" .

ENTRYPOINT ["/usr/local/bin/prometheus-azure-exporter"]
