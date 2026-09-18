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

From now on, you can use `tk prune` to remove old resources from your cluster. You can also limit pruning to a single namespace, with `tk prune --namespace <my-namespace>`.

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

## Excluding specific resources from prune

Some resources aren't created by Tanka directly, but by a third-party
operator acting on a Tanka-managed object. For example,
[external-secrets](https://external-secrets.io/) generates a `Secret` from an
`ExternalSecret`, and by default copies the `ExternalSecret`'s own metadata
onto it — including the `tanka.dev/environment` label. That makes the
generated `Secret` look like an orphaned Tanka resource, even though Tanka
never applied it and has no desired state for it.

**Owned resources are excluded automatically.** If the generated resource
carries a Kubernetes [ownerReference](https://kubernetes.io/docs/concepts/overview/working-with-objects/owners-dependents/)
with `controller: true` back to its parent — which external-secrets sets by
default (`creationPolicy: Owner`) — Tanka never considers it for pruning, no
configuration needed. Its lifecycle belongs to the owning controller, and
Kubernetes' own garbage collector already deletes it when the parent is
deleted.

This doesn't cover every case: `creationPolicy: Orphan` (or a similar setting
on another operator), for instance, intentionally omits the ownerReference so
the generated resource survives deletion of its parent. For those, set the
`tanka.dev/prune-ignore: "true"` annotation on the resource instead. Since you
usually don't control the generated resource's manifest directly, set the
annotation on the template the operator uses to build it, so it gets
propagated. For external-secrets, that's
`spec.target.template.metadata.annotations`:

```jsonnet
{
  apiVersion: 'external-secrets.io/v1',
  kind: 'ExternalSecret',
  spec: {
    target: {
      template: {
        metadata: {
          annotations: {
            'tanka.dev/prune-ignore': 'true',
          },
        },
      },
    },
  },
}
```

To see which live resources in an environment are currently protected this
way, without running a prune, use `--list-ignored`:

```bash
tk prune --list-ignored .
```
