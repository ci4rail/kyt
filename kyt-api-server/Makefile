BIN_DIR ?= ../bin
OPENAPI_GEN_DIR ?= openapi_gen

build:
	cp ${OPENAPI_GEN_DIR}/model_*.go openapi
	GOOS=linux GOARCH=386 go build -o ${BIN_DIR}/kyt-api-server main.go

clean:
	rm -f ${BIN_DIR}/kyt-api-server

.PHONY: clean
