#
# @author Jose Nidhin
#
PROJECT_NAME := $(shell basename "$(PWD)")
GO_SRC_FILES := $(shell find . -type f -name '*.go')
GO_SRC_MAIN := $(shell find . -type f -name '*.go' ! -name '*_test.go')
GO_ENVFLAGS=CGO_ENABLED=0

all: tidy vet fmt simplify test clean build

.PHONY: vet
vet:
	go vet ./...

.PHONY: fmt
fmt:
	gofmt -l -w $(GO_SRC_FILES)

.PHONY: simplify
simplify:
	gofmt -s -l -w $(GO_SRC_FILES)

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: test
test:
	go test -race -cover ./...

.PHONY: clean
clean:
	go clean -x
	rm -f cover.out coverage.html

.PHONY: build
build:
	go build -v -o $(PROJECT_NAME) $(GO_SRC_MAIN)

.PHONY: release-build
release-build:
	$(GO_ENVFLAGS) go build -v -o $(PROJECT_NAME) $(GO_SRC_MAIN)

.PHONY: container-build
container-build:
	docker image build -t $(PROJECT_NAME):latest .
