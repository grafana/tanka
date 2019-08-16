# Targets

When a project becomes bigger over time and includes a lot of Kubernetes
objects, it may become required to operate on only a subset of them (e.g. apply
only a part of an application).

Tanka helps you with this, by allowing you to limit the used objects on the command
line using the `--target` flag. Say you are deploying an `nginx` instance with a special
`nginx.conf` and want to apply the `ConfigMap` first:

```bash
# show the ConfigMap
$ tk show -t configmap/nginx .

# all good? apply!
$ tk apply -t configmap/nginx .

# and apply everything else:
$ tk apply .
```

The syntax of the `--target` / `-t` flag is `--target=<kind>/<name>`. If
multiple objects match this pattern, all of them are used.

The `--target` / `-t` flag can be specified multiple times, to work with
multiple objects.
