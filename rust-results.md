Friday project: preliminary results of integrating the Rust jsonnet implementation (https://github.com/CertainLach/jrsonnet) into Tanka using C bindings (Doing CGO on MacOS is fun!)
It's FAST. Not every works yet, but `grafana-o11y` does and it's a huge project so it's promising!

Drone env with go-jsonnet:

```
root@0447d01aca66:/Users/julienduchesne/Repos/tanka# time ./tk eval --implementation go /Users/julienduchesne/Repos/deployment_tools/ksonnet/environments/drone/ | wc -l
8191

real    0m5.515s
user    0m6.995s
sys     0m0.440s
```

Drone env with jrsonnet:

```
root@0447d01aca66:/Users/julienduchesne/Repos/tanka# time ./tk eval --implementation rust /Users/julienduchesne/Repos/deployment_tools/ksonnet/environments/drone/ | wc -l
8191

real    0m1.123s
user    0m0.242s
sys     0m0.148s
```

Drone env with go-jsonnet:

```
root@0447d01aca66:/Users/julienduchesne/Repos/tanka# time ./tk eval --implementation go /Users/julienduchesne/Repos/deployment_tools/ksonnet/environments/grafana-o11y/ | wc -l
44860

real    7m25.873s
user    10m40.486s
sys     1m5.636s
```

Drone env with jrsonnet:

```
root@0447d01aca66:/Users/julienduchesne/Repos/tanka# time ./tk eval --implementation rust /Users/julienduchesne/Repos/deployment_tools/ksonnet/environments/grafana-o11y/ | wc -l
44860

real    0m41.603s
user    0m27.297s
sys     0m9.535s
```
