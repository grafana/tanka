---
name: Installation
route: /install
---

# Installing Tanka

To install Tanka, it is usually sufficient to install the `tk` binary. It
contains the Jsonnet compiler and everything else required, apart from some
prerequesites:

* [`kubectl`](https://kubernetes.io/docs/tasks/tools/install-kubectl/): Tanka
  uses `kubectl` to communicate to your cluster. This means `kubectl` must be
  available somewhere on your `$PATH`. If you ever have worked with Kubernetes
  before, this should be the case anyways.
* `diff`: To compute differences, standard UNIX `diff(1)` is required.
* (recommended) `jb`: [#Jsonnet-bundler](#jsonnet-bundler), the Jsonnet package
  manager

## Precompiled binaries (recommended)

We provide pre-compiled binaries for Tanka at
https://github.com/grafana/tanka/releases.

Just grab the latest version from there, download it and put somewhere in your
`$PATH` (e.g. to `/usr/local/bin/tk`)

## From source

In case the above won't work for you, try with `go get`:

```bash
$ go get -u github.com/grafana/tanka/cmd/tk
```

If that won't work either, compile by hand:

```bash
$ git clone https://github.com/grafana/tanka
$ cd tanka
$ make install
```

> **Note**: You need a working `go` toolchain for this.

---

## Jsonnet-bundler
Apart from the `tk` binary, you will most probably also want to install
Jsonnet-bundler, the Jsonnet package manager:

```bash
$ go get -u github.com/jsonnet-bundler/jsonnet-bundler/cmd/jb
```
