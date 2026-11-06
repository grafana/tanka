local missing = import 'does-not-exist.libsonnet';
{ source: "root", x: missing }
