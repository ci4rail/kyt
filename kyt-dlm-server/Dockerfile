FROM golang:1.15.6 AS build
WORKDIR /go/src/github.com/ci4rail/
COPY . /go/src/github.com/ci4rail/kyt
WORKDIR /go/src/github.com/ci4rail/kyt/kyt-dlm-server

ENV CGO_ENABLED=0
ENV GOPATH=/go
ARG VERSION=dev
RUN VERSION=${VERSION} make BIN_DIR=/install

FROM scratch
ENV GIN_MODE=release
COPY --from=build /install/kyt-dlm-server /kyt-dlm-server
EXPOSE 8080/tcp
ENTRYPOINT ["/kyt-dlm-server"]
