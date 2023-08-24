package tanka

import (
	"strings"
	"testing"

	"github.com/grafana/tanka/pkg/jsonnet"
	"github.com/grafana/tanka/pkg/jsonnet/implementations/goimpl"
	"github.com/stretchr/testify/assert"
)

var jsonnetImpl = &goimpl.JsonnetGoImplementation{}

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
	json, err := evalJsonnet("testdata/cases/withtlas", jsonnetImpl, opts)
	assert.NoError(t, err)
	assert.Equal(t, `"foovalue"`, strings.TrimSpace(json))
}

func TestEvalJsonnetWithExpression(t *testing.T) {
	exprs := []string{`["testCase"]`, "testCase"}

	for _, expr := range exprs {
		t.Run(expr, func(t *testing.T) {
			opts := jsonnet.Opts{
				EvalScript: PatternEvalScript(expr),
			}

			// This will fail intermittently if TLAs are passed as positional
			// parameters.
			json, err := evalJsonnet("testdata/cases/object", jsonnetImpl, opts)
			assert.NoError(t, err)
			assert.Equal(t, `"object"`, strings.TrimSpace(json))
		})
	}
}
