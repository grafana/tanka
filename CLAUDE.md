# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What is Tanka

Tanka is a Kubernetes configuration management tool that uses **Jsonnet** instead of YAML. It provides a structured, reusable way to define and deploy Kubernetes resources. Published by Grafana Labs and used in Grafana Cloud.

## Build and Development Commands

```bash
# Build dev binary (outputs ./cmd/tk)
make dev

# Build static binary (no CGO)
make static

# Install to $GOPATH/bin
make install

# Cross-compile for all platforms (Linux amd64/arm64/arm, Darwin, Windows)
make cross

# Run all tests with benchmarks
make test

# Run a single test
go test ./pkg/some/package -run TestFunctionName

# Run linter (golangci-lint v2.9.0)
make lint

# Run acceptance/e2e tests (requires Dagger)
make acceptance-tests

# Build Docker container
make container
```

Go 1.25+ is required. The version is embedded at compile time via `-X github.com/grafana/tanka/pkg/tanka.CurrentVersion`.

## Architecture Overview

Tanka's core pipeline: **CLI Input → Find/Parse Environments → Evaluate Jsonnet → Process Resources → Apply Filters & Labels → Sort → Apply/Diff/Export to K8s**

### Key Packages

| Package | Role |
|---------|------|
| `cmd/tk/` | CLI entry point; commands: `apply`, `show`, `diff`, `prune`, `delete`, `env`, `export`, `fmt`, `eval`, `init`, `tool` |
| `pkg/tanka/` | Core programmatic API: `Load()`, `Export()`, `Find()`, `FindEnvs()` |
| `pkg/jsonnet/` | Jsonnet evaluation with caching, import resolution, linting |
| `pkg/kubernetes/` | K8s client wrapper; multiple diff strategies (native, server-side, subset) |
| `pkg/process/` | Post-evaluation pipeline: extracts K8s objects, filters, adds `tanka.dev/*` labels, sorts by dependency |
| `pkg/spec/` | Parses `spec.json` environment configs (`v1alpha1/environment.go`) |
| `pkg/helm/` | Helm integration: pull charts, template, native Jsonnet binding (`helm.template()`) |
| `pkg/kustomize/` | Kustomize wrapper with native Jsonnet binding |
| `pkg/term/` | Terminal UI utilities |
| `internal/telemetry/` | OpenTelemetry OTLP tracing (configurable via `OTEL_*` env vars) |

### Environment Configuration

Each Tanka environment has a `spec.json` file defining:
- `apiVersion`, `kind` (Kubernetes-style versioning)
- `metadata.name` — environment name
- `spec.apiServer` — target K8s cluster
- `spec.namespace`, `spec.resourceDefaults`, etc.

### Helm and Kustomize

Both are integrated as native Jsonnet bindings, allowing Helm charts and Kustomize overlays to be referenced directly from Jsonnet code. Helm charts can be vendored for reproducibility.

### External Commands

The CLI supports `tk-*` prefixed external commands discovered on `$PATH`, following a plugin pattern.

## Linting

Configuration is in `.golangci.yml`. Active linters include: `copyloopvar`, `goconst`, `gocritic`, `misspell`, `revive`, `staticcheck`, `unconvert`. Formatters: `gofmt`, `goimports`. The `third_party`, `builtin`, and `examples` directories are excluded.

## Testing

- Unit tests live alongside source files in each package
- Acceptance/e2e tests are in `acceptance-tests/` (separate Go module, run via Dagger)
- CI runs `make lint`, `make test`, and `make cross` on every PR

## Git Commits

Do not add "Co-Authored-By" or any Claude attribution to commit messages.

## Repository Layout Notes

- `acceptance-tests/` and `dagger/` are separate Go modules with their own `go.mod`
- `docs/` contains the user-facing documentation site source
- `pkg/spec/v1alpha1/` defines the versioned environment spec types
