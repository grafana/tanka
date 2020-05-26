local backendDeploy = (import './resources.jsonnet').deployment;
local namespace = (import './resources.jsonnet').namespace;
local frontendService = (import './k8s.libsonnet').service(name='frontend');
local frontendDeploy = (import './k8s.libsonnet').deployment(name='frontend');

{
  deep: {
    app: {
      namespace: namespace,
      web: {
        backend: { server: { grafana: {
          deployment: backendDeploy,
        } } },
        frontend: { nodejs: { express: {
          service: frontendService,
          deployment: frontendDeploy,
        } } },
      },
    },
  },
  flat: {
    '.app.web.backend.server.grafana.deployment': backendDeploy,
    '.app.web.frontend.nodejs.express.service': frontendService,
    '.app.web.frontend.nodejs.express.deployment': frontendDeploy,
    '.app.namespace': namespace,
  },
}
