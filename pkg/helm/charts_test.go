package helm

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddRepos(t *testing.T) {
	c, err := InitChartfile(filepath.Join(t.TempDir(), Filename))
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
	c, err := InitChartfile(filepath.Join(tempDir, Filename))
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
