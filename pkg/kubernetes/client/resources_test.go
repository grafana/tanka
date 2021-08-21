package client

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalTable(t *testing.T) {
	cases := []struct {
		name string
		tbl  string
		dest interface{}
		want interface{}
		err  error
	}{
		{
			name: "normal",
			tbl:  strings.TrimSpace(tblNormal),
			want: &Resources{
				{APIVersion: "apps/v1", Name: "Deployment", Namespaced: true},
				{APIVersion: "networking/v1", Name: "Ingress", Namespaced: true},
				{APIVersion: "v1", Name: "Namespace", Namespaced: false},
				{APIVersion: "extensions/v1", Name: "DaemonSet", Namespaced: true},
			},
			dest: &Resources{},
		},
		{
			name: "empty",
			tbl:  strings.TrimSpace(tblEmpty),
			want: &Resources{},
			dest: &Resources{
				{APIVersion: "apps/v1", Name: "Deployment", Namespaced: true},
			},
		},
		{
			name: "no-header",
			tbl:  strings.TrimSpace(tblNoHeader),
			err:  ErrorNoHeader,
		},
		{
			name: "nothing",
			tbl:  tblNothing,
			err:  ErrorNoHeader,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := UnmarshalTable(c.tbl, c.dest)
			require.Equal(t, c.err, err)
			assert.Equal(t, c.want, c.dest)
		})
	}
}

var tblNormal = `
APIVERSION     NAME        NAMESPACED
apps/v1        Deployment  true
networking/v1  Ingress     true
v1             Namespace   false
extensions/v1  DaemonSet   true
`

var tblEmpty = `
APIVERSION    NAME        NAMESPACED
`

var tblNoHeader = `
apps        Deployment  true
networking  Ingress     true
`

var tblNothing = ``
