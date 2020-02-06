package native

import (
	"encoding/json"
	"fmt"
	"testing"

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
