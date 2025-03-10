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
			name:  "url-instead-of-repo",
			input: "https://helm.releases.hashicorp.com/vault@0.19.0",
			err:   errors.New("not of form 'repo/chart@version(:path)' where repo contains no special characters"),
		},
		{
			name:  "repo-with-special-chars",
			input: "with-dashes/package@1.0.0",
			expected: &Requirement{
				Chart:     "with-dashes/package",
				Version:   "1.0.0",
				Directory: "",
			},
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
		Repo{Name: "with-dashes", URL: "https://foo.com"},
	)
	assert.NoError(t, err)

	// Only \w- characters are allowed in repo names
	err = c.AddRepos(
		Repo{Name: "re:po", URL: "https://foo.com"},
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

	err = c.Add([]string{"stable/prometheus@11.12.1"}, "")
	assert.NoError(t, err)

	// Adding again the same chart
	err = c.Add([]string{"stable/prometheus@11.12.1"}, "")
	assert.EqualError(t, err, "1 Chart(s) were skipped. Please check above logs for details")

	// Adding a chart with a different version to the same path, causes a conflict
	err = c.Add([]string{"stable/prometheus@11.12.0"}, "")
	assert.EqualError(t, err, `validation errors:
 - output directory "prometheus" is used twice, by charts "stable/prometheus@11.12.1" and "stable/prometheus@11.12.0"`)

	// Add a chart with a specific extract directory
	err = c.Add([]string{"stable/prometheus@11.12.0:prometheus-11.12.0"}, "")
	assert.NoError(t, err)

	// Add a chart while specifying a helm repo config file
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "helmConfig.yaml"), []byte(`
apiVersion: ""
generated: "0001-01-01T00:00:00Z"
repositories:
- caFile: ""
  certFile: ""
  insecure_skip_tls_verify: false
  keyFile: ""
  name: private
  pass_credentials_all: false
  password: ""
  url: https://charts.helm.sh/stable
  username: ""
`), 0644))
	err = c.Add([]string{"private/prometheus@11.12.1:private-11.12.1"}, filepath.Join(tempDir, "helmConfig.yaml"))
	assert.NoError(t, err)

	// Check file contents
	listResult, err := os.ReadDir(filepath.Join(tempDir, "charts"))
	assert.NoError(t, err)
	assert.Equal(t, 3, len(listResult))
	assert.Equal(t, "private-11.12.1", listResult[0].Name())
	assert.Equal(t, "prometheus", listResult[1].Name())
	assert.Equal(t, "prometheus-11.12.0", listResult[2].Name())

	chartContent, err := os.ReadFile(filepath.Join(tempDir, "charts", "prometheus", "Chart.yaml"))
	assert.NoError(t, err)
	assert.Contains(t, string(chartContent), `version: 11.12.1`)

	chartContent, err = os.ReadFile(filepath.Join(tempDir, "charts", "prometheus-11.12.0", "Chart.yaml"))
	assert.NoError(t, err)
	assert.Contains(t, string(chartContent), `version: 11.12.0`)

	chartContent, err = os.ReadFile(filepath.Join(tempDir, "charts", "private-11.12.1", "Chart.yaml"))
	assert.NoError(t, err)
	assert.Contains(t, string(chartContent), `version: 11.12.1`)
}

func TestAddOCI(t *testing.T) {
	tempDir := t.TempDir()
	c, err := InitChartfile(filepath.Join(tempDir, Filename))
	require.NoError(t, err)

	err = c.AddRepos(Repo{Name: "karpenter", URL: "oci://public.ecr.aws/karpenter"})
	assert.NoError(t, err)

	err = c.Add([]string{"karpenter/karpenter@v0.27.1"}, "")
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

	err = c.Add([]string{"stable/prometheus@11.12.1"}, "")
	assert.NoError(t, err)

	// Check file contents
	chartContent, err := os.ReadFile(filepath.Join(tempDir, "charts", "prometheus", "Chart.yaml"))
	assert.NoError(t, err)
	assert.Contains(t, string(chartContent), `version: 11.12.1`)

	// Delete the whole dir and revendor
	require.NoError(t, os.RemoveAll(filepath.Join(tempDir, "charts", "prometheus")))
	assert.NoError(t, c.Vendor(true, ""))

	// Check file contents
	chartContent, err = os.ReadFile(filepath.Join(tempDir, "charts", "prometheus", "Chart.yaml"))
	assert.NoError(t, err)
	assert.Contains(t, string(chartContent), `version: 11.12.1`)

	// Delete just the Chart.yaml and revendor
	require.NoError(t, os.Remove(filepath.Join(tempDir, "charts", "prometheus", "Chart.yaml")))
	assert.NoError(t, c.Vendor(true, ""))

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
			require.NoError(t, c.Add([]string{"stable/prometheus@11.12.1"}, ""))

			// Add a chart with a directory
			require.NoError(t, c.Add([]string{"stable/prometheus@11.12.1:custom-dir"}, ""))

			// Add unrelated files and folders
			require.NoError(t, os.WriteFile(filepath.Join(tempDir, "charts", "foo.txt"), []byte("foo"), 0644))
			require.NoError(t, os.Mkdir(filepath.Join(tempDir, "charts", "foo"), 0755))
			require.NoError(t, os.WriteFile(filepath.Join(tempDir, "charts", "foo", "Chart.yaml"), []byte("foo"), 0644))

			require.NoError(t, c.Vendor(prune, ""))

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

	err = c.Vendor(false, "")
	assert.EqualError(t, err, `validation errors:
 - Chart name "noslash" is not valid. Expecting a repo/name format.`)
}

