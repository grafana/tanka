# Usage Examples

## Basic Usage

Export a single environment:

```bash
tanka-export ./output ./examples/prom-grafana/environments/prom-grafana/dev
```

## Recursive Export

Export multiple environments recursively:

```bash
tanka-export --recursive ./output ./examples/prom-grafana/environments
```

## Custom Format

Use a custom filename format:

```bash
tanka-export \
  --format "{{.metadata.namespace}}/{{.kind}}/{{.metadata.name}}" \
  --extension yaml \
  ./output \
  ./examples/prom-grafana/environments/prom-grafana/dev
```

## With Merge Strategy

Export to an existing directory, replacing existing files:

```bash
tanka-export \
  --merge-strategy replace-envs \
  ./output \
  --recursive \
  ./examples/prom-grafana/environments
```

## Parallel Processing

Control the number of parallel workers:

```bash
tanka-export \
  --parallel 4 \
  --recursive \
  ./output \
  ./examples/prom-grafana/environments
```

## Filter by Name

Export only environments with a specific name:

```bash
tanka-export \
  --recursive \
  --name dev \
  ./output \
  ./examples/prom-grafana/environments
```

## With Jsonnet Files

The tool uses `jrsonnet` to evaluate Jsonnet files. It will look for `main.jsonnet` in each environment directory:

```bash
# Example environment structure:
# environments/
#   my-env/
#     main.jsonnet      <- Entry point
#     spec.json         <- Environment spec
#     lib/              <- Local libraries
```

## Testing with Tanka Examples

From the repository root:

```bash
# Build the Rust export tool
cd rust-export
cargo build --release

# Export the Tanka example environments
./target/release/tanka-export \
  --recursive \
  /tmp/tanka-output \
  ../examples/prom-grafana/environments
```

## Performance Comparison

The Rust implementation using `jrsonnet` should provide:
- Faster Jsonnet evaluation (native Rust vs Go)
- Better memory usage for large environments
- Parallel processing with Rayon

Benchmark against the Go implementation:

```bash
# Go implementation
time tk export /tmp/output-go \
  --recursive \
  examples/prom-grafana/environments

# Rust implementation
time ./rust-export/target/release/tanka-export \
  --recursive \
  /tmp/output-rust \
  examples/prom-grafana/environments
```

