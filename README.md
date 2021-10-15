<p align="center">
  <img
    width="400"
    src="https://raw.githubusercontent.com/grafana/tanka/main/docs/img/logo.svg"
    alt="Grafana Tanka Logo"
  />
</p>

<p align="center">
  <a href="https://cloud.drone.io/grafana/tanka">
    <img src="https://img.shields.io/drone/build/grafana/tanka?style=flat-square&server=https%3A%2F%2Fdrone.grafana.net">
  </a>
  <a href="https://github.com/grafana/tanka/releases">
    <img src="https://img.shields.io/github/release/grafana/tanka?style=flat-square" />
  </a>
  <img src="https://img.shields.io/github/contributors/grafana/tanka?style=flat-square" />
  <a href="https://grafana.slack.com">
    <img src="https://img.shields.io/badge/Slack-GrafanaLabs-orange?logo=slack&style=flat-square" />
  </a>
</p>

<p align="center">
  <a href="https://tanka.dev">Website</a>
  ·
  <a href="https://tanka.dev/install">Installation</a>
  ·
  <a href="https://tanka.dev/tutorial/overview">Tutorial</a>
</p>

# Grafana Tanka

<img
  src="https://raw.githubusercontent.com/grafana/tanka/main/docs/img/example.png"
  width="50%"
  align="right"
/>

**The clean, concise and super flexible alternative to YAML for your
[Kubernetes](https://k8s.io) cluster**

- **:boom: Clean**: The
  [Jsonnet language](https://jsonnet.org) expresses your apps more obviously than YAML ever did
- **:books: Reusable**: Build libraries, import them anytime and even share them on GitHub!
- **:pushpin: Concise**: Using the Kubernetes library and abstraction, you will
  never see boilerplate again!
- **:dart: Confidence**: Stop guessing and use `tk diff` to see what exactly will happen
- **:telescope: Helm**: Vendor in, modify, and export [Helm charts reproducibly](https://tanka.dev/helm#helm-support)
- **:rocket: Production ready**: Tanka deploys [Grafana Cloud](https://grafana.com/cloud) and many more production setups

<br />
<p align="center">
  <a href="https://tanka.dev/tutorial/overview"><strong>Let's kill some YAML together&nbsp;&nbsp;▶</strong></a>
</p>

## :rocket: Getting started

The very :1st_place_medal: first step of getting started is to [install Tanka](https://tanka.dev/install). 
Then follow the Tutorial to get on the track quickly.

  ### <span> <img src = 'https://raw.githubusercontent.com/grafana/tanka/main/docs/img/logo.svg' width = '30' > </span> Turorial <br>
  The following sections will explain how to deploy an example stack, ([Grafana](https://hub.docker.com/r/grafana/grafana) and [Prometheus](https://hub.docker.com/r/prom/prometheus)), to Kubernetes. 
  We will also deal with parameters, differences between dev and prod and how to stop worrying and love libraries.
  
  To do so, we have the following steps:

  - **[Deploying without Tanka first](https://tanka.dev/tutorial/refresher)**: Using good old kubectl to understand what Tanka will do for us.
  - **[Using Jsonnet](https://tanka.dev/tutorial/jsonnet)**: Doing the same thing once again, but this time with Tanka and Jsonnet.
  - **[Parameterizing](https://tanka.dev/tutorial/parameters)**: Using Variables to avoid data duplication.
  - **[Abstraction](https://tanka.dev/tutorial/abstraction)**: Splitting components into individual parts.
  - **[Environments](https://tanka.dev/tutorial/environments)**: Dealing with differences between dev and prod.
  - **[k.libsonnet](https://tanka.dev/tutorial/k-lib)**: Avoid having to remember API resources.


  Completing this gives a solid knowledge of Tanka's fundamentals. Let's get started!
  
  For more information, visit [Tanka Tutorial](https://tanka.dev/tutorial/overview). The final outcome of this tutorial can be seen here: [example/prom-grafana](https://github.com/grafana/tanka/examples/prom-grafana)

  

## :busts_in_silhouette: Community

There are several places to connect with the Tanka community:

- [GitHub Discussions](https://github.com/grafana/tanka/discussions/442): Primary support channel
- `#tanka` on [Grafana Slack](https://grafana.slack.com)
- Mailing lists
  - [`tanka-announce`](https://groups.google.com/forum/#!forum/tanka-announce):
    Low frequency list with announcements, releases, etc
  - [`tanka-users`](https://groups.google.com/forum/#!forum/tanka-users):
    General purpose group for discussions, community support and more

Please don't ask individual project members or open GitHub issues for support
requests. Use one of the above channels so everyone in the community can
participate.

Furthermore, see [`LICENSE`](./LICENSE) and [`GOVERNANCE`](./GOVERNANCE.md).

## :book: Additional resources

- https://jsonnet.org/, the official Jsonnet documentation provides lots of
  examples on how to use the language.
- https://github.com/grafana/jsonnet-libs: Grafana Labs' Jsonnet libraries are a
  rich set of configuration examples compatible with Tanka.

## :pencil: License

Tanka is an open-source project :heart:. It is free as
in beer and as in speech and this will never change.

Licensed under Apache 2.0, see [LICENSE](LICENSE).
