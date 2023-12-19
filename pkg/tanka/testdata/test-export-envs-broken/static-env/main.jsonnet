{
  deployment: {
    apiVersion: 'apps/v1',
    kind: 'Deployment',
    metadata: {
      name: 'foo',
    },
  },
  service: {
    apiVersion: 'v1',
    kind: 'Service',
    metadata: {
      // Error, this should be a string
      name: true,
    },
  },
}
