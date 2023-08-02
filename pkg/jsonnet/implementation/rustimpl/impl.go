package rustimpl

import (
	"github.com/grafana/tanka/pkg/jsonnet/implementation/types"
	"github.com/grafana/tanka/pkg/jsonnet/native"
)

// EvaluateAnonymousSnippet does the same as `EvaluateSnippet`.
// This is an alias to respect the interface and keep the CGO code intact from upstream.
func (vm *VM) EvaluateAnonymousSnippet(filename, snippet string) (string, error) {
	return vm.EvaluateSnippet(filename, snippet)
}

type JsonnetRustImplementation struct{}

func (i *JsonnetRustImplementation) MakeVM(importPaths []string, extCode map[string]string, tlaCode map[string]string, maxStack int) types.JsonnetVM {
	vm := Make()
	for i := len(importPaths) - 1; i >= 0; i-- {
		vm.JpathAdd(importPaths[i])
	}
	for key, value := range extCode {
		vm.ExtCode(key, value)
	}
	for key, value := range tlaCode {
		vm.TlaCode(key, value)
	}
	if maxStack > 0 {
		vm.MaxStack(uint(maxStack))
	}
	for _, nf := range native.Funcs() {
		params := []string{}
		for _, p := range nf.Params {
			params = append(params, string(p))
		}
		vm.NativeCallback(nf.Name, params, nf.Func)
	}

	return vm
}
