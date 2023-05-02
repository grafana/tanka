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

func callNativeFunction(t *testing.T, expectedHelmTemplateOptions TemplateOpts, inputOptionsFromJsonnet map[string]interface{}) []string {
	t.Helper()

	helmMock := &MockHelm{}

	helmMock.On(
		"ChartExists",
		"exampleChartPath",
		mock.AnythingOfType("*helm.JsonnetOpts")).
		Return("/full/chart/path", nil).
		Once()

	// this verifies that the helmMock.Template() method is called with the
	// correct arguments, i.e. includeCrds: true is set by default
	helmMock.On("Template", "exampleChartName", "/full/chart/path", expectedHelmTemplateOptions).
		Return(manifest.List{}, nil).
		Once()

	nf := NativeFunc(helmMock)
	require.NotNil(t, nf)

	// the mandatory parameters to helm.template() in Jsonnet
	params := []string{
		"exampleChartName",
		"exampleChartPath",
	}

	// mandatory parameters + the k-v pairs from the Jsonnet input
	paramsInterface := make([]interface{}, 3)
	paramsInterface[0] = params[0]
	paramsInterface[1] = params[1]
	paramsInterface[2] = inputOptionsFromJsonnet

	_, err := nf.Func(paramsInterface)

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
	expectedHelmTemplateOptions := TemplateOpts{
		KubeVersion: kubeVersion,
		IncludeCRDs: true,
	}

	// the options to helm.template(), which are a Jsonnet object, turned into a
	// go map[string]interface{}. we do not set includeCrds here, so it should
	// be true by default
	inputOptionsFromJsonnet := make(map[string]interface{})
	inputOptionsFromJsonnet["calledFrom"] = calledFrom
	inputOptionsFromJsonnet["kubeVersion"] = kubeVersion

	args := callNativeFunction(t, expectedHelmTemplateOptions, inputOptionsFromJsonnet)

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
	expectedHelmTemplateOptions := TemplateOpts{
		KubeVersion: kubeVersion,
		IncludeCRDs: false,
	}

	// the options to helm.template(), which are a Jsonnet object, turned into a
	// go map[string]interface{}. we explicitly set includeCrds to false here
	inputOptionsFromJsonnet := make(map[string]interface{})
	inputOptionsFromJsonnet["calledFrom"] = calledFrom
	inputOptionsFromJsonnet["kubeVersion"] = kubeVersion
	inputOptionsFromJsonnet["includeCrds"] = false

	args := callNativeFunction(t, expectedHelmTemplateOptions, inputOptionsFromJsonnet)

	// finally check that the actual command line arguments we will pass to
	// `helm template` don't contain the --include-crds flag
	require.NotContains(t, args, "--include-crds")
}
