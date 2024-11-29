function(bar='bar', baz='baz') {
  apiVersion: 'tanka.dev/v1alpha1',
  kind: 'Environment',
  metadata: {
    name: bar + '-' + baz,
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
