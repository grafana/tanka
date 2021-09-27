{
  project1_env1: {
    apiVersion: 'tanka.dev/v1alpha1',
    kind: 'Environment',
    metadata: {
      name: 'project1-env1',
    },
    spec: {
      apiServer: 'https://localhost',
      namespace: 'project1-env1',
    },
    data: {
      apiVersion: 'v1',
      kind: 'ConfigMap',
      metadata: { name: 'config' },
    },
  },
  project1_env2: {
    apiVersion: 'tanka.dev/v1alpha1',
    kind: 'Environment',
    metadata: {
      name: 'project1-env2',
    },
    spec: {
      apiServer: 'https://localhost',
      namespace: 'project1-env2',
    },
    data: {
      apiVersion: 'v1',
      kind: 'ConfigMap',
      metadata: { name: 'config' },
    },
  },
  project2_env1: {
    apiVersion: 'tanka.dev/v1alpha1',
    kind: 'Environment',
    metadata: {
      name: 'project2-env1',
    },
    spec: {
      apiServer: 'https://localhost',
      namespace: 'project2-env1',
    },
    data: {
      apiVersion: 'v1',
      kind: 'ConfigMap',
      metadata: { name: 'config' },
    },
  },
}
