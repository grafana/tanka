# Changelog

## 0.8.0 (2020-02-13)

The next big one is here! Feature packed with environment overriding and `tk export`. Furthermore lots of bugs were fixed, so using Tanka should be much
smoother now!

#### Highlight: Overriding `vendor` per Environment **([#198](https://github.com/grafana/tanka/pull/198))**

It is now possible, to have a `vendor/` directory managed by `jb` on an
environment basis: https://tanka.dev/libraries/overriding. This means you can
test out changes in libraries in single environments (like `dev`), without
affecting others (like `prod`).

#### Notice:

Changes done in the last release (v0.7.1) can cause indentation changes when
using `std.manifestYAMLFromJSON()`, related to bumping `gopkg.in/yaml.v2` to
`gopkg.in/yaml.v3`.  
Please encourage all your teammembers to upgrade to at least v0.7.1 to avoid
whitespace-only diffs on your projects.

### Features

- **cli**: `tk export` can be used to write all generated Kubernetes resources
  to `.yaml` files

### Bug Fixes

- **kubernetes**: Fail on `diff` when `kubectl` had an internal error
  **([#213](https://github.com/grafana/tanka/pull/213))**
- **kubernetes**: Stop injecting namespaces into wrong places:  
  Tanka was injecting the default namespace into resources of all kinds,
  regardless of whether they actually took one. This caused errors, so we
  stopped doing this. From now on, the default namespace will only be injected
  when the resource is actually namespaced.
  **([#208](https://github.com/grafana/tanka/pull/208))**

* **cli**: `tk diff` colors:  
  Before, the coloring was unstable when scrolling up and down. We fixed this by
  pressing CAPS-LOCK.  
  Furthermore, the output of `tk diff` now also works on light color schemes,
  without messing up the background color.
  **([#210](https://github.com/grafana/tanka/pull/210))**
* **cli**: Proper `--version` output:  
  The release binaries now show the real semver on `tk --version`, instead of
  the git commit sha. **([#201](https://github.com/grafana/tanka/pull/201))**
* **cli**: Print diff on apply again:  
  While refactoring, we accidentally forgot to dereference a pointer, so that
  `tk apply` showed a memory address instead of the actual differences, which
  was kinda pointless. **([#200](https://github.com/grafana/tanka/pull/200))**

## 0.7.1 (2020-02-06)

This is a smaller release focused on critical bug fixes and some other minor
enhancements. While features are included, none of them are significant, meaning
they are part of a patch release.

#### Critical: `parseYaml` works now

Before, `std.native('parseYaml')` did not work at all, a line of code got lost
during merge/rebase, resulting in `parseYaml` returning invalid data, that
Jsonnet could not process. This issue has been fixed in
**([#195](https://github.com/grafana/tanka/pull/195))**.

#### Jsonnet update

The built-in Jsonnet compiler has been upgraded to the lastest master
[`07fa4c0`](https://github.com/google/go-jsonnet/commit/07fa4c037b4ff8b5e601546cb5de4abecaf2651d).
In some cases, this should provide up to 50% more speed, especially when
`base64` is involved, which is now natively implemented.
**([#196](https://github.com/grafana/tanka/pull/196))**

### Features

- **cli**: `tk env set|add` has been extended by `--server-from-context`, which
  allows to parse `$KUBECONFIG` to find the apiServer's IP directly from that
  file, instead of having to manually specify it by hand.
  **([#184](https://github.com/grafana/tanka/pull/184))**
- **jsonnet**: `vendor` overrides:  
  It is now possible to have a `vendor/` directory per environment, so that
  updating upstream libraries can be done gradually.
  **([#185](https://github.com/grafana/tanka/pull/185))**
- **kubernetes**: disable `kubectl` validation:
  `tk apply` now takes `--validate=false` to pass that exact flag to `kubectl`
  as well, for disabling the integrated schema validation.
  **([#186](https://github.com/grafana/tanka/pull/186))**

### Bug Fixes

- **jsonnet, cli**: Stable environment name: The value of `(import "tk").env.name`
  does not anymore depend on how Tanka was invoked, but will
  always be the relative path from `<rootDir>` to the environment's directory.
  **([#182](https://github.com/grafana/tanka/pull/182))**
- **jsonnet**: The nativeFunc `parseYaml` has been fixed to actually return a
  valid result **([#195](https://github.com/grafana/tanka/pull/195))**

## 0.7.0 (2020-01-21)

The promised big update is here! In the last couple of weeks a lot has happened.

Grafana Labs [announced Tanka to the
public](https://grafana.com/blog/2020/01/09/introducing-tanka-our-way-of-deploying-to-kubernetes/),
and the project got a lot of positive feedback, shown both on HackerNews and in
a 500+ increase in GitHub stars!

While we do not ship big new features this time, we ironed out many annoyances
and made the overall experience a lot better:

#### Better website + tutorial ([#134](https://github.com/grafana/tanka/pull/134))

Our [new website](https://tanka.dev) is published! It does not only look super
sleek and performs like a supercar, we also revisited (and rewrote) the most of
the content, to provide especially new users a good experience.

This especially includes the **[new
tutorial](https://tanka.dev/tutorial/overview)**, which gives new and probably
even more experienced users a good insight into how Tanka is meant to be used.

#### :rotating_light::rotating_light: Disabling `import ".yaml"` ([#176](https://github.com/grafana/tanka/pull/176)) :rotating_light::rotating_light:

Unfortunately, we **had to disable the feature** that allowed to directly import
YAML files using the familiar `import` syntax, introduced in v0.6.0, because it
caused serious issues with `importstr`, which became unusable.

While our extensions to the Jsonnet language are cool, it is a no-brainer that
compatibility with upstream Jsonnet is more important. We will work with the
maintainers of Jsonnet to find a solution to enable both, `importstr` and
`import ".yaml"`

**Workaround:**

```diff
- import "foo.yaml"
+ std.parseYaml(importstr "foo.yaml")
```

#### `k.libsonnet` installation ([#140](https://github.com/grafana/tanka/pull/140))

Previously, installing `k.libsonnet` was no fun. While the library is required
for nearly every Tanka project, it was not possible to install it properly using
`jb`, manual work was required.

From now on, **Tanka automatically takes care of this**. A regular `tk init`
installs everything you need. In case you prefer another solution, disable this
new thing using `tk init --k8s=false`.

### Features

- **cli**, **kubernetes**: `k.libsonnet` is now automatically installed on `tk init` **([#140](https://github.com/grafana/tanka/pull/140))**:  
  Before, installing `k.libsonnet` was a time consuming manual task. Tanka now
  takes care of this, as long as `jb` is present on the `$PATH`. See
  https://tanka.dev/tutorial/k-lib#klibsonnet for more details.
- **cli**: `tk env --server-from-context`:  
  This new flag allows to infer the cluster IP from an already set up `kubectl`
  context. No need to remember IP's anymore – and they are even autocompleted on
  the shell. **([#145](https://github.com/grafana/tanka/pull/145))**
- **cli**, **jsonnet**: extCode, extVar:  
  `-e` / `--extCode` and `--extVar` allow using `std.extVar()` in Tanka as well.
  In general, `-e` is the flag to use, because it correctly handles all Jsonnet
  types (string, int, bool). Strings need quoting!
  **([#178](https://github.com/grafana/tanka/pull/178))**

* **jsonnet**: The contents of `spec.json` are now accessible from Jsonnet using
  `(import "tk").env`. **([#163](https://github.com/grafana/tanka/pull/163))**
* **jsonnet**: Lists (`[ ]`) are now fully supported, at an arbitrary level of
  nesting! **([#166](https://github.com/grafana/tanka/pull/166))**

### Bug Fixes

- **jsonnet**: `nil` values are ignored from the output. This allows to disable
  objects using the `if ... then {}` pattern, which returns nil if `false`
  **([#162](https://github.com/grafana/tanka/pull/162))**.
- **cli**: `-t` / `--target` is now case-insensitive
  **([#130](https://github.com/grafana/tanka/pull/130))**

---

## 0.6.1 (2020-01-06)

First release of the new year! This one is a quick patch that lived on master
for some time, fixing an issue with the recent "missing namespaces" enhancement
leading to `apply` being impossible when no namespace is included in Jsonnet.

More to come soon :D

---

## 0.6.0 (2019-11-27)

It has been quite some time since the last release during which Tanka has become
much more mature, especially regarding the code quality and structure.

Furthermore, Tanka has just hit the 100 Stars :tada:

Notable changes include:

#### API ([#97](https://github.com/grafana/tanka/commit/c5edb8b0153ef991765f2f555c839b0f9a487e75))

The most notable change is probably the **Go API**, available at
`https://godoc.org/github.com/grafana/tanka/pkg/tanka`, which allows to use all
features of Tanka directly from any other Golang application, without needing to
exec the binary. The API is inspired by the command line parameters and should
feel very similar.

#### Importing YAML ([#106](https://github.com/grafana/tanka/commit/8029efa44461b5f7ba83a218ccc45bd758c8a322))

It is now possible to import `.yaml` documents directly from Jsonnet. Just use
the familiar syntax `import "foo.yaml"` like you would with JSON.

#### Missing Namespaces ([#120](https://github.com/grafana/tanka/commit/3b9fac1563a75a571b512887602eb53f82e565bf))

Tanka now handles namespaces that are not yet created, in a more user friendly
way than `kubectl\*\* does natively.  
During diff, all objects of an in-existent namespace are shown as new and when
applying, namespaces are applied first to allow applying in a single step.

### Features

- **tool/imports**: import analysis using upstream jsonnet: Due to recent
  changes to google/jsonnet, we can now use the upstream compiler for static
  import analysis
  ([#84](https://github.com/grafana/tanka/commit/394cb12b28beb0ea05d065594b6cf5c3f92de5e4))
- **Array output**: The output of Jsonnet may now be an array of Manifests.
  Nested arrays are not supported yet.
  ([#112](https://github.com/grafana/tanka/commit/eb647793ff5515bc828e4f91186655c143bb6a04))

### Bug Fixes

- **Command Usage Guidelines**: Tanka now uses the [command description
  syntax](https://en.wikipedia.org/wiki/Command-line_interface#Command_description_syntax)
  ([#94](https://github.com/grafana/tanka/commit/13238e5941bd6e68f410d3938d1a285224c2f91d))
- **cli/env** resolved panic on missing `spec.json`
  ([#108](https://github.com/grafana/tanka/commit/9bd15e6b4226164efe45f50c9ed41c4a5673ea2d))

---

## 0.5.0 (2019-09-20)

This version adds a set of commands to manipulate environments (`tk env add, rm, set, list`) ([#73](https://github.com/grafana/tanka/pull/73)). The commands are
mostly `ks env` compatible, allowing `tk env` be used as a drop-in replacement
in scripts.

Furthermore, an error message has been improved, to make sure users can
differentiate between parse issues in `.jsonnet` and `spec.json`
([#71](https://github.com/grafana/tanka/pull/71)).

---

## 0.4.0 (2019-09-06)

After nearly a month, the next feature packed release of Tanka is ready!
Highlights include the new documentation website https://tanka.dev, regular
expression support for targets, diff histograms and several bug-fixes.

### Features

- **cli**: `tk show` now aborts by default, when invoked in a non-interactive
  session. Use `--dangerous-allow-redirect` to disable this safe-guard
  ([#47](https://github.com/grafana/tanka/issues/47)).
- **kubernetes**: Regexp Targets: It is now possible to use regular expressions
  when specifying the targets using `--target` / `-t`. Use it to easily select
  multiple objects at once: https://tanka.dev/targets/#regular-expressions
  ([#64](https://github.com/grafana/tanka/issues/64)).
- **kubernetes**: Diff histogram: Tanka now allows to summarize the differences
  between the live configuration and the local one, by using the unix
  `diffstat(1)` utility. Gain a sneek peek at a change using `tk diff -s .`!
  ([#67](https://github.com/grafana/tanka/issues/67))

### Bug Fixes

- **kubernetes**: Tanka does not fail anymore, when the configuration file
  `spec.json` is missing from an Environment. While you cannot apply or diff,
  the show operation works totally fine
  ([#56](https://github.com/grafana/tanka/issues/56),
  [#63](https://github.com/grafana/tanka/issues/63)).
- **kubernetes**: Errors from `kubectl` are now correctly passed to the user
  ([#61](https://github.com/grafana/tanka/issues/61)).
- **cli**: `tk diff` does not output useless empty lines (`\n`) anymore
  ([#62](https://github.com/grafana/tanka/issues/62)).

---

## 0.3.0 (2019-08-13)

Tanka v0.3.0 is here!

This version includes lots of tiny fixes and detail improvements, to make it easier for everyone to configure their Kubernetes clusters.

Enjoy target support, enhancements to the diff UX and an improved CLI experience.

### Features

The most important feature is **target support** ([#30](https://github.com/tbraack/tanka/issues/30)) ([caf205a](https://github.com/tbraack/tanka/commit/caf205a)): Using `--target=kind/name`, you can limit your working set to a subset of the objects, e.g. to do a staged rollout.

There where some other features added:

- **cli:** autoApprove, forceApply ([#35](https://github.com/tbraack/tanka/issues/35)) ([626b097](https://github.com/tbraack/tanka/commit/626b097)): allows to skip the interactive verification. Furthermore, `kubectl` can now be invoked with `--force`.
- **cli:** print deprecated warnings in verbose mode. ([#39](https://github.com/tbraack/tanka/issues/39)) ([6de170d](https://github.com/tbraack/tanka/commit/6de170d)): Warnings about the deprecated configs are only printed in verbose mode
- **kubernetes:** add namespace to apply preamble ([#23](https://github.com/tbraack/tanka/issues/23)) ([9e2d927](https://github.com/tbraack/tanka/commit/9e2d927)): The interactive verification now shows the `metadata.namespace` as well.
- **cli:** diff UX enhancements ([#34](https://github.com/tbraack/tanka/issues/34)) ([7602a19](https://github.com/tbraack/tanka/commit/7602a19)): The user experience of the `tk diff` subcommand has been improved:
  - if the output is too long to fit on a single screen, the systems `PAGER` is invoked
  - if differences are found, the exit status is set to `16`.
  - When `tk apply` is invoked, the diff is shown again, to make sure you apply what you want

### Bug Fixes

- **cli:** invalid command being executed twice ([#42](https://github.com/tbraack/tanka/issues/42)) ([28c6898](https://github.com/tbraack/tanka/commit/28c6898)): When the command failed, it was executed twice, due to an error in the error handling of the CLI.
- **cli**: config miss ([#22](https://github.com/tbraack/tanka/issues/22)) ([32bc8a4](https://github.com/tbraack/tanka/commit/32bc8a4)): It was not possible to use the new configuration format, due to an error in the config parsing.
- **cli:** remove datetime from log ([#24](https://github.com/tbraack/tanka/issues/24)) ([1e37b20](https://github.com/tbraack/tanka/commit/1e37b20))
- **kubernetes:** correct diff type on 1.13 ([#31](https://github.com/tbraack/tanka/issues/31)) ([574f946](https://github.com/tbraack/tanka/commit/574f946)): On kubernetes 1.13.0, `subset` was used, although `native` is already supported.
- **kubernetes:** Nil pointer deference in subset diff. ([#36](https://github.com/tbraack/tanka/issues/36)) ([f53c2b5](https://github.com/tbraack/tanka/commit/f53c2b5))
- **kubernetes:** sort during reconcile ([#33](https://github.com/tbraack/tanka/issues/33)) ([ab9c43a](https://github.com/tbraack/tanka/commit/ab9c43a)): The output of the reconcilation phase is now stable in ordering

---

## [0.2.0](https://github.com/tbraack/tanka/compare/v0.1.0...v0.2.0) (2019-08-07)

### Features

- **cli:** Completions ([#7](https://github.com/tbraack/tanka/issues/7)) ([aea3bdf](https://github.com/tbraack/tanka/commit/aea3bdf)): Tanka is now able auto-complete most of the command line arguments and flags. Supported shells are `bash`, `zsh` and `fish`.
- **cmd:** allow the baseDir to be passed as an argument ([#6](https://github.com/tbraack/tanka/issues/6)) ([55adf80](https://github.com/tbraack/tanka/commit/55adf80)), ([#12](https://github.com/tbraack/tanka/issues/12)) ([3248bb9](https://github.com/tbraack/tanka/commit/3248bb9)): `tk` **breaks** with the current behaviour and requires the baseDir / environment to be passed explicitely on the command line, instead of assuming it as `pwd`. This is because it allows more `go`-like UX. It is also very handy for scripts not needing to switch the directory.
- **kubernetes:** subset-diff ([#11](https://github.com/tbraack/tanka/issues/11)) ([13f6fdd](https://github.com/tbraack/tanka/commit/13f6fdd)): `tk diff` support for version below Kubernetes `1.13` is here :tada:! The strategy is called _subset diff_ and effectively compares only the fields already present in the config. This allows the (hopefully) most bloat-free experience possible without server side diff.
- **tooling:** import analysis ([#10](https://github.com/tbraack/tanka/issues/10)) ([ce2b0d3](https://github.com/tbraack/tanka/commit/ce2b0d3)): Adds `tk tool imports`, which allows to list all imports of a single file (even transitive ones). Optionally pass a git commit hash, to check whether any of the changed files is imported, to figure out which environments need to be re-applied.

---

## 0.1.0 (2019-07-31)

This release marks the begin of tanka's history :tada:!

As of now, tanka aims to nearly seemlessly connect to the point where [ksonnet](https://github.com/ksonnet/ksonnet) left.
The current feature-set is basic, but usable: The three main workflow commands are available (`show`, `diff`, `apply`), environments are supported, code-sharing is done using [`jb`](https://github.com/jsonnet-bundler/jsonnet-bundler).

Stay tuned!

### Features

- **kubernetes:** Show ([7c4bee8](https://github.com/tbraack/tanka/commit/7c4bee8)): Equivalent to `ks show`, allows previewing the generated yaml.
- **kubernetes:** Diff ([a959f38](https://github.com/tbraack/tanka/commit/a959f38)): Uses the `kubectl diff` to obtain a sanitized difference betweent the current and the desired state. Requires Kubernetes 1.13+
- **kubernetes:** Apply ([8fcb4c1](https://github.com/tbraack/tanka/commit/8fcb4c1)): Applies the changes to the cluster (like `ks apply`)
- **kubernetes:** Apply approval ([4c6414f](https://github.com/tbraack/tanka/commit/4c6414f)): Requires a typed `yes` to apply, gives the user the chance to verify cluster and context.
- **kubernetes:** Smart context ([2b3fd3c](https://github.com/tbraack/tanka/commit/2b3fd3c)): Infers the correct context from the `spec.json`. Prevents applying the correct config to the wrong cluster.
- Init Command ([ff8857c](https://github.com/tbraack/tanka/commit/ff8857c)): Initializes a new repository with the suggested directory structure.
