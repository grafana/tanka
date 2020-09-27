<p align="center">
  <img
    width="400"
    src="https://raw.githubusercontent.com/grafana/tanka/master/docs/img/logo.svg"
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
  src="https://raw.githubusercontent.com/grafana/tanka/master/docs/img/example.png"
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
- **:rocket: Production ready**: Tanka deploys [Grafana Cloud](https://grafana.com/cloud) and many more production setups

<br />
<p align="center">
  <a href="https://tanka.dev/tutorial/overview"><strong>Let's kill some YAML together&nbsp;&nbsp;▶</strong></a>
</p>

Another Grafana Tanka community call is coming up!

- :calendar: Mark the date: 2020-10-06 16:00 UTC
- :tv: Join the meet in this doc: https://bit.ly/3czxbFz

## :rocket: Getting started

To get started, [install Tanka](https://tanka.dev/install) first, and then
[follow the tutorial](https://tanka.dev/tutorial/overview). This should get you
on track quickly.

## :busts_in_silhouette: Community

The Tanka community is core to the project. Connect to us using the `#tanka`
channel on the [Grafana Slack](https://grafana.slack.com).

---

[![Community
Call](./docs/img/community-call.png)](https://docs.google.com/document/d/1mEsc0GxlnwbWAXzbIP7tBb6T5WgAI66_0gIWJzqB93o/edit)

Grafana Labs hosts a monthly community call for the Tanka project. Notes and the
meeting URL can be found in a [Google
Doc](https://docs.google.com/document/d/1mEsc0GxlnwbWAXzbIP7tBb6T5WgAI66_0gIWJzqB93o/edit).

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
