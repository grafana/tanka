package implementation

import (
	"github.com/grafana/tanka/pkg/jsonnet/implementation/goimpl"
	"github.com/grafana/tanka/pkg/jsonnet/implementation/rustimpl"
	"github.com/grafana/tanka/pkg/jsonnet/implementation/types"
)

func Get(name string) (types.JsonnetImplementation, error) {
	switch name {
	case "rust":
		return &rustimpl.JsonnetRustImplementation{}, nil
	case "go":
		return &goimpl.JsonnetGoImplementation{}, nil
	}

	// TODO: fix
	return &rustimpl.JsonnetRustImplementation{}, nil
	// return nil, fmt.Errorf("unknown implementation: %s", name)
}
