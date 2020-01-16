---
name: Directory structure
route: /directory-structure
---

# Directory structure

Tanka uses the following directories and special files:

```bash
. # the project
├── environments # code defining clusters
│   └── default
│       ├── main.jsonnet # starting point of the Jsonnet compilation
│       └── spec.json # environment's config
├── jsonnetfile.json # direct dependencies
├── jsonnetfile.lock.json # all dependencies with exact versions
├── lib # libraries for this project only
│   └── k.libsonnet # alias file for vendor/ksonnet.beta.4/k.libsonnet
└── vendor # external libraries installed using jb
    ├── ksonnet.beta.4 # kubernetes library
    │   ├── k8s.libsonnet
    │   └── k.libsonnet
    └── ksonnet-util # Grafana Labs' usability extensions to k.libsonnet
        └── kausal.libsonnet
```

## Environments
Tanka organizes configuration in environments. For the rationale behind this,
see the [section in the tutorial](/tutorial/environments).

An environment consists of at least two files:

#### spec.json
This file configures environment properties such as cluster connection
(`spec.apiServer`), default namespace (`spec.namespace`), etc.

For the full set of options, see the [Golang source
code](https://github.com/grafana/tanka/blob/master/pkg/spec/v1alpha1/config.go).

#### main.jsonnet
Like other programming languages, Jsonnet needs an entrypoint into the
evaluation, something to begin with. `main.jsonnet` is exactly this: The very
first file being evaluated, importing or directly specifying everything required
for this specific environment.


## Libraries
Tanka builds on code-reuse, by refactoring common pieces into libraries, which can be imported from two locations:

### lib
The `lib/` folder is for libraries that are meant for only this single project.
If you intend to deploy your custom e-commerce stack, you could for example have
libraries for the `auth`, `bookings`, `billing` and `inventory` here.

They are not intended to be shared and thus are a good fit for `lib/`

> **Note:** Opposing to `vendor/`, `lib/` is entirely your realm. You manage the
> contents and Tanka won't ever mess with this after `tk init`.

### vendor
Some libraries can be useful to many projects (for example ones for
[Prometheus](https://prometheus.io), [Loki](https://grafana.com/loki), etc).

These are usually published on GitHub. To use them in your project, [install
them using `jb`](/libraries/install-publish#install-a-library). This will store
a copy of the source code on the remote in the `vendor/` directory. Note that
this folder belongs to `jb` and all files not recorded in
`jsonnetfile.lock.json` will be removed on the next run. Also don't edit files
in here (your changes would be reverted anyways), use
[Shadowing](/libraries/import-paths#shadowing) instead.

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
