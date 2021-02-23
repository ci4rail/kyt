BIN_DIR ?= ../bin
VERSION ?= $(shell git describe --match=NeVeRmAtCh --always --abbrev=40 --dirty)
GO_LDFLAGS = -ldflags "-X github.com/ci4rail/kyt/kyt-dlm-server/cmd.version=$(VERSION)"

build:
	GOOS=linux GOARCH=amd64 go build $(GO_LDFLAGS) -o ${BIN_DIR}/kyt-dlm-server main.go
test:
	go test ./...
clean:
	rm -f ${BIN_DIR}/kyt-dlm-server

.PHONY: clean
