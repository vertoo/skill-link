# Makefile for skill-link

APP_NAME=skill-link
BUILD_DIR=build
SRC=./cmd/skill-link

.PHONY: all clean build build-linux build-darwin build-windows

all: clean build

build:
	@echo "Building for host OS..."
	CGO_ENABLED=0 go build -o ${APP_NAME} ${SRC}

build-linux:
	@echo "Building for Linux (x86_64/arm64)..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${BUILD_DIR}/${APP_NAME}-linux-amd64 ${SRC}
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o ${BUILD_DIR}/${APP_NAME}-linux-arm64 ${SRC}

build-darwin:
	@echo "Building for macOS (x86_64/arm64)..."
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o ${BUILD_DIR}/${APP_NAME}-darwin-amd64 ${SRC}
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o ${BUILD_DIR}/${APP_NAME}-darwin-arm64 ${SRC}

build-windows:
	@echo "Building for Windows (x86_64)..."
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o ${BUILD_DIR}/${APP_NAME}-windows-amd64.exe ${SRC}

cross-compile: build-linux build-darwin build-windows

clean:
	@echo "Cleaning up..."
	@rm -rf ${BUILD_DIR}
	@rm -f ${APP_NAME}
