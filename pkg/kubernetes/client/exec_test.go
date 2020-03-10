package client

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const patchFile = "/tmp/tk-nsPatch.yaml"

func TestPatchKubeconfig(t *testing.T) {
	cases := []struct {
		name string
		env  []string
		want []string
	}{
		{
			name: "none",
			env:  []string{},
			want: []string{
				fmt.Sprintf("KUBECONFIG=%s:%s", patchFile, filepath.Join(homeDir(), ".kube", "config")),
			},
		},
		{
			name: "custom",
			env:  []string{"KUBECONFIG=/home/user/.config/kube"},
			want: []string{"KUBECONFIG=" + patchFile + ":/home/user/.config/kube"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := patchKubeconfig(patchFile, c.env)
			assert.Equal(t, c.want, got)
		})
	}
}
