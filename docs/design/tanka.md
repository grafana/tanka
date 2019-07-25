# tanka, a configuration utility
While `json` and it's superset `yaml` are great configuration languages for
machines, they are [not really good for humans](https://youtu.be/FjdS21McgpE).

Especially in complex environments like Kubernetes, these soon become verbose,
hard to maintain and complex to understand and eventually harmfully duplicated.

#### Existing solutions
String-templating these (`envsubst`, `helm`, etc.) fall short because the
template does not really have any sense of the overall structure.  
Using full programming languages to generate the configuration feels heavy and
is likely to become overwhelming because of the many different ways a problem
could be approached (think of many lines of probably untested code that somehow
mutates dicts).

A solution to this problem is called [jsonnet](https://jsonnet.org): It is
basically `json` but with variables, conditionals, arithmetic, functions,
imports, and error propagation and especially very clever **deep-merging**.

#### Prior art
Especially [ksonnet](https://ksonnet.io) had a big impact on this idea.
While `ksonnet` really proved to be working out,
it was based on a conceptual model consisting of components, prototypes,
environments and much more. While this was fair considering the generator
approach they took, it did not experience much adoption afterwards, possibly
because of the added layer of complexity.

To face this, tanka chooses a much simpler model. It does not handle
dependency resolution of external libraries and it leaves code reuse to the
`import` feature of jsonnet.

This leaves us effectively with **`jsonnet`** and **[Environments](#environments)**

## Code Reuse
When configuration becomes more sophisticated, it inevitably results in
some degree of duplication.

Especially deploying one and the same application to multiple environments
(think `dev` and `prod`, multi-regions, etc.), more or less the same code is
required for each and every of these, just set up a little differently.

Today applications are usually considered building-blocks to compose a larger
system. However, these building blocks need configuration as well, which leads
to another layer of duplication beyond domain boundaries.

By using [`jsonnet`](https://jsonnet.org) as the underlying data templating language, tanka solves both
of these:

### Imports
`jsonnet` is able to
[import](https://jsonnet.org/learning/tutorial.html#imports) other jsonnet
snippets into the current scope, which allows to refactor commonly used code
into shared libraries in turn.

### Bundle
Once a `jsonnet` library becomes required to tackle code-reuse beyond project or
even domain boundaries, it is a common practice in programming languages to have
shared dependencies.

In the `jsonnet` world, this is provided by [`jsonnet-bundler`](https://github.com/jsonnet-bundler/jsonnet-bundler), which maintains a
`vendor` folder (the bundle) with all required libraries ready to be imported,
much like older versions of `go` used to (>1.11).  
This procedure integrates smoothly with tanka, because `vendor` is on the [`JPATH`](#jpath).

### `JPATH`
To wrap this all up, tanka sets up the `jsonnet` dependency resolution in a
special way:

| Name             | Identifier         | Description                                                                                                                           |
|------------------|--------------------|---------------------------------------------------------------------------------------------------------------------------------------|
| `rootDir`        | `jsonnetfile.json` | Every file in the tree of this folder is considered part of the project. Much like `git` has the one directory with the `.git` folder |
| `rootDir/vendor` |                    | Populated with shared dependencies by `jsonnet-bundler`                                                                               |
| `baseDir/lib`    |                    | Code that is only re-used in-tree may be put here                                                                                     |
| `baseDir`        | `main.jsonnet`     | Environment specific code can be put here. The special `main.jsonnet` is the entry point for the evaluation                           |

To resolve the `JPATH`, tanka first traverses the directory tree *upwards*, to
find a `jsonnetfile.json`, which marks the `rootDir`. Reaching `/` without a
match will result in an error.

Same applies for the `baseDir`, the tree is traversed *upwards* for a
`main.jsonnet`.

The final `JPATH` looks like the following:
```
<baseDir>
<rootDir>/lib
<rootDir>/vendor
```

Imports are relative to the `JPATH`. The earlier a directory appears in the
`JPATH`, the higher it's precedence is.

### Directory structure
In a simple setup, it is fair to have the same `rootDir` and `baseDir`:
```
.
├── jsonnetfile.json
├── lib/
├── main.jsonnet
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
└── vendor/
```

While latter one is the structure suggested by `tk init`, it is perfectly fine to
use another if it fits the use-case better. The folder does not need to be named
`environments`, either.

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
  "api_version": "tanka.dev/v1alpha1",
  "kind": "Environment",
  "metadata": {
    "name": "auto"
    "labels": {}
  },
  "spec": {}
}
```
| Field             | Description                                                         |
|-------------------|---------------------------------------------------------------------|
| `api_version`     | marks the version of the API, to allow schema changes once required |
| `kind`            | not used yet, added for completeness                                |
| `metadata.name`   | automatically set to the directory name                             |
| `metadata.labels` | use to add descriptive `key:value` pairs. Not evaluated             |
| `spec`            | [Provider](#providers) configuration                                 |

The environment object is accessible from within `jsonnet`.

## `show`, `diff`, `apply`, (repeat)
This is the main workflow of tanka: Once changes to the configuration are done,
`tk show` is used to validate that jsonnet has been correctly translated into
the target format (`json`, `yaml`, etc.)

Once this is the case, `tk diff` is used to obtain a detailed overview of how
the desired state diverts from the current state of the system and which actions
will be taken.

When the output of `diff` is satisfying, the desired state can become reality using
`tk apply`.

## Providers
Because there are a variety of systems out there that are configured using
json compatible interfaces and may have special needs, tanka is unable to
fulfill these by it's own.

While `tk eval` uses emits the raw output of the jsonnet compiler, the more
sophisticated commands `show`, `diff` and `apply` are actually fulfilled by a
so-called provider, (pluggable) pieces of logic which handle platform specific
actions.

This is best shown using an example, the *Kubernetes* provider:

| Action  | Description                                                                                                                                                                                                        |
|---------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `show`  | Once the `jsonnet` is evaluated, the provider has to reconcile it for the target system. The `jsonnet` output is a deeply nested dict of Kubernetes objects, while `kubectl` requires multi-document yaml strings. |
| `diff`  | The same reconciling happens, but instead of printing the output, it is passed to `kubectl diff -f -`. The actual `diff` is printed afterwards.                                                                    |
| `apply` | After reconciling, the output is passed to `kubectl apply -f -`                                                                                                                                                    |
