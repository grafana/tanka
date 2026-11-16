// lib/prefix-dir.suffix also has an a.libsonnet and a b.libsonnet,
// but this environment never imports them
{
  imported: import 'prefix-dir/a.libsonnet',
}
