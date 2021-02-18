BIN_DIR ?= ../bin
VERSION ?= $(shell git describe --tags --dirty --exact-match 2>/dev/null || git rev-parse --short HEAD)
GO_LDFLAGS = -ldflags "-X github.com/ci4rail/kyt/kyt-cli/cmd.version=$(VERSION)"

build:
	GOOS=linux GOARCH=386 go build $(GO_LDFLAGS) -o ${BIN_DIR}/kyt main.go
	GOOS=windows GOARCH=386 go build $(GO_LDFLAGS) -o ${BIN_DIR}/kyt.exe main.go

test:
	go test ./...

clean:
	rm -rf ${BIN_DIR}/kyt ${BIN_DIR}/kyt.exe

.PHONY: clean
