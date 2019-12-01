---
name: "Directory structure"
route: "/environments/structure"
menu: "Environments"
---

# Directory structure

Tanka works with a minimal set of assumptions on its environment. Unlike tools
such as `ksonnet`, it does expect code to be placed in strictly named
directories.

Instead, it relies on two directories to find your code:

### Root Directory (`rootDir`)

The `rootDir` marks the start of a directory tree that represents a Tanka
project. It behaves similar to a `git` repository:  
A marker file
([`jsonnetfile.json`](https://github.com/jsonnet-bundler/jsonnet-bundler), like
git has its `.git/`) indicates the beginning of the tree. Regardless of how deep
you are in it, Tanka can always discover the project by searching for a
`jsonnetfile.json` in the parent directories.

### Base Directory (`baseDir`)

The base directory is the directory that contains a file called `main.jsonnet`.
This file is used as the entrypoint for evaluating `.jsonnet` to `.yaml`.  
The `baseDir` _must_ be in the tree of the `rootDir` or the `rootDir` itself
(i.e. `jsonnetfile.json` and `main.jsonnet` in the same directory).

> Example: Minimal.  
> In this example, `rootDir` and `baseDir` are the same.
>
> ```tree
> .
> ├── jsonnetfile.json
> ├── lib/
> ├── main.jsonnet
> └── vendor/`
> ```

> Example: Environments.  
> To enable a behavior close to what `ksonnet` used to call an _Environment_,
> multiple `baseDirs` can be created in sub-directories.
>
> ```tree
> .
> ├── environments
> │   ├── dev
> │   │   └── main.jsonnet
> │   └── prod
> │       └── main.jsonnet
> ├── jsonnetfile.json
> ├── lib/
> └── vendor/
> ```

## Import paths

There are three places, imported files can come from:

#### Relative

Imports may be relative to the current file. Use a relative path
(`import "./whatever.jsonnet"`) for this.

#### Library (`rootDir/lib`)

If a folder called `lib` exists in `rootDir`, imports from here are possible as
well.

Place code here, that is used multiple times across this project. However, when
this code needs to be re-used across project boundaries, consider moving it into
its own Git repository and vendor it in using
[Jsonnet bundler](https://github.com/jsonnet-bundler/jsonnet-bundler).

#### Vendor (`rootDir/vendor`)

The purpose `vendor` folder is to hold _shared_ libraries, that are downloaded
using a package manager.

> Warning: This folder shall be managed by
> [Jsonnet bundler](https://github.com/jsonnet-bundler/jsonnet-bundler) or
> another comparable tool. **Do not** modify the files in here by hand. Change
> them on the remote if required.

### Precedence

The most specific import takes precedence:

1. Relative import
2. Local library (`lib/`)
3. Shared library (`vendor/`)

The higher an import appears in this list, the more likely it is taken.
