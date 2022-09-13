.PHONY: lint test static install uninstall cross
VERSION := $(shell git describe --tags --dirty --always)
BIN_DIR := $(GOPATH)/bin
GOX := $(BIN_DIR)/gox

lint:
	test -z $$(gofmt -s -l cmd/ pkg/ | tee /dev/stderr)
	go vet ./...

test:
	go test ./...

# Compilation
dev:
	go build -ldflags "-X main.Version=dev-${VERSION}" ./cmd/tk

LDFLAGS := '-s -w -extldflags "-static" -X github.com/grafana/tanka/pkg/tanka.CURRENT_VERSION=${VERSION}'
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