func TestConfigFileOption(t *testing.T) {
	tempDir := t.TempDir()
	c, err := InitChartfile(filepath.Join(tempDir, Filename))
	require.NoError(t, err)

	// Don't want to commit credentials so we just verify the "private" repo reference will make
	// use of this helm config since the InitChartfile does not have a reference to it.
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "helmConfig.yaml"), []byte(`
apiVersion: ""
generated: "0001-01-01T00:00:00Z"
repositories:
- caFile: ""
  certFile: ""
  insecure_skip_tls_verify: false
  keyFile: ""
  name: private
  pass_credentials_all: false
  password: ""
  url: https://charts.helm.sh/stable
  username: ""
`), 0644))
	c.Manifest.Requires = append(c.Manifest.Requires, Requirement{
		Chart:   "private/prometheus",
		Version: "11.12.1",
	})

	err = c.Vendor(false, filepath.Join(tempDir, "helmConfig.yaml"))
	assert.NoError(t, err)

	chartContent, err := os.ReadFile(filepath.Join(tempDir, "charts", "prometheus", "Chart.yaml"))
	assert.NoError(t, err)
	assert.Contains(t, string(chartContent), `version: 11.12.1`)
}

func TestChartsVersionCheck(t *testing.T) {
	tempDir := t.TempDir()
	c, err := InitChartfile(filepath.Join(tempDir, Filename))
	require.NoError(t, err)

	err = c.Add([]string{"stable/prometheus@11.12.0"}, "")
	assert.NoError(t, err)

	// Having multiple versions of the same chart should only return one update
	err = c.Add([]string{"stable/prometheus@11.11.0:old"}, "")
	assert.NoError(t, err)

	chartVersions, err := c.VersionCheck("")
	assert.NoError(t, err)

	// stable/prometheus is deprecated so only the 11.12.1 should ever be returned as latest
	latestPrometheusChartVersion := ChartSearchVersion{
		Name:        "stable/prometheus",
		Version:     "11.12.1",
		AppVersion:  "2.20.1",
		Description: "DEPRECATED Prometheus is a monitoring system and time series database.",
	}
	stableExpected := RequiresVersionInfo{
		Name:                       "stable/prometheus",
		Directory:                  "",
		CurrentVersion:             "11.12.0",
		UsingLatestVersion:         false,
		LatestVersion:              latestPrometheusChartVersion,
		LatestMatchingMajorVersion: latestPrometheusChartVersion,
		LatestMatchingMinorVersion: latestPrometheusChartVersion,
	}
	oldExpected := RequiresVersionInfo{
		Name:                       "stable/prometheus",
		Directory:                  "old",
		CurrentVersion:             "11.11.0",
		UsingLatestVersion:         false,
		LatestVersion:              latestPrometheusChartVersion,
		LatestMatchingMajorVersion: latestPrometheusChartVersion,
		LatestMatchingMinorVersion: ChartSearchVersion{
			Name:        "stable/prometheus",
			Version:     "11.11.1",
			AppVersion:  "2.19.0",
			Description: "Prometheus is a monitoring system and time series database.",
		},
	}
	assert.Equal(t, 2, len(chartVersions))
	assert.Equal(t, stableExpected, chartVersions["stable/prometheus@11.12.0"])
	assert.Equal(t, oldExpected, chartVersions["stable/prometheus@11.11.0"])
}

func TestVersionCheckWithConfig(t *testing.T) {
	tempDir := t.TempDir()
	c, err := InitChartfile(filepath.Join(tempDir, Filename))
	require.NoError(t, err)

	// Don't want to commit credentials so we just verify the "private" repo reference will make
	// use of this helm config since the InitChartfile does not have a reference to it.
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "helmConfig.yaml"), []byte(`
apiVersion: ""
generated: "0001-01-01T00:00:00Z"
repositories:
- caFile: ""
  certFile: ""
  insecure_skip_tls_verify: false
  keyFile: ""
  name: private
  pass_credentials_all: false
  password: ""
  url: https://charts.helm.sh/stable
  username: ""
`), 0644))
	c.Manifest.Requires = append(c.Manifest.Requires, Requirement{
		Chart:   "private/prometheus",
		Version: "11.12.0",
	})

	chartVersions, err := c.VersionCheck(filepath.Join(tempDir, "helmConfig.yaml"))
	assert.NoError(t, err)

	// stable/prometheus is deprecated so only the 11.12.1 should ever be returned as latest
	latestPrometheusChartVersion := ChartSearchVersion{
		Name:        "private/prometheus",
		Version:     "11.12.1",
		AppVersion:  "2.20.1",
		Description: "DEPRECATED Prometheus is a monitoring system and time series database.",
	}
	expected := RequiresVersionInfo{
		Name:                       "private/prometheus",
		Directory:                  "",
		CurrentVersion:             "11.12.0",
		UsingLatestVersion:         false,
		LatestVersion:              latestPrometheusChartVersion,
		LatestMatchingMajorVersion: latestPrometheusChartVersion,
		LatestMatchingMinorVersion: latestPrometheusChartVersion,
	}
	assert.Equal(t, 1, len(chartVersions))
	assert.Equal(t, expected, chartVersions["private/prometheus@11.12.0"])
}
