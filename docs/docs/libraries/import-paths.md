---
name: Import paths
route: /libraries/import-paths/
menu: Libraries
---

# Import paths

When using `import` or `importstr`, Tanka considers the following directories to
find a suitable file for that specific import:

| Rank | Path               | Purpose                                                                                                                      |
| ---- | ------------------ | ---------------------------------------------------------------------------------------------------------------------------- |
| 4    | `<baseDir>`        | The directory of your environment, e.g. `/environments/default`.<br /> Put things that belong to this very environment here. |
| 3    | `/lib`             | Project-global libraries, that are used in multiple environments, but are specific to this project.                          |
| 2    | `<baseDir>/vendor` | Per-environment vendor, can be used for [`vendor` overriding](/libraries/overriding#per-environment)                         |
| 1    | `/vendor`          | Global vendor, holds external libraries installed using `jb`.                                                                |

> **Note**:
>
> - If a file occurs in multiple paths, the one with the highest rank will be chosen.
> - `/` in above table means `<rootDir>`, which is your project root.
