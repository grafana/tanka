package process

import (
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"reflect"
	"regexp"
	"testing"
)

func TestFilterWithOptions(t *testing.T) {
	type args struct {
		list    manifest.List
		options []FilterOption
	}
	type test struct {
		name string
		args args
		want manifest.List
	}
	tests := []test{
		func() test {
			tdef := test{
				name: "Given no filters, no selector, expect all manifests returned",
				args: args{},
			}

			list := manifest.List{
				{
					"kind":       "Pod",
					"apiVersion": "v1",
					"metadata": map[string]interface{}{
						"name": "foo",
						"labels": map[string]interface{}{
							"app": "foo",
						},
					},
				},
				{
					"kind":       "Pod",
					"apiVersion": "v1",
					"metadata": map[string]interface{}{
						"name": "bar",
						"labels": map[string]interface{}{
							"app": "foo",
						},
					},
				},
			}

			tdef.args.list = list
			tdef.want = list

			return tdef
		}(),
		{
			name: "Given kind filter on pod, no manifests returned",
			args: args{
				list: manifest.List{
					{
						"kind":       "Deployment",
						"apiVersion": "v1",
						"metadata": map[string]interface{}{
							"name": "foo",
						},
					},
					{
						"kind":       "Deployment",
						"apiVersion": "v1",
						"metadata": map[string]interface{}{
							"name": "bar",
						},
					},
				},
				options: []FilterOption{
					func(options *FilterOptions) {
						options.Exprs = RegExps([]*regexp.Regexp{regexp.MustCompile("Pod/.*")})
					},
				},
			},
			want: manifest.List{},
		},
		func() test {
			tdef := test{
				name: "Given selector, expect app=foo manifest returned",
				args: args{},
			}

			expectedKey := "app"
			expectedValue := "foo"

			requirement, err := labels.NewRequirement(expectedKey, selection.In, []string{expectedValue})

			if err != nil {
				t.Errorf("failed to create requirement: %s", err)
			}

			selector := labels.NewSelector().Add(*requirement)

			tdef.args.options = []FilterOption{
				func(options *FilterOptions) {
					options.Selector = selector
				},
			}

			expected := manifest.Manifest{
				"kind":       "Pod",
				"apiVersion": "v1",
				"metadata": map[string]interface{}{
					"name": "foo",
					"labels": map[string]interface{}{
						expectedKey: expectedValue,
					},
				},
			}

			tdef.args.list = manifest.List{
				expected,
				{
					"kind":       "Pod",
					"apiVersion": "v1",
					"metadata": map[string]interface{}{
						"name": "bar",
						"labels": map[string]interface{}{
							"app": "bar",
						},
					},
				},
			}
			tdef.want = manifest.List{
				expected,
			}

			return tdef
		}(),
		func() test {
			tdef := test{
				name: "Given filter and selector, expect pod app=foo manifest returned",
				args: args{},
			}

			expectedKey := "app"
			expectedValue := "foo"

			requirement, err := labels.NewRequirement(expectedKey, selection.In, []string{expectedValue})

			if err != nil {
				t.Errorf("failed to create requirement: %s", err)
			}

			selector := labels.NewSelector().Add(*requirement)

			tdef.args.options = []FilterOption{
				func(options *FilterOptions) {
					options.Exprs = RegExps([]*regexp.Regexp{regexp.MustCompile("Pod/.*")})
					options.Selector = selector
				},
			}

			expected := manifest.Manifest{
				"kind":       "Pod",
				"apiVersion": "v1",
				"metadata": map[string]interface{}{
					"name": "foo",
					"labels": map[string]interface{}{
						expectedKey: expectedValue,
					},
				},
			}

			tdef.args.list = manifest.List{
				expected,
				{
					"kind":       "Deployment",
					"apiVersion": "v1",
					"metadata": map[string]interface{}{
						"name": "bar",
						"labels": map[string]interface{}{
							expectedKey: expectedValue,
						},
					},
				},
				{
					"kind":       "Pod",
					"apiVersion": "v1",
					"metadata": map[string]interface{}{
						"name": "bar",
						"labels": map[string]interface{}{
							"app": "bar",
						},
					},
				},
			}
			tdef.want = manifest.List{
				expected,
			}

			return tdef
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterWithOptions(tt.args.list, tt.args.options...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterWithOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}
