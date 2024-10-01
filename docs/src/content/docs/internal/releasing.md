---
title: "Releasing a new version"
---

Releasing a new version of Tanka requires a couple of manual and automated steps.
This guide will give you a runbook on how to do them in the right order.

## 1. Create a release tag

1. Pull the latest changes from the `main` branch to your local clone.
1. Create a new tag with the prefix `v` (e.g. `v0.28.0`).
1. Push that tag back to GitHub.

This starts multiple GitHub workflows that will produce binaries, a Docker image, and also a new GitHub release that is *marked as draft*.

## 2. Add changelog to the release notes

Once all these actions have finished, go to <https://github.com/grafana/tanka/releases> and you should see the new draft release.
Click the pencil icon ("edit") at the top-right and go to the last line of the text body.
Now hit the "Generate release notes" button to add a changelog to the end of the release notes.

## 3. Publish the release notes

Once you've check that the release looks fine (e.g. no broken links, no missing version numbers in the download paths) click the "Publish release" button.
