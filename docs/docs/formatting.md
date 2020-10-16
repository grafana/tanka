---
name: "Formatting"
route: "/formatting"
menu: "References"
---

# File Formatting

Tanka supports formatting for all `jsonnet` and `libsonnet` files using the `tk fmt` command.

By default, the command excludes all `vendor` directories and any files in the immediate directory.

```bash
# Run for current and child directories. Run this in the root of the project to format all your files.
$ tk fmt .
```
