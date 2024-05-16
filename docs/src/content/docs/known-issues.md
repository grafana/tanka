---
title: Known issues
---

Below is a list of common errors and how to address them.

## `Evaluating jsonnet: RUNTIME ERROR: Undefined external variable: __ksonnet/components`

When migrating from `ksonnet`, this error might occur, because Tanka does not
provide the global `__ksonnet` variable, nor does it strictly have the concept
of components.
You will need to use the plain Jsonnet `import` feature instead. Note that this
requires your code to be inside of one of the
[import paths](./libraries/import-paths).

## `Evaluating jsonnet: RUNTIME ERROR: couldn't open import "k.libsonnet": no match locally or in the Jsonnet library paths`

This error can occur when the `k8s-libsonnet` kubernetes libraries are missing in the
import paths. While `k8s-libsonnet` used to magically include them, Tanka follows a
more explicit approach and requires you to install them using `jb`:

```bash
jb install github.com/jsonnet-libs/k8s-libsonnet/1.21@main
echo "import 'github.com/jsonnet-libs/k8s-libsonnet/1.21/main.libsonnet'" > lib/k.libsonnet
```

This does 2 things:

1. It installs the `k8s-libsonnet` library (in `vendor/github.com/jsonnet-libs/k8s-libsonnet/1.21/`).
   You can replace the `1.21` matching the Kubernetes version you want to run against.

2. It makes an alias for libraries importing `k.libsonnet` directly. See
   [Aliasing](./tutorial/k-lib#aliasing) for the alias rationale.

## Unexpected diff if the same port number is used for UDP and TCP

A [long-standing bug in `kubectl`](https://github.com/kubernetes/kubernetes/issues/39188)
results in an incorrect diff output if the same port number is used multiple
times in differently named ports, which commonly happens if a port is specified
using both protocols, `tcp` and `udp`. Nevertheless, `tk apply` will still work
correctly.
