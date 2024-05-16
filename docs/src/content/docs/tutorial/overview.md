---
title: Overview
sidebar:
  order: 1
---

## Learning how to use Tanka

Welcome to the Tanka tutorial!
The following sections will explain how to deploy an example stack,
([Grafana](https://hub.docker.com/r/grafana/grafana) and
[Prometheus](https://hub.docker.com/r/prom/prometheus)), to Kubernetes. We will also deal with parameters, differences between `dev` and `prod` and how to stop worrying and love libraries.

To do so, we have the following steps:

1. [Deploying **without** Tanka first](./tutorial/refresher): Using good old `kubectl` to understand what Tanka will do for us.
2. [Using Jsonnet](./tutorial/jsonnet): Doing the same thing once again, but this time with Tanka and Jsonnet.
3. [Parameterizing](./tutorial/parameters): Using Variables to avoid data duplication.
4. [Abstraction](./tutorial/abstraction): Splitting components into individual parts.
5. [Environments](./tutorial/environments): Dealing with differences between `dev` and `prod`.
6. [`k.libsonnet`](./tutorial/k-lib): Avoid having to remember API resources.

Completing this gives a solid knowledge of Tanka's fundamentals. Let's get started!

## Resources

- The final outcome of this tutorial can be seen here:
  [https://github.com/grafana/tanka/examples/prom-grafana](https://github.com/grafana/tanka/tree/main/examples/prom-grafana)
