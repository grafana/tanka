{
  assert self.imported.test == 'env-vendor',
  imported: (import 'vendor-override-in-env/main.libsonnet'),
}
