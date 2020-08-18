package helmraiser

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfToArgs_noconf(t *testing.T) {
	conf := HelmConf{}
	args, tempFiles, err := confToArgs(conf)
	for _, file := range tempFiles {
		defer os.Remove(file)
	}

	assert.Equal(t, []string(nil), args)
	assert.Nil(t, err)
}

func TestConfToArgs_emptyconf(t *testing.T) {
	conf := HelmConf{
		Values: map[string]interface{}{},
		Flags:  []string{},
	}

	args, tempFiles, err := confToArgs(conf)
	for _, file := range tempFiles {
		defer os.Remove(file)
	}

	assert.Equal(t, []string(nil), args)
	assert.Nil(t, err)
}

func TestConfToArgs_flags(t *testing.T) {
	conf := HelmConf{
		Flags: []string{
			"--version=v0.1",
			"--random=arg",
		},
	}

	args, tempFiles, err := confToArgs(conf)
	for _, file := range tempFiles {
		defer os.Remove(file)
	}

	assert.Equal(t, []string{
		"--version=v0.1",
		"--random=arg",
	}, args)
	assert.Nil(t, err)
}

func TestConfToArgs_values(t *testing.T) {
	conf := HelmConf{
		Values: map[string]interface{}{
			"hasValues": "yes",
		},
	}

	args, tempFiles, err := confToArgs(conf)
	for _, file := range tempFiles {
		defer os.Remove(file)
	}

	assert.FileExists(t, tempFiles[0])
	assert.Equal(t, []string{fmt.Sprintf("--values=%s", tempFiles[0])}, args)
	assert.Nil(t, err)
}

func TestConfToArgs_flagsvalues(t *testing.T) {
	conf := HelmConf{
		Values: map[string]interface{}{
			"hasValues": "yes",
		},
		Flags: []string{
			"--version=v0.1",
			"--random=arg",
		},
	}

	args, tempFiles, err := confToArgs(conf)
	for _, file := range tempFiles {
		defer os.Remove(file)
	}

	assert.Equal(t, []string{
		fmt.Sprintf("--values=%s", tempFiles[0]),
		"--version=v0.1",
		"--random=arg",
	}, args)
	assert.Nil(t, err)
}

func TestParseYamlToMap_basic(t *testing.T) {
	yamlFile := []byte(`---
kind: testKind
metadata:
  name: testName`)
	actual, err := parseYamlToMap(yamlFile)

	expected := map[string]interface{}{
		"testname_testkind": map[string]interface{}{
			"kind": "testKind",
			"metadata": map[string]interface{}{
				"name": "testName",
			},
		},
	}
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseYamlToMap_dash(t *testing.T) {
	yamlFile := []byte(`---
kind: testKind
metadata:
  name: test-Name`)
	actual, err := parseYamlToMap(yamlFile)

	expected := map[string]interface{}{
		"test_name_testkind": map[string]interface{}{
			"kind": "testKind",
			"metadata": map[string]interface{}{
				"name": "test-Name",
			},
		},
	}
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseYamlToMap_colon(t *testing.T) {
	yamlFile := []byte(`---
kind: testKind
metadata:
  name: test:Name`)
	actual, err := parseYamlToMap(yamlFile)

	expected := map[string]interface{}{
		"test_name_testkind": map[string]interface{}{
			"kind": "testKind",
			"metadata": map[string]interface{}{
				"name": "test:Name",
			},
		},
	}
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseYamlToMap_empty(t *testing.T) {
	yamlFile := []byte(`---`)
	actual, err := parseYamlToMap(yamlFile)

	expected := map[string]interface{}{}
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseYamlToMap_multiple_files(t *testing.T) {
	yamlFile := []byte(`---
kind: testKind
metadata:
  name: testName
---
kind: testKind
metadata:
  name: testName2`)
	actual, err := parseYamlToMap(yamlFile)

	expected := map[string]interface{}{
		"testname_testkind": map[string]interface{}{
			"kind": "testKind",
			"metadata": map[string]interface{}{
				"name": "testName",
			},
		},
		"testname2_testkind": map[string]interface{}{
			"kind": "testKind",
			"metadata": map[string]interface{}{
				"name": "testName2",
			},
		},
	}
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}
