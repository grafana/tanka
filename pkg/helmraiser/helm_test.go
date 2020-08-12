package helmraiser

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfToArgs_noargs_noconf(t *testing.T) {
	conf := map[string]interface{}{}
	args, tempFiles, err := confToArgs(conf)
	for _, file := range tempFiles {
		defer os.Remove(file)
	}

	assert.Equal(t, []string(nil), args)
	assert.Nil(t, err)
}

func TestConfToArgs_args_emptyconf(t *testing.T) {
	conf := map[string]interface{}{
		"values": map[string]interface{}{},
		"flags":  []interface{}{},
	}

	args, tempFiles, err := confToArgs(conf)
	for _, file := range tempFiles {
		defer os.Remove(file)
	}

	assert.Equal(t, []string(nil), args)
	assert.Nil(t, err)
}

func TestConfToArgs_args_flags(t *testing.T) {
	conf := map[string]interface{}{
		"flags": []interface{}{
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

func TestConfToArgs_args_values(t *testing.T) {
	conf := map[string]interface{}{
		"values": map[string]interface{}{
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

func TestConfToArgs_args_flagsvalues(t *testing.T) {
	conf := map[string]interface{}{
		"values": map[string]interface{}{
			"hasValues": "yes",
		},
		"flags": []interface{}{
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
	m, err := parseYamlToMap(yamlFile)

	expected := map[string]interface{}{
		"kind": "testKind",
		"metadata": map[string]interface{}{
			"name": "testName",
		},
	}
	assert.Equal(t, expected, m["testname_testkind"])
	assert.Nil(t, err)
}

func TestParseYamlToMap_dash(t *testing.T) {
	yamlFile := []byte(`---
kind: testKind
metadata:
  name: test-Name`)
	m, err := parseYamlToMap(yamlFile)

	expected := map[string]interface{}{
		"kind": "testKind",
		"metadata": map[string]interface{}{
			"name": "test-Name",
		},
	}
	assert.Equal(t, expected, m["test_name_testkind"])
	assert.Nil(t, err)
}

func TestParseYamlToMap_colon(t *testing.T) {
	yamlFile := []byte(`---
kind: testKind
metadata:
  name: test:Name`)
	m, err := parseYamlToMap(yamlFile)

	expected := map[string]interface{}{
		"kind": "testKind",
		"metadata": map[string]interface{}{
			"name": "test:Name",
		},
	}
	assert.Equal(t, expected, m["test_name_testkind"])
	assert.Nil(t, err)
}

func TestParseYamlToMap_empty(t *testing.T) {
	yamlFile := []byte(`---`)
	m, err := parseYamlToMap(yamlFile)

	expected := map[string]interface{}{}
	assert.Equal(t, expected, m)
	assert.Nil(t, err)
}
