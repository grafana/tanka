---
name: "Exporting as YAML"
route: "/exporting"
menu: Advanced features
---

# Exporting as YAML

Tanka provides you with a day-to-day workflow for working with Kubernetes clusters:

- `tk show` for quickly checking the YAML representation looks good
- `tk diff` to ensure your changes will behave like they should
- `tk apply` makes it happen

However sometimes it can be required to integrate with other tooling that does
only support `.yaml` files.

For that case, `tk export` can be used:

```bash
#           <environment>         <outputDir>
$ tk export environments/promtail promtail
```

This will create a separate `.yaml` file for each Kubernetes resource included in your Jsonnet.

## Filenames

Tanka by default uses the following pattern:

```bash
${apiVersion}.${kind}-${metadata.name}.yaml

# examples:
apps-v1.Deployment-distributor.yaml
v1.ConfigMap-loki.yaml
v1.Service-ingester.yaml
```

If that does not fit your need, you can provide your own pattern using the `--format` flag:

```bash
tk export environments/promtail promtail --format='{{.metadata.labels.app}}-{{.metadata.name}}-{{.kind}}'
```

> The syntax is Go `text/template`. See https://golang.org/pkg/text/template/
> for reference.

This would include the label named `app`, the `name` and `kind` of the resource:

```
loki-distributor-Deployment
loki-loki-ConfigMap
loki-ingester-Service
```

You can optionally use the template function `lower` for lower-casing fields, e.g. in the above example

```bash
... --format='{{.metadata.labels.app}}-{{.metadata.name}}-{{.kind | lower}}'
```

would yield

```
loki-distributor-deployment
```

etc.

You can also use a different file extension by providing `--extension='yml'`, for example.
