{
  apiVersion: 'tanka.dev/v1alpha1',
  kind: 'Environment',
  metadata: {
    name: 'withenv',
  },
  spec: {
    apiServer: 'https://localhost',
    namespace: 'withenv',
    exportPriority: 3,
  },
  data: {
    apiVersion: 'v1',
    kind: 'ConfigMap',
    metadata: { name: 'config' },
  },
}
