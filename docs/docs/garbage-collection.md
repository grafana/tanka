---
name: "Garbage collection"
route: "/garbage-collection"
menu: Advanced features
---

# Garbage collection

Tanka can automatically delete resources from your cluster once you remove them
from Jsonnet.

> **Note:** This feature is **experimental**. Please report problems at https://github.com/grafana/tanka/issues.

To accomplish this, it appends the `tanka.dev/environment: <hash>` label to each created
resource. This is used to identify those which are missing from the local state in the
future.

> **Note:** The label value changed from the <name> to a <hash> in v0.15.0.

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
