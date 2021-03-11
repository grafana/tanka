---
name: Overriding
route: /libraries/overriding
menu: Libraries
---

# Overriding vendor

The `vendor` directory is immutable in its nature. You can't and should never
modify any files inside of it, `jb` will revert those changes on the next run anyway.

Nevertheless, it can sometimes become required to add changes there, e.g. if an
upstream library contains a bug that needs to be fixed immediately, without
waiting for the upstream maintainer to review it.

## Shadowing

Because [import paths](/libraries/import-paths) are ranked in Tanka, you can use
a technique called shadowing: By putting a file with the exact same name in a
higher ranked path, Tanka will prefer that file instead of the original in
`vendor`, which has the lowest possible rank of 1.

For example, if `/vendor/foo/bar.libsonnet` contained an error, you could create
`/lib/foo/bar.libsonnet` and fix it there.

> **Tip:** Instead of copying the file to the new location and making the edits,
> use an absolute import and [patching](/tutorial/environments#patching):
>
> ```jsonnet
> // in /lib/foo/bar.libsonnet:
> (import "../../vendor/foo/bar.libsonnet") + {
>   foo+: {
>     bar: "fixed"
>   }
> }
> ```

> **Important:** If the file you override is not the one you directly import,
> but instead imported by another file first, the override will only occur if
> the placement of the file is alongside your `main.libsonnet`.  This is due to
> the logic behind the Jsonnet importer.  Example:  We import
> `abc/main.libsonnet` located in `vendor/abc`.  Because Jsonnet first looks if
> files are locally present before considering the [import
> paths](/libraries/import-paths), you need to make sure your override is
> actually picked up. In our example, you'd need to copy the `main.libsonnet`
> into `lib/abc` as well.

## Per environment

Another common case is overriding the entire `vendor` bundle per environment.

This is handy, when you for example want to test a change of an upstream
library which is used in many environments (including `prod`) in a single one,
without affecting all the others.

For this, Tanka lets you have a separate `vendor`, `jsonnetfile.json` and
`jsonnetfile.lock.json` per environment. To do so:

#### Create `tkrc.yaml`

Tanka normally uses the `jsonnetfile.json` from your project to find its root.
As we are going to create another one of that down the tree in the next step, we
need another marker for `<rootDir>`.

For that, create an empty file called `tkrc.yaml` in your project's root,
alongside the original `jsonnetfile.json`.

> **Info**: While the name suggests that `tkrc.yaml` could be used for setting
> parameters, this is not the case yet.  
> It might however be repurposed later, in case we need such functionality

#### Add a `vendor` to your environment

In your environments folder (e.g. `/environments/default`):

```bash
# init jsonnet bundler (creates jsonnetfile.json)
$ jb init

# install the updated dependency
$ jb init github.com/foo/bar@v2
```

> **Tip**: You don't need to install everything into the new `vendor/`, as
> packages not present there can still be imported from the global `/vendor`.
