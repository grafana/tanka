---
name: "Formatting"
route: "/formatting/"
menu: "References"
---

# File Formatting

Tanka supports formatting for all `jsonnet` and `libsonnet` files using the `tk fmt` command.

By default, the command excludes all `vendor` directories.

```bash
# Run for current and child directories. Run this in the root of the project to format all your files.
tk fmt .

# Format a single file (myFile.jsonnet)
tk fmt myFile.jsonnet

# Use the `-t` tag to test (Dry run).
tk fmt -t myFile.jsonnet

# Format using verbose mode.
tk fmt -v .
```