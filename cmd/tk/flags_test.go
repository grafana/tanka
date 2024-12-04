package main

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestCliCodeParser(t *testing.T) {
	fs := pflag.NewFlagSet("test-cli-code-parser", pflag.ContinueOnError)
	parseExt, parseTLA := cliCodeParser(fs)
	err := fs.Parse([]string{
		"--ext-str", "es=1a \" \U0001f605 ' b\nc\u010f",
		"--tla-str", "ts=2a \" \U0001f605 ' b\nc\u010f",
		"--ext-code", "ec=1+2",
		"--tla-code", "tc=2+3",
		"-A", "ts2=ts2", // tla-str
		"-V", "es2=es2", // ext-str
		"--ext-str-file", `esf=e"sf.txt`,
		"--tla-str-file", `tsf=t"s"f.txt`,
		"--ext-code-file", `ecf=e"cf.json`,
		"--tla-code-file", `tcf=t"c"f.json`,
	})
	assert.NoError(t, err)
	ext := parseExt()
	assert.Equal(t, map[string]string{
		"es":  `"1a \" ` + "\U0001f605" + ` ' b\nc` + "\u010f" + `"`,
		"ec":  "1+2",
		"es2": `"es2"`,
		"esf": `importstr @"e""sf.txt"`,
		"ecf": `import @"e""cf.json"`,
	}, ext)
	tla := parseTLA()
	assert.Equal(t, map[string]string{
		"ts":  `"2a \" ` + "\U0001f605" + ` ' b\nc` + "\u010f" + `"`,
		"tc":  "2+3",
		"ts2": `"ts2"`,
		"tsf": `importstr @"t""s""f.txt"`,
		"tcf": `import @"t""c""f.json"`,
	}, tla)
}
