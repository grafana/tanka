package native

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	jsonnet "github.com/google/go-jsonnet"
	"github.com/stretchr/testify/assert"
)

// callNative calls a native function used by jsonnet VM.
func callNative(name string, data []interface{}) (res interface{}, err error, callerr error) {
	for _, fun := range Funcs() {
		if fun.Name == name {
			// Call the function
			ret, err := fun.Func(data)
			return ret, err, nil
		}
	}

	return nil, nil, fmt.Errorf("could not find native function %s", name)
}

// callVMNative calls a native function used by jsonnet VM that requires access to the VM resource
func callVMNative(name string, data []interface{}) (res interface{}, err error, callerr error) {
	vm := jsonnet.MakeVM()
	for _, fun := range VMFuncs(vm) {
		if fun.Name == name {
			// Call the function
			ret, err := fun.Func(data)
			return ret, err, nil
		}
	}

	return nil, nil, fmt.Errorf("could not find VM native function %s", name)
}

func TestSha256(t *testing.T) {
	ret, err, callerr := callNative("sha256", []interface{}{"foo"})

	assert.Empty(t, callerr)
	assert.Equal(t, "2c26b46b68ffc68ff99b453c1d30413413422d706483bfa0f98a5e886266e7ae", ret)
	assert.Empty(t, err)
}

func TestParseJSONEmptyDict(t *testing.T) {
	ret, err, callerr := callNative("parseJson", []interface{}{"{}"})

	assert.Empty(t, callerr)
	assert.Equal(t, map[string]interface{}{}, ret)
	assert.Empty(t, err)
}

func TestParseJSONkeyValuet(t *testing.T) {
	ret, err, callerr := callNative("parseJson", []interface{}{"{\"a\": 47}"})

	assert.Empty(t, callerr)
	assert.Equal(t, map[string]interface{}{"a": 47.0}, ret)
	assert.Empty(t, err)
}

func TestParseJSONInvalid(t *testing.T) {
	ret, err, callerr := callNative("parseJson", []interface{}{""})

	assert.Empty(t, callerr)
	assert.Empty(t, ret)
	assert.IsType(t, &json.SyntaxError{}, err)
}

func TestParseYAMLEmpty(t *testing.T) {
	ret, err, callerr := callNative("parseYaml", []interface{}{""})

	assert.Empty(t, callerr)
	assert.Equal(t, []interface{}{}, ret)
	assert.Empty(t, err)
}

func TestParseYAMLKeyValue(t *testing.T) {
	ret, err, callerr := callNative("parseYaml", []interface{}{"a: 47"})

	assert.Empty(t, callerr)
	assert.Equal(t, []interface{}{map[string]interface{}{"a": 47.0}}, ret)
	assert.Empty(t, err)
}

func TestParseYAMLInvalid(t *testing.T) {
	ret, err, callerr := callNative("parseYaml", []interface{}{"'"})

	assert.Empty(t, callerr)
	assert.Empty(t, ret)
	assert.NotEmpty(t, err)
}

func TestManifestJSONFromJSON(t *testing.T) {
	ret, err, callerr := callNative("manifestJsonFromJson", []interface{}{"{}", float64(4)})

	assert.Empty(t, callerr)
	assert.Equal(t, "{}\n", ret)
	assert.Empty(t, err)
}

func TestManifestJSONFromJSONReindent(t *testing.T) {
	ret, err, callerr := callNative("manifestJsonFromJson", []interface{}{"{ \"a\": 47}", float64(4)})

	assert.Empty(t, callerr)
	assert.Equal(t, "{\n    \"a\": 47\n}\n", ret)
	assert.Empty(t, err)
}

func TestManifestJSONFromJSONInvalid(t *testing.T) {
	ret, err, callerr := callNative("manifestJsonFromJson", []interface{}{"", float64(4)})

	assert.Empty(t, callerr)
	assert.Empty(t, ret)
	assert.NotEmpty(t, err)
}

func TestManifestYAMLFromJSONEmpty(t *testing.T) {
	ret, err, callerr := callNative("manifestYamlFromJson", []interface{}{"{}"})

	assert.Empty(t, callerr)
	assert.Equal(t, "{}\n", ret)
	assert.Empty(t, err)
}

func TestManifestYAMLFromJSONKeyValue(t *testing.T) {
	ret, err, callerr := callNative("manifestYamlFromJson", []interface{}{"{ \"a\": 47}"})

	assert.Empty(t, callerr)
	assert.Equal(t, "a: 47\n", ret)
	assert.Empty(t, err)
}

func TestManifestYAMLFromJSONList(t *testing.T) {
	ret, err, callerr := callNative("manifestYamlFromJson", []interface{}{`{ "list": ["a", "b", "c"]}`})

	assert.Empty(t, callerr)
	assert.Equal(t, `list:
    - a
    - b
    - c
`, ret)
	assert.Empty(t, err)
}

