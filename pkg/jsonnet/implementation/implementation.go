package implementation

import (
	"github.com/grafana/tanka/pkg/jsonnet/implementation/goimpl"
	"github.com/grafana/tanka/pkg/jsonnet/implementation/types"
)

func Get(name string) types.JsonnetImplementation {
	switch name {
	default:
		return &goimpl.JsonnetGoImplementation{}
	}
}
