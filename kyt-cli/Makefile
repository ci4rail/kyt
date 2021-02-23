BIN_DIR ?= ../bin
VERSION ?= $(shell git describe --match=NeVeRmAtCh --always --abbrev=40 --dirty)
GO_LDFLAGS = -ldflags "-X github.com/ci4rail/kyt/kyt-cli/cmd.version=$(VERSION)"
GO_ARCH = amd64
build:
	GOOS=linux GOARCH=${GO_ARCH} go build $(GO_LDFLAGS) -o ${BIN_DIR}/kyt main.go
	GOOS=windows GOARCH=${GO_ARCH} go build $(GO_LDFLAGS) -o ${BIN_DIR}/kyt.exe main.go

test:
	go test ./...

clean:
	rm -rf ${BIN_DIR}/kyt ${BIN_DIR}/kyt.exe

.PHONY: clean
