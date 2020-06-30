package spec

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/pkg/errors"

	"github.com/grafana/tanka/pkg/spec/v1alpha1"
)

// APIGroup is the prefix used for `kind`
const APIGroup = "tanka.dev"

// Specfile is the filename for the environment config
const Specfile = "spec.json"

// ParseDir parses the given environments `spec.json` into a `v1alpha1.Config`
// object with the name set to the directories name
func ParseDir(baseDir, name string) (*v1alpha1.Config, error) {
	fi, err := os.Stat(baseDir)
	if err != nil {
		return nil, err
	}
	if !fi.IsDir() {
		return nil, errors.New("baseDir is not an directory")
	}

	data, err := ioutil.ReadFile(filepath.Join(baseDir, Specfile))
	if err != nil {
		if os.IsNotExist(err) {
			c := v1alpha1.New()
			c.Metadata.Name = name
			return c, ErrNoSpec{name}
		}
		return nil, err
	}

	return Parse(data, name)
}

// Parse parses the json `data` into a `v1alpha1.Config` object.
// `name` is the name of the environment
func Parse(data []byte, name string) (*v1alpha1.Config, error) {
	config := v1alpha1.New()
	if err := json.Unmarshal(data, config); err != nil {
		return nil, errors.Wrap(err, "parsing spec.json")
	}

	// set the name field
	config.Metadata.Name = name

	if err := handleDeprecated(config, data); err != nil {
		return config, err
	}

	// default apiServer URL to https
	if !regexp.MustCompile("^.+://").MatchString(config.Spec.APIServer) {
		config.Spec.APIServer = "https://" + config.Spec.APIServer
	}

	return config, nil
}

func handleDeprecated(c *v1alpha1.Config, data []byte) error {
	var errDepr ErrDeprecated

	var msi map[string]interface{}
	if err := json.Unmarshal(data, &msi); err != nil {
		return err
	}

	// namespace -> spec.namespace
	if n, ok := msi["namespace"]; ok && c.Spec.Namespace == "" {
		n, ok := n.(string)
		if !ok {
			return ErrMistypedField{"namespace", n}
		}

		errDepr = append(errDepr, depreciation{"namespace", "spec.namespace"})
		c.Spec.Namespace = n
	}

	// server -> spec.apiServer
	if s, ok := msi["server"]; ok && c.Spec.APIServer == "" {
		s, ok := s.(string)
		if !ok {
			return ErrMistypedField{"server", s}
		}

		errDepr = append(errDepr, depreciation{"server", "spec.apiServer"})
		c.Spec.APIServer = s
	}

	// team -> metadata.labels.team
	_, hasTeam := c.Metadata.Labels["team"]
	if t, ok := msi["team"]; ok && !hasTeam {
		t, ok := t.(string)
		if !ok {
			return ErrMistypedField{"team", t}
		}

		errDepr = append(errDepr, depreciation{"team", "metadata.labels.team"})
		c.Metadata.Labels["team"] = t
	}

	if len(errDepr) != 0 {
		return errDepr
	}

	return nil
}
