package helm

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfToArgs_noconf(t *testing.T) {
	conf := TemplateOpts{}
	args, tempFiles, err := confToArgs(conf)
	for _, file := range tempFiles {
		defer os.Remove(file)
	}

	assert.Equal(t, []string(nil), args)
	assert.Nil(t, err)
}

func TestConfToArgs_emptyconf(t *testing.T) {
	conf := TemplateOpts{
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
	conf := TemplateOpts{
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
	conf := TemplateOpts{
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
	conf := TemplateOpts{
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
