---
name: "Diff strategies"
route: "/diff-strategy"
menu: "References"
---

# Diff Strategies

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
[known issue](known-issues.md#unexpected-diff-if-the-same-port-number-is-used-for-udp-and-tcp)
with `kubectl diff`, which affects ports configured to use both TCP and UDP.

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
