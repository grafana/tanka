package native

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfToArgs_noargs_noconf(t *testing.T) {
	args := []string{}
	conf := map[string]interface{}{}

	err := confToArgs(conf, &args)

	assert.Equal(t, []string{}, args)
	assert.Nil(t, err)
}

func TestConfToArgs_noargs_emptyconf(t *testing.T) {
	args := []string{}
	conf := map[string]interface{}{
		"values": map[string]interface{}{},
		"flags":  []interface{}{},
	}

	err := confToArgs(conf, &args)

	assert.Equal(t, []string{}, args)
	assert.Nil(t, err)
}

func TestConfToArgs_args_emptyconf(t *testing.T) {
	args := []string{
		"helm",
		"template",
		"name",
		"chart",
	}
	conf := map[string]interface{}{
		"values": map[string]interface{}{},
		"flags":  []interface{}{},
	}

	err := confToArgs(conf, &args)

	assert.Equal(t, []string{
		"helm",
		"template",
		"name",
		"chart",
	}, args)
	assert.Nil(t, err)
}

func TestConfToArgs_args_flags(t *testing.T) {
	args := []string{
		"helm",
		"template",
		"name",
		"chart",
	}
	conf := map[string]interface{}{
		"flags": []interface{}{
			"--version=v0.1",
			"--random=arg",
		},
	}

	err := confToArgs(conf, &args)

	assert.Equal(t, []string{
		"helm",
		"template",
		"name",
		"chart",
		"--version=v0.1",
		"--random=arg",
	}, args)
	assert.Nil(t, err)
}

func TestConfToArgs_args_values(t *testing.T) {
	args := []string{
		"helm",
		"template",
		"name",
		"chart",
	}
	conf := map[string]interface{}{
		"values": map[string]interface{}{
			"hasValues": "yes",
		},
	}

	err := confToArgs(conf, &args)

	assert.FileExists(t, strings.Split(args[len(args)-1], "=")[1])
	assert.Regexp(t, "^--values=/", args[len(args)-1])
	assert.Nil(t, err)
}

func TestConfToArgs_args_flagsvalues(t *testing.T) {
	args := []string{
		"helm",
		"template",
		"name",
		"chart",
	}
	conf := map[string]interface{}{
		"values": map[string]interface{}{
			"hasValues": "yes",
		},
		"flags": []interface{}{
			"--version=v0.1",
			"--random=arg",
		},
	}

	err := confToArgs(conf, &args)

	assert.Equal(t, 7, len(args))
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
