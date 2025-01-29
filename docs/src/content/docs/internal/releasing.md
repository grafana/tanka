---
title: 'Releasing a new version'
---

For releasing Tanka we're using [release-please][].
This workflow manages a release pull-request based on the content of the `main` branch that would update the changelog et al..
Once you want to do a release, merge that prepared pull-request.
release-please will then do all the tagging and GitHub Release creation.

[release-please]: https://github.com/googleapis/release-please-action
