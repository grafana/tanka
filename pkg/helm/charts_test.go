package helm

import (
	"path/filepath"
	"testing"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeHelm struct{}

func (f fakeHelm) Pull(chart, version string, opts PullOpts) error {
	return nil
}

func (f fakeHelm) RepoUpdate(opts Opts) error {
	return nil
}

func (f fakeHelm) Template(name, chart string, opts TemplateOpts) (manifest.List, error) {
	return nil, nil
}

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
	c.Helm = &fakeHelm{}

	err = c.Add([]string{"stable/package@1.0.0"})
	assert.NoError(t, err)

	// Adding again the same chart
	err = c.Add([]string{"stable/package@1.0.0"})
	assert.EqualError(t, err, "1 Chart(s) were skipped. Please check above logs for details")

	// Update a version
	err = c.Add([]string{"stable/package@1.1.0"})
	assert.NoError(t, err)
}
