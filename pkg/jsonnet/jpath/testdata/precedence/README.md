## `precedence` test

This directory contains some Jsonnet testdata for testing correct import precedence.

The desired precedence looks like the following (most specific wins):

1. consider baseDir (dir containing `main.jsonnet`) first: The current
   environment is most important
2. then `/lib`: project libraries are more specific than any vendors
3. then `<baseDir>/vendor`: to allow overriding project vendor on a environment level
4. finally `/vendor`: external packages are least specific

## internals

How the test works:

We basically put the same jsonnet file in multiple locations:

```jsonnet
{
    value: "<location>"
}
```

For example to check for `/lib` to precede both vendor folders, the following is used:

- `/vendor/project_lib.jsonnet`: `value: "/vendor"`
- `/<baseDir>/vendor/project_lib.jsonnet`: `value: "/baseDir-vendor"`
- `/lib/project_lib.jsonnet`: `value: "/lib"`

Then in `<baseDir>/main.jsonnet` we put `import "project_lib.jsonnet"` and
expect `value` to be `/lib`.
