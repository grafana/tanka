---
title: Garbage collection
sidebar:
  order: 1
---

Tanka can automatically delete resources from your cluster once you remove them
from Jsonnet.

:::caution
This feature is **experimental**. Please report problems at https://github.com/grafana/tanka/issues.
:::

To accomplish this, it appends the `tanka.dev/environment: <hash>` label to each created
resource. This is used to identify those which are missing from the local state in the
future.

:::note
The label value changed from the `<name>` to a `<hash>` in v0.15.0.
:::

Because the label causes a `diff` for every single object in your cluster and
not everybody wants this, it needs to be explicitly enabled. To do so, add the
following field to your `spec.json`:

```diff
{
  "spec": {
+    "injectLabels": true,
  }
}
```

Once added, run a `tk apply`, make sure the label is actually added and confirm
by typing `yes`.

From now on, you can use `tk prune` to remove old resources from your cluster.

## Filtering pruned resources

By default `tk prune` considers every resource kind labeled with the
environment. In large environments this can be expensive because it must list
every resource type from the Kubernetes API.

You can restrict pruning to a specific subset of resources using the `--target`
(`-t`) flag, which accepts the same `kind/name` regex syntax as the other
workflow commands:

```bash
# prune only orphaned StatefulSets
tk prune -t 'statefulset/.*' .

# prune only StatefulSets whose names start with "live-store"
tk prune -t 'statefulset/live-store.*' .

# prune everything except Deployments
tk prune -t '!deployment/.*' .
```

When a literal kind name is given (no regex metacharacters in the kind
position), Tanka restricts the Kubernetes API query to only that resource type,
avoiding the cost of listing every other kind.

See [Output filtering](/output-filtering) for the full filter syntax.
