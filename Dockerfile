FROM golang:1.19.3 AS builder
WORKDIR /build

COPY ./goql-plugin/go.mod ./goql-plugin/go.sum /build/
RUN go mod download

COPY ./goql-plugin/ .
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o goql ./cmd

FROM kong:3.0.1-ubuntu
USER root

RUN apt update && apt install -y curl git gcc musl-dev ca-certificates
RUN luarocks install luaossl OPENSSL_DIR=/usr/local/kong CRYPTO_DIR=/usr/local/kong

COPY --from=builder /build/goql /usr/local/bin/goql
RUN chown kong:0 /usr/local/bin/goql && chmod 755 /usr/local/bin/goql

COPY --from=builder /build/postgres /usr/local/postgres

USER kong
