FROM golang:1.22.4-alpine AS build

COPY . /sources
WORKDIR /sources

RUN go build -o ./acc-server cmd/server.go

ENTRYPOINT ["/sources/acc-server"]
