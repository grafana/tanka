package helmraiser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	jsonnet "github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v3"
)

type HelmConf struct {
	Values map[string]interface{}
	Flags  []string
}

func confToArgs(conf HelmConf) ([]string, []string, error) {
	var args []string
	var tempFiles []string

	// create file and append to args
	if len(conf.Values) != 0 {
		valuesYaml, err := yaml.Marshal(conf.Values)
		if err != nil {
			return nil, nil, err
		}
		tmpFile, err := ioutil.TempFile(os.TempDir(), "tanka-")
		if err != nil {
			return nil, nil, errors.Wrap(err, "cannot create temporary values.yaml")
		}
		tempFiles = append(tempFiles, tmpFile.Name())
		if _, err = tmpFile.Write(valuesYaml); err != nil {
			return nil, tempFiles, errors.Wrap(err, "failed to write to temporary values.yaml")
		}
		if err := tmpFile.Close(); err != nil {
			return nil, tempFiles, err
		}
		args = append(args, fmt.Sprintf("--values=%s", tmpFile.Name()))
	}

	// append custom flags to args
	args = append(args, conf.Flags...)

	if len(args) == 0 {
		args = nil
	}

	return args, tempFiles, nil
}

func parseYamlToMap(yamlFile []byte) (map[string]interface{}, error) {
	files := make(map[string]interface{})
	d := yaml.NewDecoder(bytes.NewReader(yamlFile))
	for {
		var doc, jsonDoc interface{}
		if err := d.Decode(&doc); err != nil {
			if err == io.EOF {
				break
			}
			return nil, errors.Wrap(err, "parsing manifests")
		}

		jsonRaw, err := json.Marshal(doc)
		if err != nil {
			return nil, errors.Wrap(err, "marshaling mainfests")
		}

		if err := json.Unmarshal(jsonRaw, &jsonDoc); err != nil {
			return nil, errors.Wrap(err, "unmarshaling manifests")
		}

		// Unmarshal name and kind
		kindName := struct {
			Kind     string `json:"kind"`
			Metadata struct {
				Name string `json:"name"`
			} `json:"metadata"`
		}{}
		if err := json.Unmarshal(jsonRaw, &kindName); err != nil {
			return nil, errors.Wrap(err, "subtracting kind/name through unmarshaling")
		}

		// snake_case string
		normalizeName := func(s string) string {
			s = strings.ReplaceAll(s, "-", "_")
			s = strings.ReplaceAll(s, ":", "_")
			s = strings.ToLower(s)
			return s
		}

		// create a map of resources for ease of use in jsonnet
		name := normalizeName(fmt.Sprintf("%s_%s", kindName.Metadata.Name, kindName.Kind))
		if jsonDoc != nil {
			files[name] = jsonDoc
		}
	}
	return files, nil
}

// helmTemplate wraps and runs `helm template`
// returns the generated manifests in a map
func HelmTemplate() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name: "helmTemplate",
		// Lines up with `helm template [NAME] [CHART] [flags]` except 'conf' is a bit more elaborate
		Params: ast.Identifiers{"name", "chart", "conf"},
		Func: func(data []interface{}) (interface{}, error) {
			name, chart := data[0].(string), data[1].(string)

			c, err := json.Marshal(data[2])
			if err != nil {
				return "", err
			}
			var conf HelmConf
			if err := json.Unmarshal(c, &conf); err != nil {
				return "", err
			}

			// the basic arguments to make this work
			args := []string{
				"template",
				name,
				chart,
			}

			confArgs, tempFiles, err := confToArgs(conf)
			if err != nil {
				return "", nil
			}
			for _, file := range tempFiles {
				defer os.Remove(file)
			}
			if confArgs != nil {
				args = append(args, confArgs...)
			}

			helmBinary := "helm"
			if hc := os.Getenv("TANKA_HELM_PATH"); hc != "" {
				helmBinary = hc
			}

			// convert the values map into a yaml file
			cmd := exec.Command(helmBinary, args...)
			buf := bytes.Buffer{}
			cmd.Stdout = &buf
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("while running helm %s", strings.Join(args, " ")))
			}

			return parseYamlToMap(buf.Bytes())
		},
	}
}
