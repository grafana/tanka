package client

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func collectFQNs(resIntf interface{}) []string {
	resPtr, ok := resIntf.(*Resources)
	if !ok {
		return nil
	}
	res := *resPtr
	if len(res) == 0 {
		return nil
	}
	out := make([]string, len(res))
	for pos := range res {
		out[pos] = res[pos].FQN()
	}
	return out
}

func TestUnmarshalTable(t *testing.T) {
	cases := []struct {
		name     string
		tbl      string
		dest     interface{}
		want     interface{}
		wantFQNs []string
		err      error
	}{
		{
			name: "normal",
			tbl:  strings.TrimSpace(tblv126),
			want: &Resources{
				{APIVersion: "v1", Kind: "Binding", Name: "bindings", Shortnames: "", Namespaced: true, Verbs: "create"},
				{APIVersion: "v1", Kind: "ComponentStatus", Name: "componentstatuses", Shortnames: "cs", Namespaced: false, Verbs: "get,list"},
				{APIVersion: "v1", Kind: "ConfigMap", Name: "configmaps", Shortnames: "cm", Namespaced: true, Verbs: "create,delete,deletecollection,get,list,patch,update,watch"},
				{APIVersion: "v1", Kind: "Pod", Name: "pods", Shortnames: "po", Namespaced: true, Verbs: "create,delete,deletecollection,get,list,patch,update,watch", Categories: "all"},
			},
			wantFQNs: []string{
				"bindings",
				"componentstatuses",
				"configmaps",
				"pods",
			},
			dest: &Resources{},
		},
		{
			name: "normal-v1.18-to-v.1.25",
			tbl:  strings.TrimSpace(tblv121),
			want: &Resources{
				{APIVersion: "v1", Kind: "Namespace", Name: "namespaces", Shortnames: "ns", Namespaced: false},
				{APIVersion: "apps/v1", Kind: "DaemonSet", Name: "daemonsets", Shortnames: "ds", Namespaced: true},
				{APIVersion: "apps/v1", Kind: "Deployment", Name: "deployments", Shortnames: "deploy", Namespaced: true},
				{APIVersion: "networking.k8s.io/v1", Kind: "Ingress", Name: "ingresses", Shortnames: "ing", Namespaced: true},
			},
			wantFQNs: []string{
				"namespaces",
				"daemonsets.apps",
				"deployments.apps",
				"ingresses.networking.k8s.io",
			},
			dest: &Resources{},
		},
		{
			name: "normal-v1.18-and-older",
			tbl:  strings.TrimSpace(tblv118),
			want: &Resources{
				{APIGroup: "", Kind: "Namespace", Name: "namespaces", Shortnames: "ns", Namespaced: false},
				{APIGroup: "apps", Kind: "DaemonSet", Name: "daemonsets", Shortnames: "ds", Namespaced: true},
				{APIGroup: "apps", Kind: "Deployment", Name: "deployments", Shortnames: "deploy", Namespaced: true},
				{APIGroup: "networking.k8s.io", Kind: "Ingress", Name: "ingresses", Shortnames: "ing", Namespaced: true},
			},
			wantFQNs: []string{
				"namespaces",
				"daemonsets.apps",
				"deployments.apps",
				"ingresses.networking.k8s.io",
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
			assert.Equal(t, c.wantFQNs, collectFQNs(c.dest))
		})
	}
}

// This is a snippet that was pulled from output generated with kubectl v1.26.0
// $ kubectl api-resources --cached --output=wide
var tblv126 = `
NAME                              SHORTNAMES   APIVERSION                        NAMESPACED   KIND                             VERBS                                                        CATEGORIES
bindings                                       v1                                true         Binding                          create
componentstatuses                 cs           v1                                false        ComponentStatus                  get,list
configmaps                        cm           v1                                true         ConfigMap                        create,delete,deletecollection,get,list,patch,update,watch
pods                              po           v1                                true         Pod                              create,delete,deletecollection,get,list,patch,update,watch   all
`

// this output was generated with kubectl v1.21.1
// $ kubectl api-resources | grep -e "Deployment\|DaemonSet\|Namespace\|networking.k8s.io.*Ingress$\|KIND"
var tblv121 = `
NAME                                SHORTNAMES                             APIVERSION                             NAMESPACED   KIND
namespaces                          ns                                     v1                                     false        Namespace
daemonsets                          ds                                     apps/v1                                true         DaemonSet
deployments                         deploy                                 apps/v1                                true         Deployment
ingresses                           ing                                    networking.k8s.io/v1                   true         Ingress
`

// this output was generated with kubectl v1.18.10
// $ kubectl api-resources | grep -e "Deployment\|DaemonSet\|Namespace\|networking.k8s.io.*Ingress$\|KIND"
var tblv118 = `
NAME                                SHORTNAMES                             APIGROUP                       NAMESPACED   KIND
namespaces                          ns                                                                    false        Namespace
daemonsets                          ds                                     apps                           true         DaemonSet
deployments                         deploy                                 apps                           true         Deployment
ingresses                           ing                                    networking.k8s.io              true         Ingress
`

var tblEmpty = `
APIVERSION    NAME        NAMESPACED
`

var tblNoHeader = `
apps        Deployment  true
networking  Ingress     true
`

var tblNothing = ``
