BIN_DIR ?= ../bin
VERSION ?= $(shell git describe --tags --dirty --exact-match 2>/dev/null || git rev-parse --short HEAD)
GO_LDFLAGS = -ldflags "-X github.com/ci4rail/kyt/kyt-alm-server/internal/version/version.Version=$(VERSION)"

build:
	GOOS=linux GOARCH=amd64 go build $(GO_LDFLAGS) -o ${BIN_DIR}/kyt-alm-server main.go

test:
	go test ./...
clean:
	rm -f ${BIN_DIR}/kyt-alm-server

.PHONY: clean
