package helm

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseReq(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected *Requirement
		err      error
	}{
		{
			name:  "valid",
			input: "stable/package@1.0.0",
			expected: &Requirement{
				Chart:     "stable/package",
				Version:   "1.0.0",
				Directory: "",
			},
		},
		{
			name:  "with-path",
			input: "stable/package-name@1.0.0:my-path",
			expected: &Requirement{
				Chart:     "stable/package-name",
				Version:   "1.0.0",
				Directory: "my-path",
			},
		},
		{
			name:  "with-path-with-special-chars",
			input: "stable/package@v1.24.0:my weird-path_test",
			expected: &Requirement{
				Chart:     "stable/package",
				Version:   "v1.24.0",
				Directory: "my weird-path_test",
			},
		},
		{
			name:  "with-path-with-sub-path",
			input: "stable/package@3.45.6:" + filepath.Join("myparentdir1", "mysubdir1", "mypath"),
			expected: &Requirement{
				Chart:     "stable/package",
				Version:   "3.45.6",
				Directory: filepath.Join("myparentdir1", "mysubdir1", "mypath"),
			},
		},
		{
			name:  "url-instead-of-repo",
			input: "https://helm.releases.hashicorp.com/vault@0.19.0",
			err:   errors.New("not of form 'repo/chart@version(:path)' where repo contains no special characters"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := parseReq(tc.input)
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.expected, req)
		})
	}
}

func TestAddRepos(t *testing.T) {
	c, err := InitChartfile(filepath.Join(t.TempDir(), Filename))
	require.NoError(t, err)

	err = c.AddRepos(
		Repo{Name: "foo", URL: "https://foo.com"},
		Repo{Name: "foo2", URL: "https://foo2.com"},
	)
	assert.NoError(t, err)

	// Only \w characters are allowed in repo names
	err = c.AddRepos(
		Repo{Name: "with-dashes", URL: "https://foo.com"},
	)
	assert.EqualError(t, err, "1 Repo(s) were skipped. Please check above logs for details")

	err = c.AddRepos(
		Repo{Name: "foo", URL: "https://foo.com"},
	)
	assert.EqualError(t, err, "1 Repo(s) were skipped. Please check above logs for details")
}

func TestAdd(t *testing.T) {
	tempDir := t.TempDir()
	c, err := InitChartfile(filepath.Join(tempDir, Filename))
	require.NoError(t, err)

	err = c.Add([]string{"stable/prometheus@11.12.1"})
	assert.NoError(t, err)

	// Adding again the same chart
	err = c.Add([]string{"stable/prometheus@11.12.1"})
	assert.EqualError(t, err, "1 Chart(s) were skipped. Please check above logs for details")

	// Adding a chart with a different version to the same path, causes a conflict
	err = c.Add([]string{"stable/prometheus@11.12.0"})
	assert.EqualError(t, err, `Validation errors:
 - output directory "prometheus" is used twice, by charts "stable/prometheus@11.12.1" and "stable/prometheus@11.12.0"`)

	// Add a chart with a specific extract directory
	err = c.Add([]string{"stable/prometheus@11.12.0:prometheus-11.12.0"})
	assert.NoError(t, err)

	// Add a chart with a nested extract directory
	err = c.Add([]string{"stable/prometheus@11.12.0:" + filepath.Join("zparentdir", "prometheus-11.12.0")})
	assert.NoError(t, err)

	// Check file contents
	listResult, err := os.ReadDir(filepath.Join(tempDir, "charts"))
	assert.NoError(t, err)
	assert.Equal(t, 3, len(listResult))
	assert.Equal(t, "prometheus", listResult[0].Name())
	assert.Equal(t, "prometheus-11.12.0", listResult[1].Name())
	assert.Equal(t, "zparentdir", listResult[2].Name())
	listResult, err = os.ReadDir(filepath.Join(tempDir, "charts", "zparentdir"))
	assert.NoError(t, err)
	assert.Equal(t, 1, len(listResult))
	assert.Equal(t, "prometheus-11.12.0", listResult[0].Name())

	chartContent, err := os.ReadFile(filepath.Join(tempDir, "charts", "prometheus", "Chart.yaml"))
	assert.NoError(t, err)
	assert.Contains(t, string(chartContent), `version: 11.12.1`)

	chartContent, err = os.ReadFile(filepath.Join(tempDir, "charts", "prometheus-11.12.0", "Chart.yaml"))
	assert.NoError(t, err)
	assert.Contains(t, string(chartContent), `version: 11.12.0`)

	chartContent, err = os.ReadFile(filepath.Join(tempDir, "charts", "zparentdir", "prometheus-11.12.0", "Chart.yaml"))
	assert.NoError(t, err)
	assert.Contains(t, string(chartContent), `version: 11.12.0`)
}

