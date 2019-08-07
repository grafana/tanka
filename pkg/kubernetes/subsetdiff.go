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
	name         string
	live, merged string
}

func (k Kubectl) SubsetDiff(y string) (string, error) {
	docs := []difference{}
	d := yaml.NewDecoder(strings.NewReader(y))

	rs := 0
	errs := make(chan error)
	results := make(chan difference)

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

		rs++
		go subsetDiff(k, rawShould, results, errs)

	}

	for rs > 0 {
		select {
		case d := <-results:
			docs = append(docs, d)
			rs--
		case err := <-errs:
			return "", errors.Wrap(err, "calculating subset")
		}
	}
	close(results)
	close(errs)

	s := ""
	for _, d := range docs {
		sd, err := diff(d.name, d.live, d.merged)
		if err != nil {
			return "", errors.Wrap(err, "invoking diff")
		}
		if sd != "" {
			sd += "\n"
		}
		s += sd
	}

	return s, nil
}

func subsetDiff(k Kubectl, rawShould map[interface{}]interface{}, r chan difference, e chan error) {
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
			e <- errors.Wrap(err, "getting state from cluster")
			return
		}
	}

	should, err := yaml.Marshal(rawShould)
	if err != nil {
		e <- err
		return
	}

	is, err := yaml.Marshal(subset(m, rawIs))
	if err != nil {
		e <- err
		return
	}
	if string(is) == "{}\n" {
		is = []byte("")
	}

	r <- difference{
		name:   name,
		live:   string(is),
		merged: string(should),
	}
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
