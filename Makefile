VERSION = "v0.1.0"

build:
	@go build .

.PHONY: all
all: build build-image

build-image:
	@docker build . -t cijie/goproxy:$(VERSION)