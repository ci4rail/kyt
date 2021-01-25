BIN_DIR ?= ../bin
OPENAPI_GEN_DIR ?= openapi_gen
VERSION ?= $(shell git describe --tags --dirty --exact-match 2>/dev/null || git rev-parse --short HEAD)
GO_LDFLAGS = -ldflags "-X github.com/ci4rail/kyt/kyt-api-server/cmd.version=$(VERSION)"

build:
	cp ${OPENAPI_GEN_DIR}/model_*.go openapi
	GOOS=linux GOARCH=386 go build $(GO_LDFLAGS) -o ${BIN_DIR}/kyt-api-server main.go

clean:
	rm -f ${BIN_DIR}/kyt-api-server

.PHONY: clean
