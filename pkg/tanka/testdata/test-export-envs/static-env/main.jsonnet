{
  deployment: {
    apiVersion: 'apps/v1',
    kind: 'Deployment',
    metadata: {
      name: std.extVar('deploymentName'),
    },
  },
  service: {
    apiVersion: 'v1',
    kind: 'Service',
    metadata: {
      name: std.extVar('serviceName'),
    },
  },
}
