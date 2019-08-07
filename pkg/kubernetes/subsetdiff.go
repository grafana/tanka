package kubernetes

import (
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"
	"github.com/stretchr/objx"
	yaml "gopkg.in/yaml.v2"

	"github.com/sh0rez/tanka/pkg/util"
)

type difference struct {
	live, merged string
}

func (k Kubectl) SubsetDiff(y string) (string, error) {
	docs := map[string]difference{}
	d := yaml.NewDecoder(strings.NewReader(y))
	for {

		// jsonnet output -> desired state
		var rawShould map[interface{}]interface{}
		err := d.Decode(&rawShould)
		if err == io.EOF {
			break
		}

		if err != nil {
			return "", errors.Wrap(err, "decoding yaml")
		}

		// filename
		m := objx.New(util.CleanupInterfaceMap(rawShould))
		name := strings.Replace(fmt.Sprintf("%s.%s.%s.%s",
			m.Get("apiVersion").MustStr(),
			m.Get("kind").MustStr(),
			m.Get("metadata.namespace").MustStr(),
			m.Get("metadata.name").MustStr(),
		), "/", "-", -1)

		// kubectl output -> current state
		rawIs, err := k.Get(
			m.Get("metadata.namespace").MustStr(),
			m.Get("kind").MustStr(),
			m.Get("metadata.name").MustStr(),
		)
		if err != nil {
			if _, ok := err.(ErrorNotFound); ok {
				rawIs = map[string]interface{}{}
			} else {
				return "", errors.Wrap(err, "getting state from cluster")
			}
		}

		should, err := yaml.Marshal(rawShould)
		if err != nil {
			return "", err
		}

		is, err := yaml.Marshal(subset(m, rawIs))
		if err != nil {
			return "", err
		}
		if string(is) == "{}\n" {
			is = []byte("")
		}
		docs[name] = difference{string(is), string(should)}
	}

	s := ""
	for k, v := range docs {
		d, err := diff(k, v.live, v.merged)
		if err != nil {
			return "", errors.Wrap(err, "invoking diff")
		}
		if d != "" {
			d += "\n"
		}
		s += d
	}

	return s, nil
}

// subset removes all keys from is, that are not present in should.
// It makes is a subset of should.
// Kubernetes returns more keys than we can know about.
// This means, we need to remove all keys from the kubectl output, that are not present locally.
func subset(should, is map[string]interface{}) map[string]interface{} {
	if should["namespace"] != nil {
		is["namespace"] = should["namespace"]
	}
	for k, v := range is {
		if should[k] == nil {
			delete(is, k)
			continue
		}

		switch b := v.(type) {
		case map[string]interface{}:
			if a, ok := should[k].(map[string]interface{}); ok {
				is[k] = subset(a, b)
			}
		case []map[string]interface{}:
			for i := range b {
				if a, ok := should[k].([]map[string]interface{}); ok {
					b[i] = subset(a[i], b[i])
				}
			}
		case []interface{}:
			for i := range b {
				if a, ok := should[k].([]interface{}); ok {
					aa, ok := a[i].(map[string]interface{})
					if !ok {
						continue
					}
					bb, ok := b[i].(map[string]interface{})
					if !ok {
						continue
					}
					b[i] = subset(aa, bb)
				}
			}
		}
	}
	return is
}
