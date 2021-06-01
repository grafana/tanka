local k = import 'ksonnet-util/kausal.libsonnet';

(import './config.libsonnet') +
{
  local deployment = k.apps.v1.deployment,
  local container = k.core.v1.container,
  local port = k.core.v1.containerPort,
  local service = k.core.v1.service,

  // alias our params, too long to type every time
  local c = $._config.promgrafana,

  promgrafana: {
    prometheus: {
      deployment: deployment.new(
        name=c.prometheus.name,
        replicas=1,
        containers=[
          container.new(c.prometheus.name, $._images.promgrafana.prometheus)
          + container.withPorts([port.new('api', c.prometheus.port)]),
        ],
      ),
      service: k.util.serviceFor(self.deployment),
    },

    grafana: {
      deployment: deployment.new(
        name=c.grafana.name,
        replicas=1,
        containers=[
          container.new(c.grafana.name, $._images.promgrafana.grafana)
          + container.withPorts([port.new('ui', c.grafana.port)]),
        ],
      ),
      service:
        k.util.serviceFor(self.deployment)
        + service.mixin.spec.withType('NodePort'),
    },
  },
}
