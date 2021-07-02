---
name: Directory structure
route: /directory-structure
menu: "References"
---

# Directory structure

Tanka uses the following directories and special files:

```bash
. # the project (<rootDir>)
├── environments # code defining clusters
│   └── default # <baseDir>
│       ├── main.jsonnet # starting point of the Jsonnet compilation
│       └── spec.json # environment's config
├── jsonnetfile.json # direct dependencies
├── jsonnetfile.lock.json # all dependencies with exact versions
├── lib # libraries for this project only
│   └── k.libsonnet # alias file for vendor/github.com/jsonnet-libs/k8s-libsonnet/1.21/main.libsonnet
└── vendor # external libraries installed using jb
    ├── github.com
    │   ├── grafana
    │   │   └── jsonnet-libs
    │   │       └── ksonnet-util # Grafana Labs' usability extensions to k.libsonnet
    │   │           ├── ...
    │   │           └── kausal.libsonnet
    │   └── jsonnet-libs
    │       └── k8s-libsonnet
    │           └── 1.21 # kubernetes library
    │               ├── ...
    │               └── main.libsonnet
    ├── 1.21 -> github.com/jsonnet-libs/k8s-libsonnet/1.21
    └── ksonnet-util -> github.com/grafana/jsonnet-libs/ksonnet-util
```

## Environments

Tanka organizes configuration in environments. For the rationale behind this,
see the [section in the tutorial](/tutorial/environments).

An environment consists of at least two files:

#### spec.json

This file configures environment properties such as cluster connection
(`spec.apiServer`), default namespace (`spec.namespace`), etc.

For the full set of options, see the [Golang source
code](https://github.com/grafana/tanka/blob/main/pkg/spec/v1alpha1/environment.go).

#### main.jsonnet

Like other programming languages, Jsonnet needs an entrypoint into the
evaluation, something to begin with. `main.jsonnet` is exactly this: The very
first file being evaluated, importing or directly specifying everything required
for this specific environment.

## Root and Base

When talking about directories, Tanka uses the following terms:

| Term      | Description                              | Identifier file                   |
| --------- | ---------------------------------------- | --------------------------------- |
| `rootDir` | The root of your project                 | `jsonnetfile.json` or `tkrc.yaml` |
| `baseDir` | The directory of the current environment | `main.jsonnet`                    |

Regardless what subdirectory of the project you are in, Tanka will always be
able to identify both directories, by searching for the identifier files in the
parent directories.  
Tanka needs these for correctly setting up the [import paths](/libraries/import-paths).

This is similar to how `git` always works, by looking for the `.git` directory.

## Libraries

Tanka relies heavily on code-reuse, so libraries are a natural thing. Roughly
spoken, they can be imported from two paths:

- `/lib`: Project local libraries
- `/vendor` External libraries

For more details consider the [import paths](/libraries/import-paths).

### jsonnetfile.json and the lock

`jb` records all external packages installed in a file called
`jsonnetfile.json`. This file is the source of truth about what should be
included in `vendor/`. However, it should only include what is really directly
required, all recursive dependencies will be handled just fine.

`jsonnetfile.lock.json` is generated on every run of jsonnet-bundler, including
a list of packages that must be included in `vendor/`, along with the exact
version and a `sha256` hash of the package contents.

Both files should be checked into source control: The `jsonnetfile.json`
specifies what you need and the `jsonnetfile.lock.json` is important to make
sure that subsequent `jb install` invocations always do the exact same thing.

> **Tip**: The `vendor/` directory can be safely added to `.gitignore` to keep your
> repository size down, as long as `jsonnetfile.lock.json` is checked in.
