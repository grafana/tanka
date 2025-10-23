# Tanka Export - Rust Implementation

## Overview

This is a complete Rust rewrite of Tanka's `export` command, using `jrsonnet` as the Jsonnet implementation. The tool evaluates Tanka environments written in Jsonnet and exports them as Kubernetes YAML manifests.

## Architecture

### Core Components

1. **CLI (`cli.rs`)**
   - Uses `clap` for command-line argument parsing
   - Handles all flags and options from the original Go implementation
   - Validates arguments and orchestrates the export process

2. **Environment (`environment.rs`)**
   - Defines Rust data structures mirroring Tanka's `v1alpha1.Environment` spec
   - Includes metadata, spec, and resource defaults
   - Uses `serde` for JSON serialization/deserialization

3. **Jsonnet Evaluation (`jsonnet.rs`)**
   - Uses `jrsonnet-evaluator` (Rust implementation of Jsonnet)
   - Evaluates `.jsonnet` files to JSON
   - Extracts Kubernetes manifests from evaluated data
   - Handles import paths and file resolution

4. **Manifest Handling (`manifest.rs`)**
   - Validates Kubernetes manifest structure
   - Ensures required fields (apiVersion, kind, metadata) are present
   - Converts manifests to YAML format

5. **Template Engine (`template.rs`)**
   - Uses `Tera` template engine for filename generation
   - Supports complex formatting with path separators
   - Handles the BEL character trick for subdirectory creation

6. **Export Logic (`export.rs`)**
   - Main export orchestration
   - Parallel processing using `Rayon`
   - Merge strategies (none, fail-on-conflicts, replace-envs)
   - Manifest tracking via `manifest.json`
   - Environment discovery (single or recursive)

## Key Features

### Jsonnet Evaluation with jrsonnet

- **Native Performance**: `jrsonnet` is a pure Rust implementation, providing fast evaluation
- **Standard Library**: Full Jsonnet standard library support
- **Import Resolution**: Resolves imports relative to the main file
- **Error Handling**: Detailed error messages for Jsonnet evaluation failures

### Parallel Processing

- Uses `Rayon` for parallel iteration over environments
- Configurable parallelism level via `--parallel` flag
- Thread-safe manifest collection with mutex-protected maps

### File Generation

- Template-based filename generation using Tera
- Support for nested directory structures via path separator replacement
- Atomic file writing with directory creation

### Merge Strategies

1. **None** (default): Fails if output directory is not empty
2. **FailOnConflicts**: Allows existing directory but fails on file conflicts
3. **ReplaceEnvs**: Deletes previously exported files for the same environments

### Manifest Tracking

- Generates `manifest.json` to map files to environments
- Enables incremental updates and cleanup of deleted environments
- Optional via `--skip-manifest` flag

## Dependencies

### Core Libraries

- **clap 4.5**: Command-line argument parsing with derive macros
- **tokio 1.48**: Async runtime (though most operations are sync)
- **serde 1.0**: Serialization framework
- **serde_json 1.0**: JSON handling
- **serde_yaml 0.9**: YAML output
- **anyhow 1.0**: Error handling with context
- **thiserror 1.0**: Custom error types

### Jsonnet

- **jrsonnet-evaluator 0.5.0-pre97**: Jsonnet evaluation engine
- **jrsonnet-parser 0.5.0-pre97**: Jsonnet parser
- **jrsonnet-gcmodule 0.3.9**: Garbage collection (pinned to avoid edition2024)

### Utility

- **regex 1.10**: Regular expressions for cache path matching
- **tera 1.20**: Template engine for filename formatting
- **walkdir 2.5**: Recursive directory traversal
- **rayon 1.10**: Data parallelism
- **log 0.4**: Logging facade
- **env_logger 0.11**: Simple logger implementation

## Building

### Requirements

- **Rust Nightly**: Required due to `jrsonnet` dependencies using edition2024 features
- The `rust-toolchain.toml` file specifies `nightly-2024-12-01`

### Build Commands

```bash
# Debug build
cargo build

# Release build (optimized)
cargo build --release

# Check without building
cargo check

# Run tests
cargo test
```

### Binary Location

After building, the binary is located at:
- Debug: `target/debug/tanka-export`
- Release: `target/release/tanka-export`

## Differences from Go Implementation

### Advantages

1. **Performance**: Rust's zero-cost abstractions and native jrsonnet may be faster
2. **Memory Safety**: Compile-time guarantees prevent many classes of bugs
3. **Type Safety**: Strong typing catches errors at compile time
4. **Dependency Management**: Cargo provides reliable dependency resolution

### Limitations

1. **Maturity**: The Go implementation is battle-tested in production
2. **Ecosystem**: Some Tanka-specific features may not be implemented yet
3. **Import Resolution**: jrsonnet's import system differs slightly from Go's implementation
4. **Native Functions**: Custom Jsonnet native functions would need Rust reimplementation

### Not Yet Implemented

- Caching functionality (`--cache-path`, `--cache-envs`)
- Label selectors (`--selector`)
- Resource filtering (`--targets`)
- Memory ballast allocation
- Some advanced Jsonnet features that rely on Go-specific native functions

## Future Improvements

1. **Complete Caching**: Implement evaluation caching for faster repeated exports
2. **Label Selectors**: Full support for Kubernetes label selector syntax
3. **Resource Filtering**: Implement target filtering by kind/name patterns
4. **Custom Import Resolvers**: Support for custom jpath configuration
5. **Progress Reporting**: Add progress bars for long-running operations
6. **Validation**: Add schema validation for manifests
7. **Dry-Run Mode**: Preview changes without writing files
8. **Watch Mode**: Continuous export on file changes

## Testing

### Unit Tests

Run tests with:
```bash
cargo test
```

### Integration Testing

Test with actual Tanka examples:
```bash
# Export the prom-grafana example
./target/release/tanka-export \
  --recursive \
  /tmp/test-output \
  ../examples/prom-grafana/environments

# Verify output
ls -la /tmp/test-output
cat /tmp/test-output/manifest.json
```

### Comparison Testing

Compare output with Go implementation:
```bash
# Go version
tk export /tmp/output-go -r ../examples/prom-grafana/environments

# Rust version  
./target/release/tanka-export -r /tmp/output-rust ../examples/prom-grafana/environments

# Compare
diff -r /tmp/output-go /tmp/output-rust
```

## Performance Benchmarks

TODO: Add benchmarks comparing:
- Jsonnet evaluation speed (jrsonnet vs go-jsonnet)
- Overall export time
- Memory usage
- Parallel scaling

## Contributing

When contributing to this implementation:

1. Maintain compatibility with the Go implementation's CLI interface
2. Follow Rust best practices and idioms
3. Add tests for new functionality
4. Update documentation
5. Keep dependencies minimal and well-maintained

## License

Same as Tanka - Apache License 2.0

