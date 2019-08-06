# tanka, advanced Kubernetes configuration
**Author**: @sh0rez  
**Date**: 26.07.2019

## Motivation
`json` and it's superset `yaml` are great configuration languages for
machines, but they are [not really good for humans](https://youtu.be/FjdS21McgpE).

But the Kubernetes ecosystem seems to have settled to configuring all kinds of
workloads using exactly these. While they provide a low entry barrier, it is
challenging to model advanced requirements, for example deploying the same
application to multiple clusters or to different environments.

This usually leads to duplication on two levels: The obvious one is the
application level. Once an app needs to redeployed to another environment
(`dev`, `prod`, etc.), it usually shares most of the configuration, except for
some edge-cases (secrets, etc). `kubectl` does not really provide a utility to
fully address this.

Even more duplication happens on the systems level. Today, multiple applications
are usually composed into a larger picture.  
But deploying common building blocks as `postgresql` or `nginx` requires the
same code (`Deployment`, `Service`, etc) that nearly everyone attempting to use
these will write.

While maintaining the same code for `dev` and `prod` might go well, it becomes
unconquerable once it comes to multi-region (5+) deployments or even more
versatile use-cases.

#### Existing solutions
Historically, this problem has been approached using string-templating (`helm`).
While `helm` allows code-reuse using `values.yml`, a chart is only able to provide what
has been thought of during writing. Even worse, charts are maintained inside of
the `helm/charts` repository, which makes quick or even domain-specific edits
hard. It is impossible to easily address edge-cases.

A solution to this problem is called [jsonnet](https://jsonnet.org): It is
basically `json` but with variables, conditionals, arithmetic, functions,
**imports**, and error propagation and especially very clever [**deep-merging**](#edge-cases).

#### Prior art
Especially [ksonnet](https://ksonnet.io) had a big impact on this idea.
While `ksonnet` really proved the idea to be working, is was based on the
concept of components, building blocks consisting of prototypes and parameters,
which may be composed into applications or modules which in turn may be applied
to multiple environments.

We believe such a concept overcomplicates the immediate goal of reducing
duplication while allowing edge-cases, because it handles these on a higher
conceptual level, instead of reusing the native capabilities of jsonnet.

Code-reuse and composability is
adequately provided by the native `import` feature of `jsonnet`. Sharing code
beyond application boundaries is already enabled by 
[`jsonnet-bundler`](https://github.com/jsonnet-bundler/jsonnet-bundler).

While it is possible to mimic environments with native `jsonnet`, it falls short
when it comes to actually applying it to the correct cluster. The raw `json`
being returned by the compiler needs reconciling and it must be made sure it is
applied to the correct cluster.

This leaves us effectively with **`jsonnet`** and **[Environments](#environments)**

## Code Reuse
By using [`jsonnet`](https://jsonnet.org) as the underlying data templating
language, tanka supports dynamic reusing of code, just like real programming languages do:

### Imports
`jsonnet` is able to
[import](https://jsonnet.org/learning/tutorial.html#imports) other jsonnet
snippets into the current scope, which in turn allows to refactor commonly used code
into shared libraries.

### Bundle
Once a library becomes general enough to be used beyond project or
even domain boundaries, it is a common practice in programming languages to have
shared dependencies.

In the `jsonnet` world, this is provided by
[`jsonnet-bundler`](https://github.com/jsonnet-bundler/jsonnet-bundler), which maintains a
`vendor` folder (the bundle) with all required libraries ready to be imported,
much like older versions of `go` used to (>1.11) do.  
This procedure integrates smoothly with tanka, because `vendor` is on the [`JPATH`](#jpath).

### `JPATH`
To enable a predictable developer experience, tanka uses clear rules to define
how importing works.

Imports are relative to the `JPATH`. The earlier a directory appears in the
`JPATH`, the higher it's precedence is.

To set it up, tanka makes use of the following directoies:

| Name             | Identifier         | Description                                                                                                                           |
|------------------|--------------------|---------------------------------------------------------------------------------------------------------------------------------------|
| `rootDir`        | `jsonnetfile.json` | Every file in the tree of this folder is considered part of the project. Much like `git` has the one directory with the `.git` folder |
| `rootDir/vendor` |                    | Populated with shared dependencies by `jsonnet-bundler`                                                                               |
| `baseDir/units`  |                    | Code that directly outputs Kubernetes objects (mixins, etc.)
| `baseDir/lib`    |                    | Local helper utilities and other shared code
| `baseDir`        | `main.jsonnet`     | Environment specific code can be put here. The special `main.jsonnet` is the entry point for the evaluation                           |

To resolve the `JPATH`, tanka first traverses the directory tree *upwards*, to
find a `jsonnetfile.json`, which marks the `rootDir`. Reaching `/` without a
match will result in an error.  
This is required to be able to resolve the `JPATH` regardless of how deep one is
inside of the directory tree. Think of it as a root marker, like git has its `.git` folder.  
Even if `jb` is not used, it barely harms to have an unused file with `{}` in it around.

Same applies for the `baseDir`, the tree is traversed *upwards* for a
`main.jsonnet`.

The final `JPATH` looks like the following:
```
<baseDir>
<rootDir>/units
<rootDir>/lib
<rootDir>/vendor
```


### Directory structure
In a simple setup, it is fair to have the same `rootDir` and `baseDir`:
```
.
├── jsonnetfile.json
├── lib/
├── main.jsonnet
├── units/
└── vendor/
```

However, to use [Environments](#environments), the `baseDir` could be moved
into subdirectories:
```
.
├── environments
│   ├── dev
│   │   └── main.jsonnet
│   └── prod
│       └── main.jsonnet
├── jsonnetfile.json
├── lib/
├── units/
└── vendor/
```

While latter structure is the one suggested by `tk init`, it is perfectly fine to
use another if it fits the use-case better. The folder does not need to be named
`environments`, either.

## Edge Cases
During development of `jsonnet` libraries, e.g. for applications like `mysql`,
it is impossible to think of every edge-case in before.

But thanks to the power of `jsonnet`, this is not a problem. Imagine the output
of the library being the following:
```jsonnet
local out = {
  apiVersion: "v1",
  kind: "namespace",
  metadata: {
    name: "production"
  }
};
```

When you wanted to label the namespace, but the library did not provide
functions for it, you could use deep-merging:
```jsonnet
out + {
  metadata+: {
    labels: {
      foo: "bar"
    }
  }
}
```

Note the special `+:` to enable deep merging.

This would result in the second dict being recursively merged on top of the first one:
```jsonnet
{
  apiVersion: "v1",
  kind: "namespace",
  metadata: {
    name: "production",
    labels: {
      foo: "bar"
    }
  }
}
```

## Environments
The only core concept of tanka is an `Environment`. It describes a single
context that can be configured. Such a context might be `dev` and `prod`
environments, Blue/Green deployments or the same application available
in multiple regional zones / datacenters.

An environment does not need to be created, it rather just exists as soon as a
`main.jsonnet` is found somewhere in the tree of a `rootDir`.

However, an environment may receive additional configuration, by adding a file
called `spec.json` alongside the `main.jsonnet`:
```json
{
  "apiVersion": "tanka.dev/v1alpha1",
  "kind": "Environment",
  "metadata": {
    "name": "auto",
    "labels": {}
  },
  "spec": {
    "apiServer": "https://localhost:6443",
    "namespace": "default"
  }
}
```

| Field             | Description                                                         |
|-------------------|---------------------------------------------------------------------|
| `apiVersion`     | marks the version of the API, to allow schema changes once required |
| `kind`            | not used yet, added for completeness                                |
| `metadata.name`   | automatically set to the directory name                             |
| `metadata.labels` | descriptive `key:value` pairs                                      |
| `spec.apiServer`  | The Kubernetes endpoint to use                                      |
| `spec.namespace`  | All objects will be forced into this namespace                      |

The environment object is accessible from within `jsonnet`.

