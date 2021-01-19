BIN_DIR ?= ../bin

build:
	GOOS=linux GOARCH=386 go build -o ${BIN_DIR}/kyt main.go
	GOOS=windows GOARCH=386 go build -o ${BIN_DIR}/kyt.exe main.go

clean:
	rm -rf ${BIN_DIR}/kyt ${BIN_DIR}/kyt.exe

.PHONY: clean
