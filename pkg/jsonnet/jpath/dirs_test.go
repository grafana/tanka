package jpath

import (
	"os"
	"path/filepath"
	"testing"

	"io/ioutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type scenario struct {
	name        string
	testdata    []string
	environment string

	// expected results
	baseDir string
	rootDir string
	err     error
}

// TestFindRoot asserts that baseDir and rootDir can be correctly resolved.
//
// To do so, Tanka searches the directory tree from the passed directory up twice:
// - for main.jsonnet to find the baseDir
// - for jsonnetfile.json (or tkrc.yaml) to find the rootDir
//
// This enables a git-like behaviour (regardless how deep you are in the
// project, it works)
func TestFindRoot(t *testing.T) {
	scenarios := []scenario{
		// Scenario: Missing base pointerfile. We expect an ErrorNoBase.
		{
			name:        "missing-basePointer",
			environment: "environments/default",
			testdata:    []string{"jsonnetfile.json", "environments/default/"},
			err:         ErrorNoBase,
		},
		// Scenario: Missing root pointerfile. We expect an ErrorNoRoot.
		{
			name:        "missing-rootPointer",
			environment: "environments/default",
			testdata:    []string{"environments/default/main.jsonnet"},
			err:         ErrorNoRoot,
		},

		// Make sure jsonnetfile.json works as a pointer
		scenarioPointer("jsonnetfile.json"),
		// Make sure tkrc.yaml works as a pointer
		scenarioPointer("tkrc.yaml"),

		// Per-environment vendoring is tricky, because environments get their
		// own `jsonnetfile.json`, so `rootDir` would yield the same thing as
		// `baseDir`, which is usually not wanted.
		//
		// Scenario 1: No tkrc.yaml to mark the actual root. `baseDir` and
		// `rootDir` should be equal
		scenarioLocalVendor(false),
		// Scenario 2: A tkrc.yaml is added o the actual root. `rootDir` should
		// be the actual root, `baseDir` the environment.
		scenarioLocalVendor(true),
	}

	for _, s := range scenarios {
		require.NotZero(t, s.environment)

		t.Run(s.name, func(t *testing.T) {
			dir := makeTestdata(t, s.testdata)
			defer os.RemoveAll(dir)

			_, base, root, err := Resolve(filepath.Join(dir, s.environment))
			assert.Equal(t, s.err, err)

			if err == nil {
				assert.Equal(t, filepath.Join(dir, s.baseDir), base)
				assert.Equal(t, filepath.Join(dir, s.rootDir), root)
			}
		})
	}
}

func scenarioLocalVendor(tkrc bool) scenario {
	name := "localvendor"
	td := []string{
		"jsonnetfile.json",
		"environments/default/main.jsonnet",
		"environments/default/jsonnetfile.json",
	}
	// first jsonnetfile.json is in baseDir, so it will become rootDir as well
	root := "environments/default"

	if tkrc {
		name += "-with-tkrc"
		td = append(td, "tkrc.yaml") // add tkrc
		// now root should be project_root instead
		root = "/"
	}

	return scenario{
		name:        name,
		environment: "environments/default",
		testdata:    td,

		rootDir: root,
		baseDir: "environments/default",
	}
}

func scenarioPointer(ptr string) scenario {
	return scenario{
		name:        "pointer-" + ptr,
		environment: "environments/default",

		testdata: []string{
			"environments/default/main.jsonnet",
			ptr,
		},

		baseDir: "environments/default",
		rootDir: "/",
	}
}

func makeTestdata(t *testing.T, td []string) string {
	t.Helper()

	tmp, err := ioutil.TempDir("", "tk-dirsTest")
	require.NoError(t, err)

	for _, f := range td {
		dir, file := filepath.Split(f)
		if dir != "" {
			err := os.MkdirAll(filepath.Join(tmp, dir), os.ModePerm)
			require.NoError(t, err)
		}
		if file != "" {
			_, err := os.Create(filepath.Join(tmp, dir, file))
			require.NoError(t, err)
		}
	}
	return tmp
}
