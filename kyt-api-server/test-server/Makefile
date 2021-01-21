BIN_DIR ?= ../../bin
OPENAPI_GEN_DIR ?= openapi_gen

build:
	cp ${OPENAPI_GEN_DIR}/model_*.go openapi
	GOOS=linux GOARCH=386 go build -o ${BIN_DIR}/kyt-api-server-test main.go

clean:
	rm -rf ${BIN_DIR}/kyt-api-server-test

.PHONY: clean
