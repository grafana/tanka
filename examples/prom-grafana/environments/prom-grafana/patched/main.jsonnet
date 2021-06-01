(import 'prom-grafana/prom-grafana.libsonnet') +
{
  promgrafana+: {
    prometheus+: {
      deployment+: {
        metadata+: {
          labels+: {
            foo: 'bar',
          },
        },
      },
    },
  },
}
