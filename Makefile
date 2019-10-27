build:
	sh script/build-image/build.sh

lint:
	golangci-lint run --print-issued-lines=false

.PHONY: build lint