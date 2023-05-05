package types

type JsonnetVM interface {
	EvaluateAnonymousSnippet(filename, snippet string) (string, error)
	EvaluateFile(filename string) (string, error)
}

type JsonnetImplementation interface {
	MakeVM(importPaths []string, extCode map[string]string, tlaCode map[string]string, maxStack int) JsonnetVM
}
