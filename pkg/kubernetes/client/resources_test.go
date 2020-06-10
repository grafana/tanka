package client

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grafana/tanka/pkg/kubernetes/resources"
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
			want: &resources.Store{
				{APIGroup: "apps", Name: "Deployment", Namespaced: true},
				{APIGroup: "networking", Name: "Ingress", Namespaced: true},
				{APIGroup: "", Name: "Namespace", Namespaced: false},
				{APIGroup: "extensions", Name: "DaemonSet", Namespaced: true},
			},
			dest: &resources.Store{},
		},
		{
			name: "empty",
			tbl:  strings.TrimSpace(tblEmpty),
			want: &resources.Store{},
			dest: &resources.Store{
				{APIGroup: "apps", Name: "Deployment", Namespaced: true},
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
APIGROUP    NAME        NAMESPACED
apps        Deployment  true
networking  Ingress     true
            Namespace   false
extensions  DaemonSet   true
`

var tblEmpty = `
APIGROUP    NAME        NAMESPACED
`

var tblNoHeader = `
apps        Deployment  true
networking  Ingress     true
`

var tblNothing = ``
