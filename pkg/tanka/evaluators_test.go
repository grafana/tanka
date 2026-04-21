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
	json, err := evalJsonnet(t.Context(), "testdata/cases/withtlas", jsonnetImpl, opts)
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
			json, err := evalJsonnet(t.Context(), "testdata/cases/object", jsonnetImpl, opts)
			assert.NoError(t, err)
			assert.Equal(t, `"object"`, strings.TrimSpace(json))
		})
	}
}

// An EvalScript with a top-level function containing only optional arguments
// should be evaluated as a function even if no TLAs are provided.
func TestEvalWithOptionalTlas(t *testing.T) {
	opts := jsonnet.Opts{
		EvalScript: "main.metadata.name",
	}
	json, err := evalJsonnet(t.Context(), "testdata/cases/with-optional-tlas/main.jsonnet", jsonnetImpl, opts)
	assert.NoError(t, err)
	assert.Equal(t, `"bar-baz"`, strings.TrimSpace(json))
}

// An EvalScript with a top-level function containing should allow passing only
// a subset of the TLAs.
func TestEvalWithOptionalTlasSpecifiedArg2(t *testing.T) {
	opts := jsonnet.Opts{
		EvalScript: "main.metadata.name",
		TLACode:    jsonnet.InjectedCode{"baz": "'changed'"},
	}
	json, err := evalJsonnet(t.Context(), "testdata/cases/with-optional-tlas/main.jsonnet", jsonnetImpl, opts)
	assert.NoError(t, err)
	assert.Equal(t, `"bar-changed"`, strings.TrimSpace(json))
}

// An EvalScript with a top-level function having no arguments should be
// evaluated as a function even if no TLAs are provided.
func TestEvalFunctionWithNoTlas(t *testing.T) {
	opts := jsonnet.Opts{
		EvalScript: "main.metadata.name",
	}
	json, err := evalJsonnet(t.Context(), "testdata/cases/function-with-zero-params/main.jsonnet", jsonnetImpl, opts)
	assert.NoError(t, err)
	assert.Equal(t, `"inline"`, strings.TrimSpace(json))
}

// An EvalScript with a top-level function should return an understandable
// error message if an incorrect TLA is provided.
func TestInvalidTlaArg(t *testing.T) {
	opts := jsonnet.Opts{
		EvalScript: "main",
		TLACode:    jsonnet.InjectedCode{"foo": "'bar'"},
	}
	json, err := evalJsonnet(t.Context(), "testdata/cases/function-with-zero-params/main.jsonnet", jsonnetImpl, opts)
	assert.Contains(t, err.Error(), "function has no parameter foo")
	assert.Equal(t, "", json)
}

// Providing a TLA to an EvalScript with a non-function top level mainfile
// should not return an error.
func TestTlaWithNonFunction(t *testing.T) {
	opts := jsonnet.Opts{
		EvalScript: "main",
		TLACode:    jsonnet.InjectedCode{"foo": "'bar'"},
	}
	json, err := evalJsonnet(t.Context(), "testdata/cases/withenv/main.jsonnet", jsonnetImpl, opts)
	assert.NoError(t, err)
	assert.NotEmpty(t, json)
}

func TestBuildEvalScript_WindowsPath(t *testing.T) {
	got := buildEvalScript(`C:\Users\foo\bar\main.jsonnet`, "main", nil, false)

	assert.Contains(t, got, `import 'C:/Users/foo/bar/main.jsonnet'`,
		"entrypoint backslashes must be converted to forward slashes")
}

func TestBuildEvalScript_Function(t *testing.T) {
	tlas := []string{"foo", "bar"}
	got := buildEvalScript(`C:\proj\main.jsonnet`, "main", tlas, true)

	assert.Contains(t, got, `function(foo, bar)`)
	assert.Contains(t, got, `import 'C:/proj/main.jsonnet'`)
	assert.Contains(t, got, `(foo=foo, bar=bar)`)
}

func TestBuildEvalScript_UnixPathUnchanged(t *testing.T) {
	got := buildEvalScript(`/home/user/proj/main.jsonnet`, "main", nil, false)
	assert.Contains(t, got, `import '/home/user/proj/main.jsonnet'`)
}
