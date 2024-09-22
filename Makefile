BINARY_NAME=xping
VERSION=1.0.0
BUILD_DIR=bin

.PHONY: all build build-linux clean

all: build

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

build: $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) -v

build-linux: $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)_linux_amd64 -v

clean:
	go clean
	rm -rf $(BUILD_DIR)
