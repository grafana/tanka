{
  apiVersion: 'tanka.dev/v1alpha1',
  kind: 'Environment',
  metadata: {
    name: 'withenv',
  },
  spec: {
    apiServer: 'https://localhost',
    namespace: 'withenv',
  },
  data: {
    testCase: 'object',
  },
}
