package tanka

import (
	"strings"
	"testing"

	"github.com/grafana/tanka/pkg/jsonnet"
	"github.com/stretchr/testify/assert"
)

func TestEvalJsonnet(t *testing.T) {
	var tlaCode jsonnet.InjectedCode
	// Pass in the mandatory parameters as TLA codes, but note that only `foo`
	// contains `data`, which is a valid key inside the `o` object defined in
	// testdata/cases/withtlas/main.jsonnet. If they are passed as positional
	// parameters, then their names are ignored, which will lead to arbitrary
	// failures because the order in which they're passed is random.
	// `evalJsonnet` has been updated to pass them as named parameters.
	tlaCode.Set("foo", "'data'")
	tlaCode.Set("bar", "'kaboom'")
	tlaCode.Set("baz", "'kaboom'")

	opts := jsonnet.Opts{
		EvalScript: "main",
		TLACode:    tlaCode,
	}

	// This will fail intermittently if TLAs are passed as positional
	// parameters.
	json, err := evalJsonnet("testdata/cases/withtlas", opts)
	assert.NoError(t, err)
	assert.Equal(t, `"foovalue"`, strings.TrimSpace(json))
}