func TestAddOCI(t *testing.T) {
	tempDir := t.TempDir()
	c, err := InitChartfile(filepath.Join(tempDir, Filename))
	require.NoError(t, err)

	err = c.AddRepos(Repo{Name: "karpenter", URL: "oci://public.ecr.aws/karpenter"})
	assert.NoError(t, err)

	err = c.Add([]string{"karpenter/karpenter@v0.27.1"})
	assert.NoError(t, err)

	// Check file contents
	listResult, err := os.ReadDir(filepath.Join(tempDir, "charts"))
	assert.NoError(t, err)
	assert.Equal(t, 1, len(listResult))
	assert.Equal(t, "karpenter", listResult[0].Name())
}

func TestRevendorDeletedFiles(t *testing.T) {
	tempDir := t.TempDir()
	c, err := InitChartfile(filepath.Join(tempDir, Filename))
	require.NoError(t, err)

	err = c.Add([]string{"stable/prometheus@11.12.1"})
	assert.NoError(t, err)

	// Check file contents
	chartContent, err := os.ReadFile(filepath.Join(tempDir, "charts", "prometheus", "Chart.yaml"))
	assert.NoError(t, err)
	assert.Contains(t, string(chartContent), `version: 11.12.1`)

	// Delete the whole dir and revendor
	require.NoError(t, os.RemoveAll(filepath.Join(tempDir, "charts", "prometheus")))
	assert.NoError(t, c.Vendor(true))

	// Check file contents
	chartContent, err = os.ReadFile(filepath.Join(tempDir, "charts", "prometheus", "Chart.yaml"))
	assert.NoError(t, err)
	assert.Contains(t, string(chartContent), `version: 11.12.1`)

	// Delete just the Chart.yaml and revendor
	require.NoError(t, os.Remove(filepath.Join(tempDir, "charts", "prometheus", "Chart.yaml")))
	assert.NoError(t, c.Vendor(true))

	// Check file contents
	chartContent, err = os.ReadFile(filepath.Join(tempDir, "charts", "prometheus", "Chart.yaml"))
	assert.NoError(t, err)
	assert.Contains(t, string(chartContent), `version: 11.12.1`)
}

func TestPrune(t *testing.T) {
	for _, prune := range []bool{false, true} {
		t.Run(fmt.Sprintf("%t", prune), func(t *testing.T) {
			tempDir := t.TempDir()
			c, err := InitChartfile(filepath.Join(tempDir, Filename))
			require.NoError(t, err)

			// Add a chart
			require.NoError(t, c.Add([]string{"stable/prometheus@11.12.1"}))

			// Add a chart with a directory
			require.NoError(t, c.Add([]string{"stable/prometheus@11.12.1:custom-dir"}))

			// Add unrelated files and folders
			require.NoError(t, os.WriteFile(filepath.Join(tempDir, "charts", "foo.txt"), []byte("foo"), 0644))
			require.NoError(t, os.Mkdir(filepath.Join(tempDir, "charts", "foo"), 0755))
			require.NoError(t, os.WriteFile(filepath.Join(tempDir, "charts", "foo", "Chart.yaml"), []byte("foo"), 0644))

			require.NoError(t, c.Vendor(prune))

			// Check if files are pruned
			listResult, err := os.ReadDir(filepath.Join(tempDir, "charts"))
			assert.NoError(t, err)
			if prune {
				assert.Equal(t, 2, len(listResult))
				assert.Equal(t, "custom-dir", listResult[0].Name())
				assert.Equal(t, "prometheus", listResult[1].Name())
			} else {
				assert.Equal(t, 4, len(listResult))
				chartContent, err := os.ReadFile(filepath.Join(tempDir, "charts", "foo", "Chart.yaml"))
				assert.NoError(t, err)
				assert.Contains(t, string(chartContent), `foo`)
			}
		})
	}
}

func TestInvalidChartName(t *testing.T) {
	tempDir := t.TempDir()
	c, err := InitChartfile(filepath.Join(tempDir, Filename))
	require.NoError(t, err)

	c.Manifest.Requires = append(c.Manifest.Requires, Requirement{
		Chart:   "noslash",
		Version: "1.0.0",
	})

	err = c.Vendor(false)
	assert.EqualError(t, err, `Validation errors:
 - Chart name "noslash" is not valid. Expecting a repo/name format.`)
}
