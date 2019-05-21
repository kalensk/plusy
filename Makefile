ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
SRC_DIR=$(ROOT_DIR)/src
BINARY_NAME=plusy
FILES=$$(go list -f {{.Dir}} ./...)
PACKAGES=$$(go list ./...)
TESTFILES?=$$(go list ./... | grep -v acceptance)
ACCPT_TESTFILES=$$(go list ./... | grep acceptance )

default: build

# Bootstrap
bootstrap:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(go env GOPATH)/bin v1.15.0
	go get github.com/golang/lint/golint
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	dep ensure -v

# Formatting and Linting
format:
	@# exclude the top level directory from goimports since it does not like it, that is:
	@# /home/kalensk/go/src/github.com/kalensk/plusy/src
	goimports -l -w $(shell go list -f {{.Dir}} ./... | grep -v 'src$$')
	@# The above command is equivant to: golangci-lint run --disable-all --enable goimports

	go fmt $(PACKAGES)
	@# equivant to: golangci-lint run --disable-all --enable gofmt

lint:
	golangci-lint run

# Build and install
build: format
	go build -o $(BINARY_NAME) -race $(SRC_DIR)/main.go

install: build
	cp $(BINARY_NAME) $(GOPATH)/bin/

# type = [redis | postgres | neo4j]
database:
	bash $(ROOT_DIR)/docker/run-$(type).sh

# Tests
unit:
	go test $(TESTARGS) $(TESTFILES)

acceptance:
	go test $(TESTARGS) $(ACCPT_TESTFILES)

test: unit acceptance

check: test

# Coverage
coverage:

# Clean
clean:
	rm -f $(GOPATH)/bin/$(BINARY_NAME)