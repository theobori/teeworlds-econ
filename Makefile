# Go files to format
GOFMT_FILES ?= $(shell find . -name "*.go")

default: fmt

fmt:
	gofmt -w $(GOFMT_FILES)

clean:
	go clean -testcache

test: clean
	go test -v ./...

.PHONY: \
	fmt \
	test \
	clean
