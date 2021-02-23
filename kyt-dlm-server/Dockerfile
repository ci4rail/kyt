FROM golang:1.16 AS build
WORKDIR /go/src/github.com/ci4rail/
COPY . /go/src/github.com/ci4rail/kyt
ENV CGO_ENABLED=0
ENV GOPATH=/go
ARG VERSION=dev

WORKDIR /go/src/github.com/ci4rail/kyt/kyt-server-common
RUN make test

WORKDIR /go/src/github.com/ci4rail/kyt/kyt-dlm-server
RUN make test && VERSION=${VERSION} make BIN_DIR=/install

FROM mcr.microsoft.com/azure-cli:2.19.1
RUN az extension add --name azure-iot
RUN apk add --update-cache \
    ca-certificates && \
    rm -rf /var/cache/apk/*
ENV GIN_MODE=release
COPY --from=build /install/kyt-dlm-server /kyt-dlm-server
EXPOSE 8080/tcp
ENTRYPOINT ["/kyt-dlm-server"]
