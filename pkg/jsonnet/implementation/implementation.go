package implementation

import (
	"fmt"

	"github.com/grafana/tanka/pkg/jsonnet/implementation/goimpl"
	"github.com/grafana/tanka/pkg/jsonnet/implementation/types"
)

func Get(name string) (types.JsonnetImplementation, error) {
	if name == "go" || name == "" {
		return &goimpl.JsonnetGoImplementation{}, nil
	}

	return nil, fmt.Errorf("unknown jsonnet implementation: %s", name)
}
