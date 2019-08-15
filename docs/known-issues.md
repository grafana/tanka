# Known Errors

Below is a list of common error messages and how to address them.

### `Evaluating jsonnet: RUNTIME ERROR: Undefined external variable: __ksonnet/components`
When migrating from `ksonnet`, this error might occur, because Tanka does not
provide the global `__ksonnet` variable, nor does it strictly have the concept
of components.  
You will need to use the plain Jsonnet `import` feature instead. Note that this
requires your code to be inside of one of the [import
paths](directory-structure.md/#import-paths).
