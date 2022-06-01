package helm

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddRepos(t *testing.T) {
	c, err := InitChartfile(filepath.Join(t.TempDir(), Filename), false)
	require.NoError(t, err)

	err = c.AddRepos(
		Repo{Name: "foo", URL: "https://foo.com"},
		Repo{Name: "foo2", URL: "https://foo2.com"},
	)
	assert.NoError(t, err)

	err = c.AddRepos(
		Repo{Name: "foo", URL: "https://foo.com"},
	)
	assert.EqualError(t, err, "1 Repo(s) were skipped. Please check above logs for details")
}

func TestAdd(t *testing.T) {
	tempDir := t.TempDir()
	c, err := InitChartfile(filepath.Join(tempDir, Filename), false)
	require.NoError(t, err)

	err = c.Add([]string{"stable/prometheus@11.12.1"})
	assert.NoError(t, err)

	// Adding again the same chart
	err = c.Add([]string{"stable/prometheus@11.12.1"})
	assert.EqualError(t, err, "1 Chart(s) were skipped. Please check above logs for details")

	// Add a chart with a specific extract directory
	err = c.Add([]string{"stable/prometheus@11.12.0:prometheus-11.12.0"})
	assert.NoError(t, err)

	// Check file contents
	listResult, err := os.ReadDir(filepath.Join(tempDir, "charts"))
	assert.NoError(t, err)
	assert.Equal(t, 2, len(listResult))
	assert.Equal(t, "prometheus", listResult[0].Name())
	assert.Equal(t, "prometheus-11.12.0", listResult[1].Name())

	chartContent, err := os.ReadFile(filepath.Join(tempDir, "charts", "prometheus", "Chart.yaml"))
	assert.NoError(t, err)
	assert.Contains(t, string(chartContent), `version: 11.12.1`)

	chartContent, err = os.ReadFile(filepath.Join(tempDir, "charts", "prometheus-11.12.0", "Chart.yaml"))
	assert.NoError(t, err)
	assert.Contains(t, string(chartContent), `version: 11.12.0`)
}

func TestPrune(t *testing.T) {
	for _, prune := range []bool{false, true} {
		t.Run(fmt.Sprintf("%t", prune), func(t *testing.T) {
			tempDir := t.TempDir()
			c, err := InitChartfile(filepath.Join(tempDir, Filename), prune)
			require.NoError(t, err)

			// Add a chart
			err = c.Add([]string{"stable/prometheus@11.12.1"})
			require.NoError(t, err)

			// Add unrelated files and folders
			err = os.WriteFile(filepath.Join(tempDir, "charts", "foo.txt"), []byte("foo"), 0644)
			err = os.Mkdir(filepath.Join(tempDir, "charts", "foo"), 0755)
			require.NoError(t, err)
			err = os.WriteFile(filepath.Join(tempDir, "charts", "foo", "Chart.yaml"), []byte("foo"), 0644)
			require.NoError(t, err)

			err = c.Vendor()
			require.NoError(t, err)

			// Check if files are pruned
			listResult, err := os.ReadDir(filepath.Join(tempDir, "charts"))
			assert.NoError(t, err)
			if prune {
				assert.Equal(t, 1, len(listResult))
			} else {
				assert.Equal(t, 3, len(listResult))
				chartContent, err := os.ReadFile(filepath.Join(tempDir, "charts", "foo", "Chart.yaml"))
				assert.NoError(t, err)
				assert.Contains(t, string(chartContent), `foo`)
			}
		})
	}
}
