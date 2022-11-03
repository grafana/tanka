{
  // jsonnet supports going one level lower than files really are
  first: import '../relative-imported/main.jsonnet',
  second: import '../../relative-imported2/main.jsonnet',

  externalFile: importstr '../../other-files/test.txt',
  externalFile2: importstr '../../../other-files/test2.txt',
}
