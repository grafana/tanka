---
name: Installation
route: /install
---

# Installing Tanka

To install Tanka, it is usually sufficient to install the `tk` binary. It
contains the Jsonnet compiler and everything else required, apart from some
prerequesites:

- [`kubectl`](https://kubernetes.io/docs/tasks/tools/install-kubectl/): Tanka
  uses `kubectl` to communicate to your cluster. This means `kubectl` must be
  available somewhere on your `$PATH`. If you ever have worked with Kubernetes
  before, this should be the case anyways.
- `diff`: To compute differences, standard UNIX `diff(1)` is required.
- (recommended) `jb`: [#Jsonnet-bundler](#jsonnet-bundler), the Jsonnet package
  manager

## Using a package manager (recommended)

We maintain Tanka packages for some operation systems. Installing these is
recommeded, as updates are automatically distributed.

### macOS

Tanka is in
[`Homebrew/homebrew-core`](https://github.com/Homebrew/homebrew-core/blob/master/Formula/tanka.rb),
so you can just install it using `brew`:

```bash
$ brew install tanka
```

### ArchLinux

We maintain two AUR packages, one building [from
source](https://aur.archlinux.org/packages/tanka/) and another one using a
[pre-compiled binary](https://aur.archlinux.org/packages/tanka-bin/). These can
be installed using any AUR helper, e.g. `yay`:

```bash
# from source:
$ yay tanka

# using pre-compiled binary:
$ yay tanka-bin
```

## Precompiled binaries

For all other operating systems, we provide pre-compiled binaries for Tanka at
https://github.com/grafana/tanka/releases.

Just grab the latest version from there, download it and put somewhere in your
`$PATH` (e.g. to `/usr/local/bin/tk`)

## From source

In case the above won't work for you, you can try building the most recent
release using `go get`:

```bash
$ GO111MODULE=on go get github.com/grafana/tanka/cmd/tk
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
