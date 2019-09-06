# Changelog

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
