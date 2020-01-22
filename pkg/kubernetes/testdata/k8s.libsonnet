// a very basic ripoff from `k.libsonnet`, because we can't vendor in tests

{
  deployment(name="grafana", image="grafana/grafana"):: {
    apiVersion: "apps/v1",
    kind: "Deployment",
    metadata: { name: name },
    spec: {
      replicas: 1,
      template: {
        containers: [{
            name: name,
            image: image,
        }],
        metadata: { labels: { app: name }}
      }
    }
  },
  service(name="grafana", image="grafana/grafana"):: {
    apiVersion: "v1",
    kind: "Service",
    metadata: { name: name },
    spec: {
      selector: { app: name },
      ports: [{
        name: name,
        port: 3000,
        targetPort: 3000
      }]
    }
  },
  namespace(name="default"):: {
    apiVersion: "v1",
    kind: "Namespace",
    metadata: { name: name }
  }
}
