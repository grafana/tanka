# Changelog

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
way than `kubectl** does natively.  
During diff, all objects of an in-existent namespace are shown as new and when
applying, namespaces are applied first to allow applying in a single step.

### Features
* **tool/imports**: import analysis using upstream jsonnet: Due to recent
  changes to google/jsonnet, we can now use the upstream compiler for static
  import analysis
  ([#84](https://github.com/grafana/tanka/commit/394cb12b28beb0ea05d065594b6cf5c3f92de5e4))
* **Array output**: The output of Jsonnet may now be an array of Manifests.
  Nested arrays are not supported yet.
  ([#112](https://github.com/grafana/tanka/commit/eb647793ff5515bc828e4f91186655c143bb6a04))


### Bug Fixes
* **Command Usage Guidelines**: Tanka now uses the [command description
  syntax](https://en.wikipedia.org/wiki/Command-line_interface#Command_description_syntax)
  ([#94](https://github.com/grafana/tanka/commit/13238e5941bd6e68f410d3938d1a285224c2f91d))
* **cli/env** resolved panic on missing `spec.json`
  ([#108](https://github.com/grafana/tanka/commit/9bd15e6b4226164efe45f50c9ed41c4a5673ea2d))

## 0.5.0 (2019-09-20)
This version adds a set of commands to manipulate environments (`tk env add, rm,
set, list`) ([#73](https://github.com/grafana/tanka/pull/73)). The commands are
mostly `ks env` compatible, allowing `tk env` be used as a drop-in replacement
in scripts.

Furthermore, an error message has been improved, to make sure users can
differentiate between parse issues in `.jsonnet` and `spec.json`
([#71](https://github.com/grafana/tanka/pull/71)).

## 0.4.0 (2019-09-06)
After nearly a month, the next feature packed release of Tanka is ready!
Highlights include the new documentation website https://tanka.dev, regular
expression support for targets, diff histograms and several bug-fixes.

### Features
* **cli**: `tk show` now aborts by default, when invoked in a non-interactive
  session. Use `--dangerous-allow-redirect` to disable this safe-guard
  ([#47](https://github.com/grafana/tanka/issues/47)).
* **kubernetes**: Regexp Targets: It is now possible to use regular expressions
  when specifying the targets using `--target` / `-t`. Use it to easily select
  multiple objects at once: https://tanka.dev/targets/#regular-expressions
  ([#64](https://github.com/grafana/tanka/issues/64)).
* **kubernetes**: Diff histogram: Tanka now allows to summarize the differences
  between the live configuration and the local one, by using the unix
  `diffstat(1)` utility. Gain a sneek peek at a change using `tk diff -s .`!
  ([#67](https://github.com/grafana/tanka/issues/67))

### Bug Fixes
* **kubernetes**: Tanka does not fail anymore, when the configuration file
  `spec.json` is missing from an Environment. While you cannot apply or diff,
  the show operation works totally fine
  ([#56](https://github.com/grafana/tanka/issues/56),
  [#63](https://github.com/grafana/tanka/issues/63)).
* **kubernetes**: Errors from `kubectl` are now correctly passed to the user
  ([#61](https://github.com/grafana/tanka/issues/61)).
* **cli**: `tk diff` does not output useless empty lines (`\n`) anymore
  ([#62](https://github.com/grafana/tanka/issues/62)).

## 0.3.0 (2019-08-13)

Tanka v0.3.0 is here!

This version includes lots of tiny fixes and detail improvements, to make it easier for everyone to configure their Kubernetes clusters.

Enjoy target support, enhancements to the diff UX and an improved CLI experience.

### Features

The most important feature is **target support** ([#30](https://github.com/tbraack/tanka/issues/30)) ([caf205a](https://github.com/tbraack/tanka/commit/caf205a)): Using `--target=kind/name`, you can limit your working set to a subset of the objects, e.g. to do a staged rollout.

There where some other features added:
* **cli:** autoApprove, forceApply ([#35](https://github.com/tbraack/tanka/issues/35)) ([626b097](https://github.com/tbraack/tanka/commit/626b097)): allows to skip the interactive verification. Furthermore, `kubectl` can now be invoked with `--force`.
* **cli:** print deprecated warnings in verbose mode. ([#39](https://github.com/tbraack/tanka/issues/39)) ([6de170d](https://github.com/tbraack/tanka/commit/6de170d)): Warnings about the deprecated configs are only printed in verbose mode
* **kubernetes:** add namespace to apply preamble ([#23](https://github.com/tbraack/tanka/issues/23)) ([9e2d927](https://github.com/tbraack/tanka/commit/9e2d927)): The interactive verification now shows the `metadata.namespace` as well.
* **cli:** diff UX enhancements ([#34](https://github.com/tbraack/tanka/issues/34)) ([7602a19](https://github.com/tbraack/tanka/commit/7602a19)): The user experience of the `tk diff` subcommand has been improved:
  * if the output is too long to fit on a single screen, the systems `PAGER` is invoked
  * if differences are found, the exit status is set to `16`.
  * When `tk apply` is invoked, the diff is shown again, to make sure you apply what you want

### Bug Fixes

* **cli:** invalid command being executed twice ([#42](https://github.com/tbraack/tanka/issues/42)) ([28c6898](https://github.com/tbraack/tanka/commit/28c6898)): When the command failed, it was executed twice, due to an error in the error handling of the CLI.
* **cli**: config miss ([#22](https://github.com/tbraack/tanka/issues/22)) ([32bc8a4](https://github.com/tbraack/tanka/commit/32bc8a4)): It was not possible to use the new configuration format, due to an error in the config parsing.
* **cli:** remove datetime from log ([#24](https://github.com/tbraack/tanka/issues/24)) ([1e37b20](https://github.com/tbraack/tanka/commit/1e37b20))
* **kubernetes:** correct diff type on 1.13 ([#31](https://github.com/tbraack/tanka/issues/31)) ([574f946](https://github.com/tbraack/tanka/commit/574f946)): On kubernetes 1.13.0, `subset` was used, although `native` is already supported.
* **kubernetes:** Nil pointer deference in subset diff. ([#36](https://github.com/tbraack/tanka/issues/36)) ([f53c2b5](https://github.com/tbraack/tanka/commit/f53c2b5))
* **kubernetes:** sort during reconcile ([#33](https://github.com/tbraack/tanka/issues/33)) ([ab9c43a](https://github.com/tbraack/tanka/commit/ab9c43a)): The output of the reconcilation phase is now stable in ordering


## [0.2.0](https://github.com/tbraack/tanka/compare/v0.1.0...v0.2.0) (2019-08-07)

### Features

* **cli:** Completions ([#7](https://github.com/tbraack/tanka/issues/7)) ([aea3bdf](https://github.com/tbraack/tanka/commit/aea3bdf)): Tanka is now able auto-complete most of the command line arguments and flags. Supported shells are `bash`, `zsh` and `fish`.
* **cmd:** allow the baseDir to be passed as an argument ([#6](https://github.com/tbraack/tanka/issues/6)) ([55adf80](https://github.com/tbraack/tanka/commit/55adf80)), ([#12](https://github.com/tbraack/tanka/issues/12)) ([3248bb9](https://github.com/tbraack/tanka/commit/3248bb9)): `tk` **breaks** with the current behaviour and requires the baseDir / environment to be passed explicitely on the command line, instead of assuming it as `pwd`. This is because it allows more `go`-like UX. It is also very handy for scripts not needing to switch the directory.
* **kubernetes:** subset-diff ([#11](https://github.com/tbraack/tanka/issues/11)) ([13f6fdd](https://github.com/tbraack/tanka/commit/13f6fdd)): `tk diff` support for version below Kubernetes `1.13` is here :tada:! The strategy is called *subset diff* and effectively compares only the fields already present in the config. This allows the (hopefully) most bloat-free experience possible without server side diff.
* **tooling:** import analysis ([#10](https://github.com/tbraack/tanka/issues/10)) ([ce2b0d3](https://github.com/tbraack/tanka/commit/ce2b0d3)): Adds `tk tool imports`, which allows to list all imports of a single file (even transitive ones). Optionally pass a git commit hash, to check whether any of the changed files is imported, to figure out which environments need to be re-applied.

## 0.1.0 (2019-07-31)

This release marks the begin of tanka's history :tada:!

As of now, tanka aims to nearly seemlessly connect to the point where [ksonnet](https://github.com/ksonnet/ksonnet) left.
The current feature-set is basic, but usable: The three main workflow commands are available (`show`, `diff`, `apply`), environments are supported, code-sharing is done using [`jb`](https://github.com/jsonnet-bundler/jsonnet-bundler).

Stay tuned!

### Features

* **kubernetes:** Show ([7c4bee8](https://github.com/tbraack/tanka/commit/7c4bee8)): Equivalent to `ks show`, allows previewing the generated yaml.
* **kubernetes:** Diff ([a959f38](https://github.com/tbraack/tanka/commit/a959f38)): Uses the `kubectl diff` to obtain a sanitized difference betweent the current and the desired state. Requires Kubernetes 1.13+
* **kubernetes:** Apply ([8fcb4c1](https://github.com/tbraack/tanka/commit/8fcb4c1)): Applies the changes to the cluster (like `ks apply`)
* **kubernetes:** Apply approval ([4c6414f](https://github.com/tbraack/tanka/commit/4c6414f)): Requires a typed `yes` to apply, gives the user the chance to verify cluster and context.
* **kubernetes:** Smart context ([2b3fd3c](https://github.com/tbraack/tanka/commit/2b3fd3c)): Infers the correct context from the `spec.json`. Prevents applying the correct config to the wrong cluster.
* Init Command ([ff8857c](https://github.com/tbraack/tanka/commit/ff8857c)): Initializes a new repository with the suggested directory structure.
