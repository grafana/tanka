.PHONY: lint test static cross
VERSION := $(shell git describe --tags --dirty --always)
BIN_DIR := $(GOPATH)/bin
GOX := $(BIN_DIR)/gox

lint:
	test -z $$(gofmt -s -l cmd/ pkg/)
	go vet ./...

test:
	go test ./...

# Compilation
dev:
	go build -ldflags "-X main.Version=dev-${VERSION}" ./cmd/tk

LDFLAGS := '-s -w -extldflags "-static" -X main.Version=${VERSION}'
static:
	CGO_ENABLED=0 GOOS=linux go build -ldflags=${LDFLAGS} ./cmd/tk

$(GOX):
	go get -u github.com/mitchellh/gox
cross: $(GOX)
	CGO_ENABLED=0 gox -output="dist/{{.Dir}}_{{.OS}}_{{.Arch}}" -ldflags=${LDFLAGS} -arch="amd64 arm64 arm" -os="linux" -osarch="darwin/amd64" ./cmd/tk

# Docker container
container: static
	docker build -t shorez/tanka .

# CI
drone:
	jsonnet .drone/drone.jsonnet  | jq .drone -r | yq -y . > .drone/drone.yml

changelog:
	standard-version --skip.commit --skip.tag --preMajor

release:
	standard-version --skip.changelog --preMajor
