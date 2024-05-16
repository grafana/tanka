---
title: Diff strategies
sidebar:
  order: 5
---

Tanka supports two different ways of computing differences between the local
configuration and the live cluster state: Either **native** `kubectl diff -f -`
is used, which gives the best possible results, but is only possible for
clusters with
[server-side diff](https://kubernetes.io/blog/2019/01/14/apiserver-dry-run-and-kubectl-diff/)
support (Kubernetes 1.13+).

When this is not available, Tanka falls back to `subset` mode.

You can specify the diff-strategy to use on the command line as well:

```bash
# native
tk diff --diff-strategy=native .

# validate: Like native but with a server-side validation
tk diff --diff-strategy=validate .

# server-side
tk diff --diff-strategy=server .

# subset
tk diff --diff-strategy=subset .
```

## Native

The native diff mode is recommended, because it uses `kubectl diff` underneath,
which sends the objects to the Kubernetes API server and computes the
differences over there.

This has the huge benefit that all possible changes by webhooks and other
internal components of Kubernetes can be encountered as well.

However, this is a fairly new feature and only available on Kubernetes 1.13 or
later. Only the API server (master nodes) needs to have that
version, worker nodes do not matter.

There is a
[known issue](./known-issues#unexpected-diff-if-the-same-port-number-is-used-for-udp-and-tcp)
with `kubectl diff`, which affects ports configured to use both TCP and UDP.

### Server-side diffs

There are two additional modes which extend `native`: `validate` and `server`.
While all `kubectl diff` commands are sent to the API server, these two
methods take advantage of an additional server-side diff mode (which uses the
`kubectl diff --server-side` flag, complementing the
[server-side apply](./server-side-apply) mode).

Since a plain `server` diff often produces cruft, and wouldn't be representative
of a client-side apply, the `validate` method allows the server-side diff to
check that all models are valid server-side, but still displays the `native`
diff output to the user.

## Subset

If native diffing is not supported by your cluster, Tanka provides subset diff
as a fallback method.

**Subset diff only compares fields present in the local configuration and
ignores all other fields**. When you remove a field locally, you will see no
differences.

This is required, because Kubernetes adds dynamic fields to the state during
runtime, which we cannot know of on the client side. To produce a somewhat
usable output, we can effectively only compare what we already know about.

If this is a problem for you, consider switching to [native](#native) mode.

## External diff utilities

You can use external diff utilities by setting the environment variable
`KUBECTL_EXTERNAL_DIFF`. If you want to use a GUI or interactive diff utility
you must also set `KUBECTL_INTERACTIVE_DIFF=1` to prevent Tanka from capturing
stdout.
