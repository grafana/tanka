package client

import (
	"testing"

	"k8s.io/apimachinery/pkg/util/sets"
)

func TestKubectl_deleteCtl(t *testing.T) {
	info := Info{
		Kubeconfig: Config{
			Context: Context{
				Name: "foo-context",
			},
		},
	}

	type args struct {
		ns   string
		kind string
		name string
		opts DeleteOpts
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
				ns:   "foo-ns",
				kind: "deploy",
				name: "foo-deploy",
				opts: DeleteOpts{},
			},
			expectedArgs:   []string{"--context", info.Kubeconfig.Context.Name, "-n", "foo-ns", "deploy", "foo-deploy"},
			unExpectedArgs: []string{"--force", "--dry-run=server"},
		},
		{
			name: "test dry-run",
			args: args{
				opts: DeleteOpts{DryRun: "server"},
			},
			expectedArgs: []string{"--dry-run=server"},
		},
		{
			name: "test force",
			args: args{
				opts: DeleteOpts{Force: true},
			},
			expectedArgs: []string{"--force"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := Kubectl{
				info: info,
			}
			got := k.deleteCtl(tt.args.ns, tt.args.kind, tt.args.name, tt.args.opts)
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
