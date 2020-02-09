FROM golang:1.12-alpine

RUN set -ex; \
    apk update; \
    apk add --no-cache git

ENV GO111MODULE=on \
    GCO_ENABLED=0
WORKDIR /go/src/github.com/thikonom/tinyurl/
COPY . .
RUN go get -d -v

CMD go run ./scripts/tinyurl.go
