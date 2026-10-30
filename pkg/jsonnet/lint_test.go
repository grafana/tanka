package jsonnet

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLint(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		opts := &LintOpts{Parallelism: 4}
		err := Lint([]string{"testdata/importTree"}, opts)
		assert.NoError(t, err)
	})

	t.Run("environment-relative import from a nested file", func(t *testing.T) {
		buf := &bytes.Buffer{}
		opts := &LintOpts{Out: buf, Parallelism: 4}
		err := Lint([]string{"testdata/lintEnvRelativeImport/environments/default"}, opts)
		assert.NoError(t, err)
		assert.Empty(t, buf.String())
	})

	t.Run("custom entrypoint with an environment vendor dir", func(t *testing.T) {
		buf := &bytes.Buffer{}
		opts := &LintOpts{Out: buf, Parallelism: 4}
		err := Lint([]string{"testdata/lintCustomEntrypoint/environments/custom/custom.jsonnet"}, opts)
		assert.NoError(t, err)
		assert.Empty(t, buf.String())
	})

	t.Run("error", func(t *testing.T) {
		buf := &bytes.Buffer{}
		opts := &LintOpts{Out: buf, Parallelism: 4}
		err := Lint([]string{"testdata/lintingError"}, opts)
		assert.EqualError(t, err, "Linting has failed for at least one file")
		assert.Equal(t, absPath(t, "testdata/lintingError/main.jsonnet")+`:1:7-22 Unused variable: unused

local unused = 'test';


`, buf.String())
	})
}
