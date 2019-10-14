# Known Errors

Below is a list of common error messages and how to address them.

### `Evaluating jsonnet: RUNTIME ERROR: Undefined external variable: __ksonnet/components`
When migrating from `ksonnet`, this error might occur, because Tanka does not
provide the global `__ksonnet` variable, nor does it strictly have the concept
of components.  
You will need to use the plain Jsonnet `import` feature instead. Note that this
requires your code to be inside of one of the [import
paths](directory-structure.md/#import-paths).

### `Evaluating jsonnet: RUNTIME ERROR: couldn't open import "k.libsonnet": no match locally or in the Jsonnet library paths`
This error can occur when the `ksonnet` kubernetes libraries are missing in the import paths. While `ksonnet` used to magically include them, Tanka follows a more explicit approach and requires you to install them using `jb`:

```bash
$ jb install github.com/ksonnet/ksonnet-lib/ksonnet.beta.3
```

This installs version `beta.3` of the libraries, matching Kubernetes version
`1.8.0`. If you need another version, take a look at
https://github.com/ksonnet/ksonnet-lib. When a pre-compiled version is
available, install it using `jb`, otherwise compile it yourself and place it
under `lib/`.
