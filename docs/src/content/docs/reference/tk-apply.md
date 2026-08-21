---
title: tk apply
---

Apply the Jsonnet configuration to the connected Kubernetes cluster. Tanka will
show a diff of the changes and prompt for interactive approval before applying.

## Synopsis

```
tk apply <path> [flags]
```

## Examples

```bash
# Apply the default environment
tk apply environments/default

# Apply only a specific resource
tk apply -t deployment/grafana environments/default

# Dry-run (client-side)
tk apply --dry-run=client environments/default

# Apply without interactive approval (for CI)
tk apply --auto-approve=always environments/default
```

## Options

### `--apply-strategy`

`string` — Force the apply strategy to use. Automatically chosen if not set.

### `--auto-approve`

`string` — Skip interactive approval. Only for automation! Allowed values: `always`, `never`, `if-no-changes`.

### `--color`

`string` — Controls color in diff output. Must be `auto`, `always`, or `never`. (default: `auto`)

### `--diff-strategy`

`string` — Force the diff strategy to use. Automatically chosen if not set.

### `--dry-run`

`string` — `--dry-run` parameter passed to kubectl. Must be `none`, `server`, or `client`.

### `--ext-code`

`stringArray` — Set code value of extVar (Format: `key=<code>`).

### `--ext-code-file`

`stringArray` — Set code value of extVar from file (Format: `key=filename`).

### `-V` / `--ext-str`

`stringArray` — Set string value of extVar (Format: `key=value`).

### `--ext-str-file`

`stringArray` — Set string value of extVar from file (Format: `key=filename`).

### `--force`

Force applying (`kubectl apply --force`).

### `-h` / `--help`

Help for apply.

### `--jsonnet-implementation`

Use `go` to use the native go-jsonnet implementation and `binary:<path>` to delegate evaluation to a binary (with the same API as the regular `jsonnet` binary, see the BinaryImplementation docstrings for more details). (default: `go`)

### `--log-level`

`string` — Possible values: `disabled`, `fatal`, `error`, `warn`, `info`, `debug`, `trace`. (default: `info`)

### `--max-stack`

`int` — Jsonnet VM max stack. The default value is the value set in the go-jsonnet library. Increase this if you get: max stack frames exceeded.

### `--name`

`string` — String that only a single inline environment contains in its name.

### `-t` / `--target`

`strings` — Regex filter on `<kind>/<name>`. See [output filtering](/output-filtering/).

### `--tla-code`

`stringArray` — Set code value of top level function (Format: `key=<code>`).

### `--tla-code-file`

`stringArray` — Set code value of top level function from file (Format: `key=filename`).

### `-A` / `--tla-str`

`stringArray` — Set string value of top level function (Format: `key=value`).

### `--tla-str-file`

`stringArray` — Set string value of top level function from file (Format: `key=filename`).

### `--validate`

Validation of resources (`kubectl --validate=false`). (default: `true`)
