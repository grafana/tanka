package goimpl

import (
	"github.com/google/go-jsonnet"
	"github.com/grafana/tanka/pkg/jsonnet/implementation/types"
)

type JsonnetGoVM struct {
	vm *jsonnet.VM
}

func (vm *JsonnetGoVM) EvaluateAnonymousSnippet(filename, snippet string) (string, error) {
	return vm.vm.EvaluateAnonymousSnippet(filename, snippet)
}

func (vm *JsonnetGoVM) EvaluateFile(filename string) (string, error) {
	return vm.vm.EvaluateFile(filename)
}

type JsonnetGoImplementation struct{}

func (i *JsonnetGoImplementation) MakeVM(importPaths []string, extCode map[string]string, tlaCode map[string]string, maxStack int) types.JsonnetVM {
	return &JsonnetGoVM{
		vm: MakeRawVM(importPaths, extCode, tlaCode, maxStack),
	}
}
