.PHONY: lint test static install uninstall cross
GOPATH := $(shell go env GOPATH)
VERSION := $(shell git describe --tags --dirty --always)
BIN_DIR := $(GOPATH)/bin
GOX := $(BIN_DIR)/gox
GOLINTER := $(GOPATH)/bin/golangci-lint

$(GOLINTER):
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.49.0

lint: $(GOLINTER)
	$(GOLINTER) run

test:
	go test ./... -bench=. -benchmem

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

# CI
drone:
	drone jsonnet --source .drone/drone.jsonnet --target .drone/drone.yml --stream --format
	drone lint .drone/drone.yml
	drone sign --save grafana/tanka .drone/drone.yml
