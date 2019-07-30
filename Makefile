.PHONY: lint test
VERSION := $(shell git describe --tags --dirty --always)

lint:
	test -z $$(gofmt -s -l cmd/ pkg/)
	go vet ./...

test:
	go test ./...

# Compilation
dev:
	go build -ldflags "-X main.Version=dev-${VERSION}" ./cmd/tk

static:
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-s -w -extldflags "-static" -X main.Version=${VERSION}' ./cmd/tk

# Docker container
container: static
	docker build -t shorez/tanka .

# CI
drone:
	jsonnet .drone.jsonnet  | jq .drone -r | yq -y . > .drone.yml
