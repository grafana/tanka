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

Tanka can also export multiple inline environments, as showcased in [Use case: consistent inline
environments](/inline-environments#use-case-consistent-inline-environments). This follows the same
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

## Using a memory ballast

_Read [this blog post](https://blog.twitch.tv/en/2019/04/10/go-memory-ballast-how-i-learnt-to-stop-worrying-and-love-the-heap/) for more information about memory ballasts._

For large environments that load lots of data into memory on evaluation, a memory ballast can dramatically improve performance. This feature is exposed through the `--mem-ballast-size-bytes` flag on the export command.

Anecdotally (Grafana Labs), environments that took around a minute to load were able to load in around 45 secs with a ballast of 5GB (`--mem-ballast-size-bytes=5368709120`). Decreasing the ballast size resulted in negative impact on performance, and increasing it more did not result in any noticeable impact.

## Caching

Tanka can also cache the results of the export. This is useful if you often export the same files and want to avoid recomputing them. The cache key is calculated from the main file and all of its transitive imports, so any change to any file possibly used in an environment will invalidate the cache.

This is configured by two flags:

- `--cache-path`: The local filesystem path where the cache will be stored. The cache is a flat directory of json files (one per environment).
- `--cache-envs`: If exporting multiple environments, this flag can be used to specify, with regexes, which environments to cache. If not specified, all environments are cached.

Notes:

- Using the cache might be slower than evaluating jsonnet directy. It is only recommended for environments that are very CPU intensive to evaluate.
- To use object storage, you can point the `--cache-path` to a FUSE mount, such as [`s3fs`](https://github.com/s3fs-fuse/s3fs-fuse)
