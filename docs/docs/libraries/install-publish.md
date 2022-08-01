---
name: Installing and publishing
route: /libraries/install-publish/
menu: Libraries
---

# Installing and publishing
The tool for dealing with libraries is
[`jsonnet-bundler`](https://github.com/jsonnet-bundler/jsonnet-bundler). It can
install packages from any git source using `ssh` and GitHub over `https`.

## Install a library
To install a library from GitHub, use one of the following:

```bash
$ jb install github.com/<user>/<repo>
$ jb install github.com/<user>/<repo>/<subdir>
$ jb install github.com/<user>/<repo>/<subdir>@<version>
```

Otherwise, use the ssh syntax:

```bash
$ jb install git+ssh://git@mycode.server:<path-to-repo>.git
$ jb install git+ssh://git@mycode.server:<path-to-repo>.git/<subdir>
$ jb install git+ssh://git@mycode.server:<path-to-repo>.git/<subdir>@<version>
```

> **Note**: `version` may be any git ref, such as commits, tags or branches

## Publish to Git(Hub)
Publishing is as easy as committing and pushing to a git remote.
[GitHub](https://github.com) is recommended, as it is most common and supports
faster installing using http archives.
