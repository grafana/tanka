{
  base_env: {
    apiVersion: 'tanka.dev/v1alpha1',
    kind: 'Environment',
    metadata: {
      name: 'base',
    },
    spec: {
      apiServer: 'https://localhost',
      namespace: 'base',
    },
    data: {
      apiVersion: 'v1',
      kind: 'ConfigMap',
      metadata: { name: 'config' },
    },
  },
  other_env: {
    apiVersion: 'tanka.dev/v1alpha1',
    kind: 'Environment',
    metadata: {
      name: 'base-and-more',
    },
    spec: {
      apiServer: 'https://localhost',
      namespace: 'base-and-more',
    },
    data: {
      apiVersion: 'v1',
      kind: 'ConfigMap',
      metadata: { name: 'config' },
    },
  },
}
