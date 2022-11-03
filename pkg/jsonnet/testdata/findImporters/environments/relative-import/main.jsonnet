{
  // jsonnet supports going one level lower than files really are
  first: import '../relative-imported/main.jsonnet',
  second: import '../../relative-imported2/main.jsonnet',
}
