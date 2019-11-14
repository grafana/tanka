package kubernetes

import (
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"
	"github.com/stretchr/objx"
	yaml "gopkg.in/yaml.v3"

	"github.com/grafana/tanka/pkg/kubernetes/client"
)

type difference struct {
	name         string
	live, merged string
}

func SubsetDiffer(c client.Client) Differ {
	return func(state client.Manifests) (*string, error) {
		docs := []difference{}
		d := yaml.NewDecoder(strings.NewReader(state.String()))

		routines := 0
		errCh := make(chan error)
		resultCh := make(chan difference)

		for {
			// jsonnet output -> desired state
			var rawShould map[string]interface{}
			err := d.Decode(&rawShould)
			if err == io.EOF {
				break
			}

			if err != nil {
				return nil, errors.Wrap(err, "decoding yaml")
			}

			routines++
			go parallelSubsetDiff(c, rawShould, resultCh, errCh)
		}

		var lastErr error
		for i := 0; i < routines; i++ {
			select {
			case d := <-resultCh:
				docs = append(docs, d)
			case err := <-errCh:
				lastErr = err
			}
		}
		close(resultCh)
		close(errCh)

		if lastErr != nil {
			return nil, errors.Wrap(lastErr, "calculating subset")
		}

		var diffs string
		for _, d := range docs {
			diffStr, err := diff(d.name, d.live, d.merged)
			if err != nil {
				return nil, errors.Wrap(err, "invoking diff")
			}
			if diffStr != "" {
				diffStr += "\n"
			}
			diffs += diffStr
		}
		diffs = strings.TrimSuffix(diffs, "\n")

		return &diffs, nil
	}
}

func parallelSubsetDiff(c client.Client, rawShould map[string]interface{}, r chan difference, e chan error) {
	diff, err := subsetDiff(c, rawShould)
	if err != nil {
		e <- err
	}
	r <- *diff
}

func subsetDiff(c client.Client, rawShould map[string]interface{}) (*difference, error) {
	m := objx.New(rawShould)
	name := strings.Replace(fmt.Sprintf("%s.%s.%s.%s",
		m.Get("apiVersion").MustStr(),
		m.Get("kind").MustStr(),
		m.Get("metadata.namespace").MustStr(),
		m.Get("metadata.name").MustStr(),
	), "/", "-", -1)

	// kubectl output -> current state
	rawIs, err := c.Get(
		m.Get("metadata.namespace").MustStr(),
		m.Get("kind").MustStr(),
		m.Get("metadata.name").MustStr(),
	)
	if err != nil {
		if _, ok := err.(ErrorNotFound); ok {
			rawIs = map[string]interface{}{}
		} else {
			return nil, errors.Wrap(err, "getting state from cluster")
		}
	}

	should, err := yaml.Marshal(rawShould)
	if err != nil {
		return nil, err
	}

	is, err := yaml.Marshal(subset(m, rawIs))
	if err != nil {
		return nil, err
	}
	if string(is) == "{}\n" {
		is = []byte("")
	}

	return &difference{
		name:   name,
		live:   string(is),
		merged: string(should),
	}, nil
}

// subset removes all keys from is, that are not present in should.
// It makes is a subset of should.
// Kubernetes returns more keys than we can know about.
// This means, we need to remove all keys from the kubectl output, that are not present locally.
func subset(should, is map[string]interface{}) map[string]interface{} {
	if should["namespace"] != nil {
		is["namespace"] = should["namespace"]
	}

	// just ignore the apiVersion for now, too much bloat
	if should["apiVersion"] != nil && is["apiVersion"] != nil {
		is["apiVersion"] = should["apiVersion"]
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
					if i >= len(a) {
						// slice in config shorter than in live. Abort, as there are no entries to diff anymore
						break
					}

					// value not a dict, no recursion needed
					cShould, ok := a[i].(map[string]interface{})
					if !ok {
						continue
					}

					// value not a dict, no recursion needed
					cIs, ok := b[i].(map[string]interface{})
					if !ok {
						continue
					}
					b[i] = subset(cShould, cIs)
				}
			}
		}
	}
	return is
}
