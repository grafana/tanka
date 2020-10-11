{
  envs: [
    {
      apiVersion: 'tanka.dev/v1alpha1',
      kind: 'Environment',
      metadata: {
        name: 'withenv1',
      },
      spec: {
        apiServer: 'https://localhost',
        namespace: 'withenv',
      },
      data: {
        testCase: 'object',
      },
    },
    {
      apiVersion: 'tanka.dev/v1alpha1',
      kind: 'Environment',
      metadata: {
        name: 'withenv2',
      },
      spec: {
        apiServer: 'https://localhost',
        namespace: 'withenv',
      },
      data: {
        testCase: 'object',
      },
    },
  ],
}
