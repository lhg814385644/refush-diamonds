# Makefile
PREFIX=$(shell pwd)
OUTPUT_DIR=${PREFIX}/bin

build:
	@echo "building in linux"
	go build -o ${OUTPUT_DIR}/app

clean:
	rm -rf ${OUTPUT_DIR}/*

build-windows:
	@echo "building in windows"
	GOOS=windows GOARCH=amd64 go build -o ${OUTPUT_DIR}/app.exe