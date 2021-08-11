APP := $(shell basename $(CURDIR))
VERSION := $(shell git describe --tags --always --dirty)
GOPATH := $(CURDIR)/Godeps/_workspace:$(GOPATH)
PATH := $(GOPATH)/bin:$(PATH)

# HELP
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
#.PHONY: bin/$(APP) bin clean start test

help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

# Go TASKS
build:  bin/$(APP)  test ## Build Go

bin/$(APP): bin
	go build -v -o $@ -ldflags "-X main.Version='${VERSION}'"

bin: clean
	mkdir -p bin

clean: ## Clean Go
	rm -rf bin

lint: ## Go Lint
	golangci-lint run

shiftleft: ## Adopt Shift Left
	npm --prefix shift-left install

test:  ## Test Go
	go test ./...

test-iaac: ## End to End Test for the Dev Container
	goss -g .devcontainer/go-goss.yaml validate --format tap

test-watch:  ## Test Watch
	gotestsum --watch --format testname

test-cover: ## Go Coverage
	go test ./... --cover -coverprofile coverage/coverage.out
	go tool cover -html=coverage/coverage.out

test-summary: ## Prints formatted test output, and a summary of the test run
	gotestsum --format testname -- -coverprofile=coverage/coverage.out ./...
