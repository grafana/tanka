.PHONY: all lint test static install uninstall cross acceptance-tests dagger-develop
GOPATH := $(shell go env GOPATH)
VERSION := $(shell git describe --tags --dirty --always)
BIN_DIR := $(GOPATH)/bin
GOX := $(BIN_DIR)/gox
GOLINTER := $(GOPATH)/bin/golangci-lint

all: lint test static

$(GOLINTER):
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.9.0

lint: $(GOLINTER)
	$(GOLINTER) run

test:
	go test ./... -bench=. -benchmem

acceptance-tests:
	dagger call acceptance-tests --root-dir . --acceptance-tests-dir ./acceptance-tests

# Compilation
dev:
	go build -ldflags "-X main.Version=dev-${VERSION}" ./cmd/tk

LDFLAGS := '-s -w -extldflags "-static" -X github.com/grafana/tanka/pkg/tanka.CurrentVersion=${VERSION}'
static:
	CGO_ENABLED=0 go build -ldflags=${LDFLAGS} ./cmd/tk

install:
	CGO_ENABLED=0 go install -ldflags=${LDFLAGS} ./cmd/tk

uninstall:
	go clean -i ./cmd/tk

$(GOX):
	go get -u github.com/mitchellh/gox
	go install github.com/mitchellh/gox

cross: $(GOX)
	CGO_ENABLED=0 $(BIN_DIR)/gox -output="dist/{{.Dir}}-{{.OS}}-{{.Arch}}" -ldflags=${LDFLAGS} -arch="amd64 arm64 arm" -os="linux" -osarch="darwin/amd64" -osarch="darwin/arm64" -osarch="windows/amd64" ./cmd/tk

# Docker container
container: static
	docker build -t grafana/tanka .

dagger-develop:
	@cp dagger/.gitignore dagger/.gitignore.bak
	@dagger develop --silent
	@mv dagger/.gitignore.bak dagger/.gitignore
