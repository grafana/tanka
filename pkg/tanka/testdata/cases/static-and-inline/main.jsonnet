{
  apiVersion: 'tanka.dev/v1alpha1',
  kind: 'Environment',
  metadata: {
    name: 'inline',
  },
  spec: {
    apiServer: 'https://localhost',
    namespace: 'inline',
  },
  data: {
    apiVersion: 'v1',
    kind: 'ConfigMap',
    metadata: { name: 'config' },
  },
}
