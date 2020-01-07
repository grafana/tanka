---
name: Import paths
route: /libraries/import-paths
menu: Libraries
---

# Import paths

When using `import` or `importstr`, Tanka considers the following directories to
find a suitable file for that specific import:

1. `<baseDir>`: The directory of your environment (`/environments/default`,
   etc). Put things that only belong to a single environment here.
2. `/lib`: Libraries created for this very project, not meant to be shared
   otherwise. Put everything you need across multiple environments here.
3. `/vendor`: Shared libraries installed using `jsonnet-bundler`. Do not modify
   this folder by hand, your changes will be overwritten by `jb` anways.

> **Note**: The directories are visited in the above order. For example, when a
> file is present in both, `/lib` and `/vendor`, the one from `/lib` will be
> taken, as it occurs higher in the list.

### Shadowing
It is possible to shadow certain files (overlay them with another version), by
putting a file with the exact same name and into a higher ranked import path.
This can be handy if you need to do temporary changes to a vendored library by
overlaying the to-be-changed files using new ones in `lib/`.

For example, to shadow `/vendor/my/lib/file.libsonnet`, copy it to
`/lib/my/lib/file.libsonnet` and do your changes. Tanka will take the file in
`lib/` instead of `vendor/` from now on.
