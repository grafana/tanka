(import "ksonnet-util/kausal.libsonnet") +
{
  _config:: {
    grafana: {
      port: 3000,
      name: "grafana",
    },
    prometheus: {
      port: 9090,
      name: "prometheus"
    }
  },

  local deployment = $.apps.v1.deployment,
  local container = $.core.v1.container,
  local port = $.core.v1.containerPort,
  local service = $.core.v1.service,

  prometheus: {
    deployment: deployment.new(
      name=$._config.prometheus.name, replicas=1,
      containers=[
        container.new($._config.prometheus.name, "prom/prometheus")
        + container.withPorts([port.new("api", $._config.prometheus.port)]),
      ],
    ),
    service: $.util.serviceFor(self.deployment),
  },
  grafana: {
    deployment: deployment.new(
      name=$._config.grafana.name, replicas=1,
      containers=[
        container.new($._config.grafana.name, "grafana/grafana")
        + container.withPorts([port.new("ui", $._config.grafana.port)]),
      ],
    ),
    service: $.util.serviceFor(self.deployment) + service.mixin.spec.withType("NodePort"),
  },
}
