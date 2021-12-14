package client

import (
	"testing"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
)

func TestKubectl_applyCtl(t *testing.T) {
	info := Info{
		Kubeconfig: Config{
			Context: Context{
				Name: "foo-context",
			},
		},
	}

	type args struct {
		data manifest.List
		opts ApplyOpts
	}
	tests := []struct {
		name           string
		args           args
		expectedArgs   []string
		unExpectedArgs []string
	}{
		{
			name: "test default",
			args: args{
				opts: ApplyOpts{Validate: true},
			},
			expectedArgs:   []string{"--context", info.Kubeconfig.Context.Name},
			unExpectedArgs: []string{"--force", "--dry-run=server", "--validate=false"},
		},
		{
			name: "test force",
			args: args{
				opts: ApplyOpts{Validate: true, Force: true},
			},
			expectedArgs:   []string{"--force"},
			unExpectedArgs: []string{"--validate=false"},
		},
		{
			name: "test validate",
			args: args{
				opts: ApplyOpts{Validate: false},
			},
			expectedArgs: []string{"--validate=false"},
		},
		{
			name: "test dry-run",
			args: args{
				opts: ApplyOpts{Validate: true, DryRun: "server"},
			},
			expectedArgs:   []string{"--dry-run=server"},
			unExpectedArgs: []string{"--validate=false"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := Kubectl{
				info: info,
			}
			got := k.applyCtl(tt.args.data, tt.args.opts)
			gotSet := sets.NewString(got.Args...)
			if !gotSet.HasAll(tt.expectedArgs...) {
				t.Errorf("Kubectl.applyCtl() = %v doesn't have (all) expectedArgs='%v'", got.Args, tt.expectedArgs)
			}
			if gotSet.HasAny(tt.unExpectedArgs...) {
				t.Errorf("Kubectl.applyCtl() = %v has (any) unExpectedArgs='%v'", got.Args, tt.unExpectedArgs)
			}

		})
	}
}
