FROM golang:1.15.6 AS build
WORKDIR /go/src/github.com/ci4rail/
COPY . /go/src/github.com/ci4rail/kyt
WORKDIR /go/src/github.com/ci4rail/kyt/kyt-dlm-devinfo-static

ENV CGO_ENABLED=0
ENV GOPATH=/go
ARG VERSION=dev
RUN VERSION=${VERSION} make BIN_DIR=/install

FROM scratch
COPY --from=build /install/kyt-dlm-devinfo-static /kyt-dlm-devinfo-static
ENTRYPOINT ["/kyt-dlm-devinfo-static"]
