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
#           <outputDir> <environment>
$ tk export promtail/   environments/promtail
```

> **Note:** The arguments flipped in v0.14.0, the `<outputDir>` comes first now.

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
tk export promtail environments/promtail --format='{{.metadata.labels.app}}-{{.metadata.name}}-{{.kind}}'
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


## Multiple environments

Tanka can also export of multiple inline environments, as showcased in [Use case: consistent inline
environments](https://tanka.dev/inline-environments#use-case-consistent-inline-environments). This follows the same
principles as describe before with the addition that you can also refer to Environment specific data through the `env`
keyword.

For example an export might refer to data from the Environment spec:

```bash
# Format based on environment {{env.<...>}}
$ tk export exportDir environments/dev/ \
  --format '{{env.metadata.labels.cluster}}/{{env.spec.namespace}}//{{.kind}}-{{.metadata.name}}'
```

Even more advanced use cases allow you to export multiple environments in a single execution:

```bash
# Export multiple environments
$ tk export exportDir environments/dev/ environments/qa/
# Recursive export
$ tk export exportDir environments/ --recursive
# Recursive export with labelSelector
$ tk export exportDir environments/ -r -l team=infra
```
