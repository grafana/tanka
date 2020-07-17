package tanka

import (
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEvalJsonnet(t *testing.T) {
	m := make(map[string]string)
	_,e:=evalJsonnet("./testdata/",v1alpha1.New(),m)
	assert.NoError(t,e)
}
