BIN_DIR ?= ../bin

build:
	GOOS=linux GOARCH=386 go build -o ${BIN_DIR}/kyt-api-server main.go
	GOOS=windows GOARCH=386 go build -o ${BIN_DIR}/kyt-api-server.exe main.go

clean:
	rm -rf ${BIN_DIR}/*

.PHONY: clean
