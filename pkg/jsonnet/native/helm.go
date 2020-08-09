package native

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

// helmTemplate wraps and runs `helm template`
// returns the generated manifests in a map
func helmTemplate() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   "helmTemplate",
		Params: ast.Identifiers{"name", "chart", "conf"},
		Func: func(data []interface{}) (interface{}, error) {
			name, chart := data[0].(string), data[1].(string)
			conf := data[2].(map[string]interface{})

			// the basic arguments to make this work
			args := []string{
				"template",
				name,
				chart,
			}

			// convert the values map into a yaml file
			if val, ok := conf["values"]; ok {
				dataValues := val.(interface{})
				valuesYaml, err := yaml.Marshal(dataValues)
				if err != nil {
					return "", err
				}
				tmpFile, err := ioutil.TempFile(os.TempDir(), "tanka-")
				if err != nil {
					return "", errors.Wrap(err, "cannot create temporary values.yaml")
				}
				defer os.Remove(tmpFile.Name())
				if _, err = tmpFile.Write(valuesYaml); err != nil {
					return "", errors.Wrap(err, "failed to write to temporary values.yaml")
				}
				if err := tmpFile.Close(); err != nil {
					return "", err
				}
				args = append(args, fmt.Sprintf("--values=%s", tmpFile.Name()))
			}

			// append custom flags to the arguments
			if val, ok := conf["flags"]; ok {
				dataFlags := val.([]interface{})
				flags := make([]string, 0)
				for _, f := range dataFlags {
					flags = append(flags, f.(string))
				}
				args = append(args, flags...)
			}

			cmd := exec.Command("helm", args...)
			buf := bytes.Buffer{}
			cmd.Stdout = &buf
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("whil running helm %s", strings.Join(args, " ")))
			}

			files := make(map[string]interface{})
			d := yaml.NewDecoder(bytes.NewReader(buf.Bytes()))
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

				// name_kind will be the key used in the map
				kindName := struct {
					Kind     string `json:"kind"`
					Metadata struct {
						Name string `json:"name"`
					} `json:"metadata"`
				}{}
				if err := json.Unmarshal(jsonRaw, &kindName); err != nil {
					return nil, errors.Wrap(err, "subtracting kind/name through unmarshaling")
				}

				// it is easier to use snake_case in jsonnet
				normalizeName := func(s string) string {
					s = strings.ReplaceAll(s, "-", "_")
					s = strings.ReplaceAll(s, ":", "_")
					s = strings.ToLower(s)
					return s
				}

				// create a map of resources instead of a list
				// this makes the files easier to modify in jsonnet
				if jsonDoc != nil {
					files[normalizeName(
						fmt.Sprintf("%s_%s",
							kindName.Metadata.Name,
							kindName.Kind,
						),
					)] = jsonDoc
				}
			}

			return files, nil
		},
	}
}
