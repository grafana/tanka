package spec

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/grafana/tanka/pkg/spec/v1alpha1"
)

// list of deprecated config keys and their alternatives
// however, they still work and are aliased internally
var deprecated = []depreciation{
	{old: "namespace", new: "spec.namespace"},
	{old: "server", new: "spec.apiServer"},
	{old: "team", new: "metadata.labels.team"},
}

// Parse parses the json `data` into a `v1alpha1.Config` object.
// `baseDir` is the path of the environment
func Parse(data []byte, baseDir string) (*v1alpha1.Config, error) {
	v := viper.New()
	v.SetConfigType("json")
	if err := v.ReadConfig(bytes.NewReader(data)); err != nil {
		return nil, err
	}
	return parse(v, baseDir)
}

// ParseDir parses the given environments `spec.json` into a `v1alpha1.Config`
// object with the name set to the directories name
func ParseDir(baseDir string) (*v1alpha1.Config, error) {
	fi, err := os.Stat(baseDir)
	if err != nil {
		return nil, err
	}
	if !fi.IsDir() {
		return nil, errors.New("baseDir is not an directory")
	}

	v := viper.New()
	v.SetConfigName("spec")
	v.AddConfigPath(baseDir)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	return parse(v, baseDir)
}

// parse accepts a viper.Viper already loaded with the actual config and
// unmarshals it onto a v1alpha1.Config
func parse(v *viper.Viper, baseDir string) (*v1alpha1.Config, error) {
	var errDepr ErrDeprecated

	// handle deprecated ksonnet spec
	for _, d := range deprecated {
		if v.IsSet(d.old) && !v.IsSet(d.new) {
			if errDepr == nil {
				errDepr = ErrDeprecated{d}
			} else {
				errDepr = append(errDepr, d)
			}
			v.Set(d.new, v.Get(d.old))
		}
	}

	config := v1alpha1.New()
	if err := v.Unmarshal(config); err != nil {
		return nil, errors.Wrap(err, "parsing spec.json")
	}

	// set the name field
	config.Metadata.Name = filepath.Base(baseDir)
	config.Metadata.Labels["path"] = baseDir

	// return depreciation notes in case any exist as well
	return config, errDepr
}
