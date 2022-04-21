TARGET=lightbull

all: build-x64-linux< build-armv7-linux

prepare:
	mkdir -p ./build

clean:
	rm -rf ./build

build-x64-linux: prepare
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/$(TARGET)-x64-linux .

# Raspberry Pi 3
build-armv7-linux: prepare
	GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build -o build/$(TARGET)-armv7-linux .
