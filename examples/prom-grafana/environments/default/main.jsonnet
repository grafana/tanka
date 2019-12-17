local k = import "ksonnet-util/kausal.libsonnet";

k + {
  local deployment = $.apps.v1.deployment,
  local container = $.core.v1.container,
  local port = $.core.v1.containerPort,
  local service = $.core.v1.service,

  prometheus: {
    deployment: deployment.new(
      name="prometheus", replicas=1,
      containers=[
        container.new("prometheus", "prom/prometheus")
        + container.withPorts([port.new("api", 9090)]),
      ],
    ),
    service: $.util.serviceFor(self.deployment),
  },
  grafana: {
    deployment: deployment.new(
      name="grafana", replicas=1,
      containers=[
        container.new("grafana", "grafana/grafana")
        + container.withPorts([port.new("ui", 3000)]),
      ],
    ),
    service: $.util.serviceFor(self.deployment) + service.mixin.spec.withType("NodePort"),
  },
}
