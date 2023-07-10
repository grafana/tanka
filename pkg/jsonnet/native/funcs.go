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

// wrapNativeFunc returns a wrapper around impl, which must be a function with
// return types (interface{}, error).
// The wrapper converts the untyped parameters from jsonnet to the types impl
// expects, passes them to impl (note: the values are converted, not just
// assigned). When unexpected types are passed on the jsonnet side, the wrapper
// never calls impl and returns a human-readable error instead.
// The wrapper is meant to be set as the given NativeFunction's Func field.
func wrapNativeFunc(f *jsonnet.NativeFunction, impl interface{}) func([]interface{}) (interface{}, error) {
	implV := reflect.ValueOf(impl)
	implT := implV.Type()
	if implV.Kind() != reflect.Func || implV.IsNil() {
		panic(fmt.Errorf("wrapNativeFunc(%s): not a non-nil function", f.Name))
	}
	if implT.NumIn() != len(f.Params) {
		panic(fmt.Errorf("wrapNativeFunc(%s): wrong number of input parameters", f.Name))
	}
	var goodOutTypesFunc func() (interface{}, error)
	outTypesT := reflect.TypeOf(goodOutTypesFunc)
	if implT.NumOut() != outTypesT.NumOut() ||
		implT.Out(0) != outTypesT.Out(0) || implT.Out(1) != outTypesT.Out(1) {
		panic(fmt.Errorf("wrapNativeFunc(%s): incorrect return parameters", f.Name))
	}
	return func(params []interface{}) (interface{}, error) {
		if len(params) != implT.NumIn() {
			// Bug. Jsonnet should call us with the correct parameters.
			panic(fmt.Errorf("%s(): wrong number of parameters", f.Name))
		}
		callParams := make([]reflect.Value, len(params))
		for i, v := range params {
			if v == nil {
				return nil, fmt.Errorf("%s(): argument %#v is null", f.Name, f.Params[i])
			}
			dstT := implT.In(i)
			srcV := reflect.ValueOf(v)
			if !srcV.CanConvert(dstT) {
				return nil, fmt.Errorf("%s(): argument %#v has unexpected type", f.Name, f.Params[i])
			}
			callParams[i] = srcV.Convert(dstT)
		}
		results := implV.Call(callParams)
		var outErr error
		reflect.ValueOf(&outErr).Elem().Set(results[1])
		return results[0].Interface(), outErr
	}
}

// parseJSON wraps `json.Unmarshal` to convert a json string into a dict
func parseJSON() *jsonnet.NativeFunction {
	f := &jsonnet.NativeFunction{
		Name:   "parseJson",
		Params: ast.Identifiers{"json"},
	}
	f.Func = wrapNativeFunc(f,
		func(data []byte) (res interface{}, err error) {
			err = json.Unmarshal(data, &res)
			return
		})
	return f
}

func hashSha256() *jsonnet.NativeFunction {
	f := &jsonnet.NativeFunction{
		Name:   "sha256",
		Params: ast.Identifiers{"str"},
	}
	f.Func = wrapNativeFunc(f,
		func(data []byte) (interface{}, error) {
			h := sha256.New()
			h.Write(data)
			return fmt.Sprintf("%x", h.Sum(nil)), nil
		},
	)
	return f
}

// parseYAML wraps `yaml.Unmarshal` to convert a string of yaml document(s) into a (set of) dicts
func parseYAML() *jsonnet.NativeFunction {
	f := &jsonnet.NativeFunction{
		Name:   "parseYaml",
		Params: ast.Identifiers{"yaml"},
	}
	f.Func = wrapNativeFunc(f,
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
		})
	return f
}

// manifestJSONFromJSON reserializes JSON which allows to change the indentation.
func manifestJSONFromJSON() *jsonnet.NativeFunction {
	f := &jsonnet.NativeFunction{
		Name:   "manifestJsonFromJson",
		Params: ast.Identifiers{"json", "indent"},
	}
	f.Func = wrapNativeFunc(f,
		func(data []byte, indent int) (interface{}, error) {
			data = bytes.TrimSpace(data)
			buf := bytes.Buffer{}
			if err := json.Indent(&buf, data, "", strings.Repeat(" ", indent)); err != nil {
				return "", err
			}
			buf.WriteString("\n")
			return buf.String(), nil
		})
	return f
}

// manifestYamlFromJSON serializes a JSON string as a YAML document
func manifestYAMLFromJSON() *jsonnet.NativeFunction {
	f := &jsonnet.NativeFunction{
		Name:   "manifestYamlFromJson",
		Params: ast.Identifiers{"json"},
	}
	f.Func = wrapNativeFunc(f,
		func(data []byte) (interface{}, error) {
			var input interface{}
			if err := json.Unmarshal(data, &input); err != nil {
				return "", err
			}
			output, err := yaml.Marshal(input)
			return string(output), err
		})
	return f
}

// escapeStringRegex escapes all regular expression metacharacters
// and returns a regular expression that matches the literal text.
func escapeStringRegex() *jsonnet.NativeFunction {
	f := &jsonnet.NativeFunction{
		Name:   "escapeStringRegex",
		Params: ast.Identifiers{"str"},
	}
	f.Func = wrapNativeFunc(f,
		func(s string) (interface{}, error) {
			return regexp.QuoteMeta(s), nil
		})
	return f
}

// regexMatch returns whether the given string is matched by the given re2 regular expression.
func regexMatch() *jsonnet.NativeFunction {
	f := &jsonnet.NativeFunction{
		Name:   "regexMatch",
		Params: ast.Identifiers{"regex", "string"},
	}
	f.Func = wrapNativeFunc(f,
		func(regex, s string) (interface{}, error) {
			return regexp.MatchString(regex, s)
		})
	return f
}

// regexSubst replaces all matches of the re2 regular expression with another string.
func regexSubst() *jsonnet.NativeFunction {
	f := &jsonnet.NativeFunction{
		Name:   "regexSubst",
		Params: ast.Identifiers{"regex", "src", "repl"},
	}
	f.Func = wrapNativeFunc(f,
		func(regex, src, repl string) (interface{}, error) {
			r, err := regexp.Compile(regex)
			if err != nil {
				return "", err
			}
			return r.ReplaceAllString(src, repl), nil
		})
	return f
}