func TestManifestYAMLFromJSONInvalid(t *testing.T) {
	ret, err, callerr := callNative("manifestYamlFromJson", []interface{}{""})

	assert.Empty(t, callerr)
	assert.Empty(t, ret)
	assert.NotEmpty(t, err)
}

func TestEscapeStringRegex(t *testing.T) {
	ret, err, callerr := callNative("escapeStringRegex", []interface{}{""})

	assert.Empty(t, callerr)
	assert.Equal(t, "", ret)
	assert.Empty(t, err)
}

func TestEscapeStringRegexValue(t *testing.T) {
	ret, err, callerr := callNative("escapeStringRegex", []interface{}{"([0-9]+).*\\s"})

	assert.Empty(t, callerr)
	assert.Equal(t, "\\(\\[0-9\\]\\+\\)\\.\\*\\\\s", ret)
	assert.Empty(t, err)
}

func TestEscapeStringRegexInvalid(t *testing.T) {
	ret, err, callerr := callNative("escapeStringRegex", []interface{}{"([0-9]+"})

	assert.Empty(t, callerr)
	assert.Equal(t, "\\(\\[0-9\\]\\+", ret)
	assert.Empty(t, err)
}

func TestRegexMatch(t *testing.T) {
	ret, err, callerr := callNative("regexMatch", []interface{}{"", "a"})

	assert.Empty(t, callerr)
	assert.Equal(t, true, ret)
	assert.Empty(t, err)
}

func TestRegexMatchNoMatch(t *testing.T) {
	ret, err, callerr := callNative("regexMatch", []interface{}{"a", "b"})

	assert.Empty(t, callerr)
	assert.Equal(t, false, ret)
	assert.Empty(t, err)
}

func TestRegexMatchInvalidRegex(t *testing.T) {
	ret, err, callerr := callNative("regexMatch", []interface{}{"[0-", "b"})

	assert.Empty(t, callerr)
	assert.Empty(t, ret)
	assert.NotEmpty(t, err)
}

func TestRegexSubstNoChange(t *testing.T) {
	ret, err, callerr := callNative("regexSubst", []interface{}{"a", "b", "c"})

	assert.Empty(t, callerr)
	assert.Equal(t, "b", ret)
	assert.Empty(t, err)
}

func TestRegexSubstValid(t *testing.T) {
	ret, err, callerr := callNative("regexSubst", []interface{}{"p[^m]*", "pm", "poe"})

	assert.Empty(t, callerr)
	assert.Equal(t, "poem", ret)
	assert.Empty(t, err)
}

func TestRegexSubstInvalid(t *testing.T) {
	ret, err, callerr := callNative("regexSubst", []interface{}{"p[^m*", "pm", "poe"})

	assert.Empty(t, callerr)
	assert.Empty(t, ret)
	assert.NotEmpty(t, err)
}

func TestImportFiles(t *testing.T) {
	cwd, err := os.Getwd()
	assert.NoError(t, err)
	tempDir, err := os.MkdirTemp("", "importFilesTest")
	assert.NoError(t, err)
	defer func() {
		if err := os.Chdir(cwd); err != nil {
			panic(err)
		}
		os.RemoveAll(tempDir)
	}()
	err = os.Chdir(tempDir)
	assert.NoError(t, err)
	importDirName := "imports"
	importDir := filepath.Join(tempDir, importDirName)
	err = os.Mkdir(importDir, 0750)
	assert.NoError(t, err)
	importFiles := []string{"test1.libsonnet", "test2.libsonnet"}
	excludeFiles := []string{"skip1.libsonnet", "skip2.libsonnet"}
	for i, fName := range append(importFiles, excludeFiles...) {
		fPath := filepath.Join(importDir, fName)
		content := fmt.Sprintf("{ test: %d }", i)
		err = os.WriteFile(fPath, []byte(content), 0644)
		assert.NoError(t, err)
	}
	opts := make(map[string]interface{})
	opts["calledFrom"] = filepath.Join(tempDir, "main.jsonnet")
	opts["exclude"] = excludeFiles
	ret, err, callerr := callVMNative("importFiles", []interface{}{importDirName, opts})
	assert.NoError(t, err)
	assert.Nil(t, callerr)
	importMap, ok := ret.(map[string]interface{})
	assert.True(t, ok)
	for i, fName := range importFiles {
		content, ok := importMap[fName]
		assert.True(t, ok)
		cMap, ok := content.(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, cMap["test"], float64(i))
	}
	// Make sure excluded files were not imported
	for _, fName := range excludeFiles {
		_, ok = importMap[fName]
		assert.False(t, ok)
	}
}
