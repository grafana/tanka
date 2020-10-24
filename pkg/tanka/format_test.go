package tanka

import (
	"testing"

	"github.com/gobwas/glob"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatFiles_IgnoresVendor(t *testing.T) {
	var files []string
	_, err := FormatFiles([]string{"./testdata/cases/format/"}, &FormatOpts{
		Excludes: []glob.Glob{glob.MustCompile("**/vendor/**")},
		OutFn: func(f, _ string) error {
			files = append(files, f)
			return nil
		},

		PrintNames: false,
	})

	require.NoError(t, err)
	assert.Contains(t, files, "testdata/cases/format/a.jsonnet")
	assert.Contains(t, files, "testdata/cases/format/b.libsonnet")
	assert.Contains(t, files, "testdata/cases/format/foo/a.jsonnet")
	assert.Contains(t, files, "testdata/cases/format/foo/b.libsonnet")
	assert.NotContains(t, files, "testdata/cases/format/vendor/a.jsonnet")
	assert.NotContains(t, files, "testdata/cases/format/vendor/b.libsonnet")
}
