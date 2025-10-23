# Tanka Export (Rust Implementation)

A Rust rewrite of the Tanka `export` command using `jrsonnet` as the Jsonnet implementation.

## Requirements

- **Rust Nightly**: This project requires Rust nightly due to dependencies on `jrsonnet` which uses edition2024 features.
- The `rust-toolchain.toml` file in this directory will automatically use the correct nightly version when building.

## Import Resolution

The implementation properly configures jrsonnet with import paths following the [official jrsonnet CLI approach](https://github.com/CertainLach/jrsonnet/blob/master/cmds/jrsonnet/src/main.rs). The tool:

- Automatically discovers `lib/` and `vendor/` directories by walking up from the environment path
- Finds the project root (where `jsonnetfile.json` exists) 
- Configures `FileImportResolver` with discovered paths
- Includes the standard library via `ContextInitializer`

This matches how the Go implementation handles imports.

## Building

```bash
cargo build --release
```

## Usage

```bash
tanka-export <output_dir> <path> [<path>...] [OPTIONS]

Export Tanka environments to YAML files

Arguments:
  <OUTPUT_DIR>  Output directory for exported manifests
  <PATH>...     Path(s) to Tanka environments

Options:
  --format <FORMAT>
          Format string for output filenames [default: {{.apiVersion}}.{{.kind}}-{{or .metadata.name .metadata.generateName}}]

  --extension <EXTENSION>
          File extension for exported files [default: yaml]

  -p, --parallel <PARALLEL>
          Number of environments to process in parallel [default: 8]

  -c, --cache-path <CACHE_PATH>
          Local file path where cached evaluations should be stored

  -e, --cache-envs <CACHE_ENVS>
          Regexes which define which environment should be cached

  --mem-ballast-size-bytes <MEM_BALLAST_SIZE_BYTES>
          Size of memory ballast to allocate (bytes) [default: 0]

  --merge-strategy <MERGE_STRATEGY>
          What to do when exporting to an existing directory. Values: 'fail-on-conflicts', 'replace-envs'

  --merge-deleted-envs <MERGE_DELETED_ENVS>
          Tanka main files that have been deleted

  --skip-manifest
          Skip generating manifest.json file

  -r, --recursive
          Look recursively for Tanka environments

  --targets <TARGETS>
          Filter by resource kind/name (e.g., deployment/my-app)

  --name <NAME>
          Filter environments by name

  --selector <SELECTOR>
          Label selector to filter environments

  -h, --help
          Print help

  -V, --version
          Print version
```

## Features

- **Jsonnet Evaluation**: Uses `jrsonnet`, a Rust implementation of Jsonnet, for fast native evaluation
- **Parallel Processing**: Exports multiple environments in parallel using Rayon
- **Template-based Filenames**: Supports Tera templates for flexible file naming
- **Merge Strategies**: Handle exporting to existing directories with conflict resolution
- **Manifest Tracking**: Generates `manifest.json` to track exported files

## Implementation Details

### Modules

- `cli.rs`: Command-line argument parsing using clap
- `environment.rs`: Tanka environment data structures
- `export.rs`: Main export logic and orchestration
- `jsonnet.rs`: Jsonnet evaluation using jrsonnet
- `manifest.rs`: Kubernetes manifest handling
- `template.rs`: Filename templating using Tera

### Differences from Go Implementation

- Uses `jrsonnet` instead of the Go jsonnet implementation
- Uses `Rayon` for parallelism instead of Go's goroutines
- Uses `Tera` for templating instead of Go's text/template
- Async runtime with Tokio

## License

Same as Tanka (Apache 2.0)

