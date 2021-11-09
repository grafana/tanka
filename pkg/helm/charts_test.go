package helm

import (
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
	c, err := InitChartfile(filepath.Join(t.TempDir(), Filename))
	require.NoError(t, err)

	err = c.Add([]string{
		"stable/foo@1.0.0",
		"stable/foo2@1.0.0",
	}, false)
	assert.NoError(t, err)

	// Adding again the same chart
	err = c.Add([]string{
		"stable/foo@1.0.0",
	}, false)
	assert.EqualError(t, err, "1 Chart(s) were skipped. Please check above logs for details")

	// Update a version
	err = c.Add([]string{
		"stable/foo@1.1.0",
	}, false)
	assert.NoError(t, err)
}
