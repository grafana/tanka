package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubset(t *testing.T) {
	tests := []struct {
		name       string
		should, is map[string]interface{}
		want       map[string]interface{}
	}{
		{
			name: "simple",
			should: map[string]interface{}{
				"foo": "bar",
				"bam": "boo",
			},
			is: map[string]interface{}{
				"foo": "baz",
				"baz": "bar",
			},
			want: map[string]interface{}{
				"foo": "baz",
			},
		},
		{
			name: "nested",
			should: map[string]interface{}{
				"foo": "bar",
				"baz": map[string]interface{}{
					"foo": "bar",
					"bar": "boo",
				},
			},
			is: map[string]interface{}{
				"foo": "bam",
				"baz": map[string]interface{}{
					"rab": "bar",
					"bar": "foo",
				},
			},
			want: map[string]interface{}{
				"foo": "bam",
				"baz": map[string]interface{}{
					"bar": "foo",
				},
			},
		},
		{
			name: "slice",
			should: map[string]interface{}{
				"foo": []map[string]interface{}{
					{
						"foo": "bar",
					},
				},
			},
			is: map[string]interface{}{
				"foo": []map[string]interface{}{
					{
						"foo": "baz",
						"bam": "baz",
					},
				},
			},
			want: map[string]interface{}{
				"foo": []map[string]interface{}{
					{
						"foo": "baz",
					},
				},
			},
		},
		{
			name: "heterogeneous_slice",
			should: map[string]interface{}{
				"foo": []map[string]interface{}{
					{
						"foo": "bar",
						"het": []interface{}{
							"foobar",
							map[string]interface{}{
								"bam": "baz",
							},
						},
					},
				},
			},
			is: map[string]interface{}{
				"foo": []map[string]interface{}{
					{
						"foo": "baz",
						"het": []interface{}{
							"foobam",
							map[string]interface{}{
								"bam": "bloo",
								"boo": "boar",
							},
							map[string]interface{}{
								"a": "b",
							},
						},
					},
				},
			},
			want: map[string]interface{}{
				"foo": []map[string]interface{}{
					{
						"foo": "baz",
						"het": []interface{}{
							"foobam",
							map[string]interface{}{
								"bam": "bloo",
							},
							map[string]interface{}{
								"a": "b",
							},
						},
					},
				},
			},
		},
		{
			name: "namespace",
			should: map[string]interface{}{
				"namespace": "loki",
			},
			is: map[string]interface{}{},
			want: map[string]interface{}{
				"namespace": "loki",
			},
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, subset(c.should, c.is))
		})
	}
}
