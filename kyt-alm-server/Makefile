BIN_DIR ?= ../bin
VERSION ?= $(shell git describe --tags --dirty --exact-match 2>/dev/null || git rev-parse --short HEAD)
GO_LDFLAGS = -ldflags "-X github.com/ci4rail/kyt/kyt-alm-server/cmd.version=$(VERSION)"

build:
	GOOS=linux GOARCH=386 go build $(GO_LDFLAGS) -o ${BIN_DIR}/kyt-alm-server main.go

clean:
	rm -f ${BIN_DIR}/kyt-alm-server

.PHONY: clean
