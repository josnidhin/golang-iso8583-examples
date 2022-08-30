#
# @author Jose Nidhin
#
GO_SRC_FILES := $(shell find . -type f -name '*.go')
GO_SRC_MAIN := $(shell ls *.go)

all: tidy vet fmt simplify test clean run 

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

.PHONY: run
run:
	go run $(GO_SRC_MAIN)
