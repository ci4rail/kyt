BIN_DIR ?= ../../bin

build:
	GOOS=linux GOARCH=386 go build -o ${BIN_DIR}/kyt-dlm-server-test main.go

clean:
	rm -rf ${BIN_DIR}/kyt-dlm-server-test

.PHONY: clean
