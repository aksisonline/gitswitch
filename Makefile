.PHONY: build install clean

VERSION ?= dev

build:
	go build -trimpath -ldflags="-s -w -X main.version=$(VERSION)" -o bin/gitswitch ./cmd/gitswitch

install: build
	cp bin/gitswitch /usr/local/bin/

install-brew:
	brew install --build-from-source ./Formula/gitswitch.rb

clean:
	rm -rf bin/
	go clean

test:
	go test ./...

fmt:
	go fmt ./...
