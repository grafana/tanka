package telemetry

import (
	"fmt"

	"github.com/grafana/tanka/pkg/spec/v1alpha1"
	"go.opentelemetry.io/otel/attribute"
)

func AttrPath(v string) attribute.KeyValue {
	return attribute.String("tanka.path", v)
}

func AttrEnv(v *v1alpha1.Environment) []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("tanka.env.id", fmt.Sprintf("%s@%s", v.Metadata.Name, v.Spec.APIServer)),
	}
}
