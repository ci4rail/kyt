BIN_DIR ?= ../bin
VERSION ?= $(shell git describe --tags --dirty --exact-match 2>/dev/null || git rev-parse --short HEAD)
GO_LDFLAGS = -ldflags "-X github.com/ci4rail/kyt/kyt-dlm-devinfo-static/cmd.version=$(VERSION)"

build:
	GOOS=linux GOARCH=arm64 go build $(GO_LDFLAGS) -o ${BIN_DIR}/kyt-dlm-devinfo-static main.go

clean:
	rm -f ${BIN_DIR}/kyt-dlm-devinfo-static

.PHONY: clean
