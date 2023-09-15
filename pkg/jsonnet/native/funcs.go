package native

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strings"

	jsonnet "github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
	"github.com/grafana/tanka/pkg/helm"
	"github.com/grafana/tanka/pkg/kustomize"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v3"
)

// Funcs returns a slice of native Go functions that shall be available
// from Jsonnet using `std.nativeFunc`
func Funcs() []*jsonnet.NativeFunction {
	return []*jsonnet.NativeFunction{
		// Parse serialized data into dicts
		parseJSON(),
		parseYAML(),

		// Convert serializations
		manifestJSONFromJSON(),
		manifestYAMLFromJSON(),

		// Regular expressions
		escapeStringRegex(),
		regexMatch(),
		regexSubst(),

		// Hash functions
		hashSha256(),

		helm.NativeFunc(helm.ExecHelm{}),
		kustomize.NativeFunc(kustomize.ExecKustomize{}),
	}
}

// wrapNativeFunc takes a function name, names of the parameters (in a single
// comma-separated string), and an implementation - a function with return
// types (interface{}, error).
// It produces a jsonnet.NativeFunction which has:
//   - Name from given name
//   - Params derived from given paramNamesStr
//   - Func that converts the untyped parameters from the jsonnet side to the
//     types impl expects, passes them to impl (note: the values are converted,
//     not just assigned). When unexpected number and/or types of parameters
//     are passed on the jsonnet side, the Func never calls impl and returns a
//     human-readable error instead.
func wrapNativeFunc(name, paramNamesStr string, impl interface{}) *jsonnet.NativeFunction {
	implV := reflect.ValueOf(impl)
	implT := implV.Type()
	if implV.Kind() != reflect.Func || implV.IsNil() {
		panic(fmt.Errorf("wrapNativeFunc(%s): not a non-nil function", name))
	}
	var paramNames ast.Identifiers
	for _, name := range strings.Split(paramNamesStr, ",") {
		paramNames = append(paramNames, ast.Identifier(name))
	}
	if implT.NumIn() != len(paramNames) {
		panic(fmt.Errorf("wrapNativeFunc(%s): wrong number of input parameters", name))
	}
	var goodOutTypesFunc func() (interface{}, error)
	outTypesT := reflect.TypeOf(goodOutTypesFunc)
	if implT.NumOut() != outTypesT.NumOut() ||
		implT.Out(0) != outTypesT.Out(0) || implT.Out(1) != outTypesT.Out(1) {
		panic(fmt.Errorf("wrapNativeFunc(%s): incorrect return parameters", name))
	}
	return &jsonnet.NativeFunction{
		Name:   name,
		Params: paramNames,
		Func: func(params []interface{}) (interface{}, error) {
			if len(params) != implT.NumIn() {
				// Bug. Jsonnet should call us with the correct parameters.
				panic(fmt.Errorf("%s(): wrong number of parameters", name))
			}
			callParams := make([]reflect.Value, len(params))
			for i, v := range params {
				if v == nil {
					return nil, fmt.Errorf("%s(): argument %#v is null", name, paramNames[i])
				}
				dstT := implT.In(i)
				srcV := reflect.ValueOf(v)
				if !srcV.CanConvert(dstT) {
					return nil, fmt.Errorf("%s(): argument %#v has unexpected type", name, paramNames[i])
				}
				callParams[i] = srcV.Convert(dstT)
			}
			results := implV.Call(callParams)
			var outErr error
			reflect.ValueOf(&outErr).Elem().Set(results[1])
			return results[0].Interface(), outErr
		},
	}
}

// parseJSON wraps `json.Unmarshal` to convert a json string into a dict
func parseJSON() *jsonnet.NativeFunction {
	return wrapNativeFunc(
		"parseJson",
		"json",
		func(data []byte) (res interface{}, err error) {
			err = json.Unmarshal(data, &res)
			return
		},
	)
}

func hashSha256() *jsonnet.NativeFunction {
	return wrapNativeFunc(
		"sha256",
		"str",
		func(data []byte) (interface{}, error) {
			h := sha256.New()
			h.Write(data)
			return fmt.Sprintf("%x", h.Sum(nil)), nil
		},
	)
}

// parseYAML wraps `yaml.Unmarshal` to convert a string of yaml document(s) into a (set of) dicts
func parseYAML() *jsonnet.NativeFunction {
	return wrapNativeFunc(
		"parseYaml",
		"yaml",
		func(data []byte) (interface{}, error) {
			ret := []interface{}{}

			d := yaml.NewDecoder(bytes.NewReader(data))
			for {
				var doc, jsonDoc interface{}
				if err := d.Decode(&doc); err != nil {
					if err == io.EOF {
						break
					}
					return nil, errors.Wrap(err, "parsing yaml")
				}

				jsonRaw, err := json.Marshal(doc)
				if err != nil {
					return nil, errors.Wrap(err, "converting yaml to json")
				}

				if err := json.Unmarshal(jsonRaw, &jsonDoc); err != nil {
					return nil, errors.Wrap(err, "converting yaml to json")
				}

				ret = append(ret, jsonDoc)
			}

			return ret, nil
		},
	)
}

// manifestJSONFromJSON reserializes JSON which allows to change the indentation.
func manifestJSONFromJSON() *jsonnet.NativeFunction {
	return wrapNativeFunc(
		"manifestJsonFromJson",
		"json,indent",
		func(data []byte, indent int) (interface{}, error) {
			data = bytes.TrimSpace(data)
			buf := bytes.Buffer{}
			if err := json.Indent(&buf, data, "", strings.Repeat(" ", indent)); err != nil {
				return "", err
			}
			buf.WriteString("\n")
			return buf.String(), nil
		},
	)
}

// manifestYamlFromJSON serializes a JSON string as a YAML document
func manifestYAMLFromJSON() *jsonnet.NativeFunction {
	return wrapNativeFunc(
		"manifestYamlFromJson",
		"json",
		func(data []byte) (interface{}, error) {
			var input interface{}
			if err := json.Unmarshal(data, &input); err != nil {
				return "", err
			}
			output, err := yaml.Marshal(input)
			return string(output), err
		},
	)
}

// escapeStringRegex escapes all regular expression metacharacters
// and returns a regular expression that matches the literal text.
func escapeStringRegex() *jsonnet.NativeFunction {
	return wrapNativeFunc(
		"escapeStringRegex",
		"str",
		func(s string) (interface{}, error) {
			return regexp.QuoteMeta(s), nil
		},
	)
}

// regexMatch returns whether the given string is matched by the given re2 regular expression.
func regexMatch() *jsonnet.NativeFunction {
	return wrapNativeFunc(
		"regexMatch",
		"regex,string",
		func(regex, s string) (interface{}, error) {
			return regexp.MatchString(regex, s)
		},
	)
}

// regexSubst replaces all matches of the re2 regular expression with another string.
func regexSubst() *jsonnet.NativeFunction {
	return wrapNativeFunc(
		"regexSubst",
		"regex,src,repl",
		func(regex, src, repl string) (interface{}, error) {
			r, err := regexp.Compile(regex)
			if err != nil {
				return "", err
			}
			return r.ReplaceAllString(src, repl), nil
		},
	)
}
