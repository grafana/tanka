# Changelog

## 0.12 (2020-10-05)

Like good wine, some things need time. After 3 months of intense development we
have another Tanka release ready:

#### :wheel_of_dharma: Helm support

This one is huge! Tanka can now **load Helm Charts**:

- [`helm-util`](https://github.com/grafana/jsonnet-libs/tree/master/helm-util)
  provides `helm.template()` to load them from inside Jsonnet
- Declarative vendoring using `tk tool charts`
- Jsonnet-native overwriting of chart contents

Just by upgrading to 0.12, you have access to every single Helm chart on the
planet, right inside of Tanka! Read more on https://tanka.dev/helm

#### :house: Top Level Arguments

Tanka now supports the `--tla-str` and `--tla-code` flags from the `jsonnet` cli
to late-bind data into the evaluation in a well-defined way. See
https://tanka.dev/jsonnet/injecting-values for more details.

#### :sparkles: Inline Eval

Ever wanted to pull another value out of Jsonnet that does not comply to the
Kubernetes object rules Tanka imposes onto everything? Wait no longer and use
`tk eval -e`:

```console
$ tk eval environments/prometheus -e prometheus_rules
```

Above returns `$.prometheus_rules` as JSON. Every Jsonnet selector is supported:

```console
$ tk eval environments/prometheus -e 'foo.bar[0]'
```

### Features

- **k8s, jsonnet** :sparkles:: Support for [Helm](https://helm.sh). In combination with
  [`helm-util`](https://github.com/grafana/jsonnet-libs/tree/master/helm-util),
  Tanka can now load resources from Helm Charts.
  **([#336](https://github.com/grafana/tanka/pull/336))**
- **k8s**: Default metadata from `spec.json`
  **([#366](https://github.com/grafana/tanka/pull/366))**

* **helm**: Charttool: Adds `tk tool charts` for easy management of vendored
  Helm charts **([#367](https://github.com/grafana/tanka/pull/367))**,
  **([#369](https://github.com/grafana/tanka/pull/369))**
* **helm**: Require Helm Charts to be available locally
  **([#370](https://github.com/grafana/tanka/pull/370))**
* **helm**: Configurable name format
  **([#381](https://github.com/grafana/tanka/pull/381))**

- **cli**: Filtering (`-t`) now supports negative expressions (`-t !deployment/.*`) to exclude resources
  **([#339](https://github.com/grafana/tanka/pull/339))**
- **cli** :sparkles:: Inline eval (Use `tk eval -e` to extract nested fields)
  **([#378](https://github.com/grafana/tanka/pull/378))**
- **cli**: Custom paging (`PAGER` env var)
  **([#373](https://github.com/grafana/tanka/pull/373))**
- **cli**: Predict plain directories if outside a project
  **([#357](https://github.com/grafana/tanka/pull/357))**

* **jsonnet** :sparkles:: Top Level Arguments can now be specified using `--tla-str` and
  `--tla-code` **([#340](https://github.com/grafana/tanka/pull/340))**

### Bug Fixes

- **yaml**: Pin yaml library to v2.2.8 to avoid whitespace changes
  **([#386](https://github.com/grafana/tanka/pull/386))**
- **cli**: Actually respect `TANKA_JB_PATH`
  **([#350](https://github.com/grafana/tanka/pull/350))**
- **k8s**: Update `kubectl v1.18.0` warning
  **([#371](https://github.com/grafana/tanka/pull/371))**

* **jsonnet**: Load `main.jsonnet` using full path. This makes `std.thisFile`
  usable **([#370](https://github.com/grafana/tanka/pull/370))**
* **jsonnet**: Import path resolution now works on Windows
  **([#331](https://github.com/grafana/tanka/pull/331))**
* **jsonnet**: Arrays are now supported at the top level
  **([#321](https://github.com/grafana/tanka/pull/321))**

### BREAKING

- **api**: Struct based Go API: Modifies our Go API
  (`github.com/grafana/tanka/pkg/tanka`) to be based on structs instead of
  variadic arguments. This has no impact on daily usage of Tanka.
  **([#376](https://github.com/grafana/tanka/pull/376))**
- **jsonnet**: ExtVar flags are now `--ext-str` and `--ext-code` (were `--extVar` and `--extCode`)
  **([#340](https://github.com/grafana/tanka/pull/340))**

## 0.11.1 (2020-07-17)

This is a minor release with one bugfix and one minor feature.

### Features

- **process**: With 0.11.0, tanka started automatically adding namespaces to _all_ manifests it processed. We updated this to _not_
  add a namespace to cluster-wide object types in order to make handling of these resources more consistent in different workflows. **([#320](https://github.com/grafana/tanka/pull/320))**

### Bug Fixes

- **export**: Fix inverted logic while checking if a file already exists. This broke `tk export` entirely.
  **([#317](https://github.com/grafana/tanka/pull/317))**

## 0.11.0 (2020-07-07)

2 months later and here we are with another release! Packed with many
detail-improvements, this is what we want to highlight:

#### :sparkles: Enhanced Kubernetes resource handling

From now on, Tanka handles the resources it extracts from your Jsonnet output in
an enhanced way:

1. **Lists**: Contents of lists, such as `RoleBindingList` are automatically
   flattened into the resource stream Tanka works with. This makes sure they are
   properly labeled for garbage collection, etc.
2. **Default namespaces**: While you could always define the default namespace
   (the one for resources without an explicit one) in `spec.json`, this
   information is now also persisted into the YAML returned by `tk show` and `tk export`.
   See https://tanka.dev/namespaces for more information.

#### :hammer: More powerful exporting

`tk export` can now do even more than just writing YAML files to disk:

1. `--extension` can be used to control the file-extension (defaults to `.yaml`)
2. When you put a `/` in your `--format` for the filename, Tanka creates a
   directory or you. This allows e.g. sorting by namespace:
   `--format='{{.metadata.namespace}}/{{.kind}}-{{.metadata.name}}'`
3. Using `--merge`, you can export multiple environments into the same directory
   tree, so you get the full YAML picture of your entire cluster!

#### :fax: Easier shell scripting

The `tk env list` command now has a `--names` option making it easy to operate on multiple environments:

```bash
# diff all environments:
for e in $(tk env list --names); do
  tk diff $e;
done
```

Also, to use a more granular subset of your environments, you can now use
`--selector` / `-l` to match against `metadata.labels` of defined in your
`spec.json`:

```bash
$ tk env list -l status=dev
```

### Features

- **cli**: `tk env list` now supports label selectors, similar to `kubectl get -l` **([#295](https://github.com/grafana/tanka/pull/295))**
- **cli**: If `spec.apiServer` of `spec.json` lacks a protocol, it now defaults
  to `https` **([#289](https://github.com/grafana/tanka/pull/289))**
- **cli**: `tk delete` command to teardown environments
  **([#313](https://github.com/grafana/tanka/pull/313))**

* **cli**: Support different file-extensions than `.yaml` for `tk export`
  **([#294](https://github.com/grafana/tanka/pull/394))** (**@marthjod**)
* **cli**: Support creating sub-directories in `tk export`
  **([#300](https://github.com/grafana/tanka/pull/300))** (**@simonfrey**)
* **cli**: Allow writing into existing folders during `tk export`
  **([#314](https://github.com/grafana/tanka/pull/314))**

- **tooling**: `tk tool imports` now follows symbolic links
  **([#302](https://github.com/grafana/tanka/pull/302))**,
  **([#303](https://github.com/grafana/tanka/pull/303))**

* **process**: `List` types are now unwrapped by Tanka itself
  **([#306](https://github.com/grafana/tanka/pull/306))**
* **process**: Automatically set `metadata.namespace` to the value of
  `spec.namespace` if not set from Jsonnet
  **([#312](https://github.com/grafana/tanka/pull/312))**

### Bug Fixes

- **jsonnet**: Using `import "tk"` twice no longer panics
  **([#290](https://github.com/grafana/tanka/pull/290))**
- **tooling**: `tk tool imports` no longer gets stuck when imports are recursive
  **([#298](https://github.com/grafana/tanka/pull/298))**
- **process**: Fully deterministic recursion, so that error messages are
  consistent **([#307](https://github.com/grafana/tanka/pull/307))**

## 0.10.0 (2020-05-07)

New month, new release! And this one ships with a long awaited feature:

#### :sparkles: Garbage collection

Tanka can finally clean up behind itself. By optionally attaching a
`tanka.dev/environment` label to each resource it creates, we can find these
afterwards and purge those removed from the Jsonnet code. No more dangling
resources!

> :warning: Keep in mind this is still experimental!

To get started, enable labeling in your environment's `spec.json`:

```diff
  "spec": {
+   "injectLabels": true,
  }
```

Don't forget to `tk apply` afterwards! From now on, Tanka can clean up using `tk prune`.

Docs: https://tanka.dev/garbage-collection

#### :boat: Logo

Tanka now has it's very own logo, and here it is:

<img src="docs/img/logo.svg" width="400px" />

#### :package: Package managers

Tanka is now present in some package managers, notably `brew` for macOS and the
AUR of ArchLinux! See the updated [install
instructions](https://tanka.dev/install#using-a-package-manager-recommended) to
make sure to use these if possible.

### Features:

- **cli**: `TANKA_JB_PATH` environment variable introduced to set the `jb`
  binary if required **([#272](https://github.com/grafana/tanka/pull/272))**.
  Thanks [@qckzr](https://github.com/qckzr)

* **kubernetes**: Garbage collection
  **([#251](https://github.com/grafana/tanka/pull/251))**

### Bug Fixes

- **kubernetes**: Resource sorting is now deterministic
  **([#259](https://github.com/grafana/tanka/pull/259))**

## 0.9.2 (2020-04-19)

Mini-release to fix an issue with our Makefile (required for packaging). No
changes in functionality.

### Bug Fixes

- **build**: Enable `static` Makefile target on all operating systems
  ([#262](https://github.com/grafana/tanka/pull/262))

## 0.9.1 (2020-04-08)

Small patch release to fix a `panic` issue with `tk apply`.

### Bug Fixes

- **kubernetes**: don't panic on failed diff
  **([#256](https://github.com/grafana/tanka/pull/256))**

## 0.9.0 (2020-04-07)

**This release includes a critical fix, update ASAP**.

Another Tanka release is here, just in time for Easter. Enjoy the built-in
[formatter](#sparkles-highlight-jsonnet-formatter-tk-fmt), much [more
intelligent apply](#rocket-highlight-sorting-during-apply) and several important
bug fixes.

#### :rotating_light: Alert: `kubectl diff` changes resources :rotating_light:

The recently released `kubectl` version `v1.18.0` includes a **critical issue**
that causes `kubectl diff` (and so `tk diff` as well) to **apply** the changes.

This can be very **harmful**, so Tanka decided to require you to **downgrade**
to `v1.17.x`, until the fix in `kubectl` version `v1.18.1` is released.

- Upstream issue: https://github.com/kubernetes/kubernetes/issues/89762)
- Unreleased fix: https://github.com/kubernetes/kubernetes/pull/89795

#### :sparkles: Highlight: Jsonnet formatter (`tk fmt`)

Since `jsonnetfmt` was [rewritten in Go
recently](https://github.com/google/go-jsonnet/pull/388), Tanka now ships it as
`tk fmt`. Just run `tk fmt .` to keep all Jsonnet files recursively formatted.

#### :rocket: Highlight: Sorting during apply

When using `tk apply`, Tanka now automatically **sorts** your objects
based on **dependencies** between them, so that for example
`CustomResourceDefinitions` created before being used, all in the same run. No
more partly failed applies!

### Features

- **kubernetes** :sparkles:: Objects are now sorted by dependency before `apply`
  **([#244](https://github.com/grafana/tanka/pull/244))**
- **cli**: Env var `TANKA_KUBECTL_PATH` can now be used to set a custom
  `kubectl` binary
  **([#221](https://github.com/grafana/tanka/pull/221))**
- **jsonnet** :sparkles: : Bundle `jsonnetfmt` as `tk fmt`
  **([#241](https://github.com/grafana/tanka/pull/241))**

* **docker**: The Docker image now includes GNU `less`, instead of the BusyBox
  one **([#232](https://github.com/grafana/tanka/pull/232))**
* **docker**: Added `kubectl`, `jsonnet-bundler`, `coreutils`, `git` and
  `diffutils` to the Docker image, so Tanka can be fully used in there.
  **([#243](https://github.com/grafana/tanka/pull/243))**

### Bug Fixes

- **cli**: The diff shown on `tk apply` is now colored again
  **([#216](https://github.com/grafana/tanka/pull/216))**

* **client**: The namespace patch file saved to a temporary location is now
  removed after run **([#225](https://github.com/grafana/tanka/pull/225))**
* **client**: Scanning for the correct context won't panic anymore, but print a
  proper error **([#228](https://github.com/grafana/tanka/pull/228))**
* **client**: Use `os.PathListSeparator` during context patching, so that Tanka
  also works on non-UNIX platforms (e.g. Windows)
  **([#242](https://github.com/grafana/tanka/pull/242))**

- **kubernetes** :rotating_light:: Refuse to diff on `kubectl` version `v1.18.0`, because of
  above mentioned unfixed issue
  **([#254](https://github.com/grafana/tanka/pull/254))**
- **kubernetes**: Apply no longer aborts when diff fails
  **([#231](https://github.com/grafana/tanka/pull/231))**
- **kubernetes** :sparkles:: Namespaces that will be created in the same run are now
  properly handled during `diff`
  **([#237](https://github.com/grafana/tanka/pull/237))**

### Other

- **cli**: Migrates from `spf13/cobra` to much smaller `go-clix/cli`. This cuts
  our dependencies to a minimum.
  **([#235](https://github.com/grafana/tanka/pull/235))**

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
