package helm

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
)

const calledFrom = "/my/path/here"

type MockHelm struct {
	mock.Mock
}

// fulfill the Helm interface
func (m *MockHelm) Pull(chart, version string, opts PullOpts) error {
	args := m.Called(chart, version, opts)
	return args.Error(0)
}

func (m *MockHelm) RepoUpdate(opts Opts) error {
	args := m.Called(opts)
	return args.Error(0)
}

func (m *MockHelm) Template(name, chart string, opts TemplateOpts) (manifest.List, error) {
	args := m.Called(name, chart, opts)

	// figure out what arguments `helm template` would be called with and save
	// them
	execHelm := &ExecHelm{}
	cmdArgs := execHelm.templateCommandArgs(name, chart, opts)
	m.TestData().Set("templateCommandArgs", cmdArgs)

	return args.Get(0).(manifest.List), args.Error(1)
}

func (m *MockHelm) ChartExists(chart string, opts *JsonnetOpts) (string, error) {
	args := m.Called(chart, opts)
	return args.String(0), args.Error(1)
}

func callNativeFunction(t *testing.T, templateOpts TemplateOpts, parameters []interface{}) []string {
	t.Helper()

	helmMock := &MockHelm{}

	helmMock.On(
		"ChartExists",
		"chart",
		mock.AnythingOfType("*helm.JsonnetOpts")).
		Return("/full/chart/path", nil).
		Once()

	// this verifies that the helmMock.Template() method is called with the
	// correct arguments, i.e. includeCrds: true is set by default
	helmMock.On("Template", "name", "/full/chart/path", templateOpts).
		Return(manifest.List{}, nil).
		Once()

	nf := NativeFunc(helmMock)
	require.NotNil(t, nf)

	params := []string{
		"name",
		"chart",
	}

	opts := make(map[string]interface{})
	opts["calledFrom"] = calledFrom

	// params + opts
	paramsInterface := make([]interface{}, 3)
	paramsInterface[0] = params[0]
	paramsInterface[1] = params[1]
	paramsInterface[2] = opts

	_, err := nf.Func(parameters)

	require.NoError(t, err)

	helmMock.AssertExpectations(t)

	return helmMock.TestData().Get("templateCommandArgs").StringSlice()
}

// TestDefaultCommandineFlagsIncludeCrds tests that the includeCrds flag is set
// to true by default
func TestDefaultCommandineFlagsIncludeCrds(t *testing.T) {
	kubeVersion := "1.18.0"

	// we will check that the template function is called with these options,
	// i.e. that includeCrds got set to true. This is not us passing an input,
	// we are asserting here that the template function is called with these
	// options.
	templateOpts := TemplateOpts{
		KubeVersion: kubeVersion,
		IncludeCRDs: true,
	}

	params := []string{
		"name",
		"chart",
	}

	// we do not set includeCrds here, so it should be true by default
	opts := make(map[string]interface{})
	opts["calledFrom"] = calledFrom
	opts["kubeVersion"] = kubeVersion

	// params + opts
	paramsInterface := make([]interface{}, 3)
	paramsInterface[0] = params[0]
	paramsInterface[1] = params[1]
	paramsInterface[2] = opts

	args := callNativeFunction(t, templateOpts, paramsInterface)

	// finally check that the actual command line arguments we will pass to
	// `helm template` contain the --include-crds flag
	require.Contains(t, args, "--include-crds")
}

// TestIncludeCrdsFalse tests that the includeCrds flag is can be set to false,
// and this makes it to the helm.Template() method call
func TestIncludeCrdsFalse(t *testing.T) {
	kubeVersion := "1.18.0"

	// we will check that the template function is called with these options,
	// i.e. that includeCrds got set to false. This is not us passing an input,
	// we are asserting here that the template function is called with these
	// options.
	templateOpts := TemplateOpts{
		KubeVersion: kubeVersion,
		IncludeCRDs: false,
	}

	params := []string{
		"name",
		"chart",
	}

	// we explicitly set includeCrds to false here
	opts := make(map[string]interface{})
	opts["calledFrom"] = calledFrom
	opts["kubeVersion"] = kubeVersion
	opts["includeCrds"] = false

	// params + opts
	paramsInterface := make([]interface{}, 3)
	paramsInterface[0] = params[0]
	paramsInterface[1] = params[1]
	paramsInterface[2] = opts

	args := callNativeFunction(t, templateOpts, paramsInterface)

	// finally check that the actual command line arguments we will pass to
	// `helm template` don't contain the --include-crds flag
	require.NotContains(t, args, "--include-crds")
}
