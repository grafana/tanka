---
name: "Server-Side Apply"
route: "/server-side-apply"
menu: Advanced features
---

# Server-Side Apply

Tanka supports
[server-side apply](https://kubernetes.io/docs/reference/using-api/server-side-apply/),
which requires at least Kubernetes 1.16+, and was promoted to stable status in 1.22.

To enable server-side diff in tanka, add the following field to `spec.json`:

```diff
{
  "spec": {
+    "applyStrategy": "server",
  }
}
```

While server-side apply doesn't have any effect on the resources being applied
and is intended to be a general in-place upgrade to client-side apply, there are
differences in how fields are managed that can make converting existing cluster
resources a non-trival change.

Identifying and fixing these changes are beyond the scope of this guide, but
many can be found before an apply by using the `validate` or `server`
[diff strategy](/diff-strategy).

## Field conflicts

As part of the changes, you may encounter error messages which
recommend the use of the `--force-conflicts` flag. Using `tk apply --force`
in server-side mode will enable that flag for kubectl instead of
`kubectl --force`, which no longer has any effect in server-side mode.
