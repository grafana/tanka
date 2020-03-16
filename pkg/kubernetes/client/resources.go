package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
)

// Resources the Kubernetes API knows
type Resources []Resource

// Namespaced returns whether a resource is namespace-specific or cluster-wide
func (r Resources) Namespaced(m manifest.Manifest) bool {
	for _, res := range r {
		if m.Kind() == res.Kind {
			return res.Namespaced
		}
	}

	return false
}

// Resource is a Kubernetes API Resource
type Resource struct {
	ApiGroup   string `json:"APIGROUP"`
	Kind       string `json:"KIND"`
	Name       string `json:"NAME"`
	Namespaced bool   `json:"NAMESPACED,string"`
	Shortnames string `json:"SHORTNAMES"`
}

// Resources returns all API resources known to the server
func (k Kubectl) Resources() (Resources, error) {
	cmd := k.ctl("api-resources")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	var res Resources
	if err := UnmarshalTable(out.String(), &res); err != nil {
		return nil, errors.Wrap(err, "parsing table")
	}

	return res, nil
}

// UnmarshalTable unmarshals a raw CLI table into ptr. `json` package is used
// for getting the dat into the ptr, `json:` struct tags can be used.
func UnmarshalTable(raw string, ptr interface{}) error {
	raw = strings.TrimSpace(raw)

	lines := strings.Split(raw, "\n")
	if len(lines) < 2 {
		return errors.New("table has less than 2 lines. No content found")
	}

	headerStr := lines[0]
	lines = lines[1:]

	spc := regexp.MustCompile(`[A-Z]+\s*`)
	header := spc.FindAllString(headerStr, -1)

	var tbl []map[string]string
	for _, l := range lines {
		elems := splitRow(l, header)
		if len(elems) != len(header) {
			return fmt.Errorf("header and row have different element count: %v != %v", len(header), len(elems))
		}

		row := make(map[string]string)
		for i, e := range elems {
			key := strings.TrimSpace(header[i])
			row[key] = strings.TrimSpace(e)
		}
		tbl = append(tbl, row)
	}

	j, err := json.Marshal(tbl)
	if err != nil {
		return err
	}

	return json.Unmarshal(j, ptr)
}

func splitRow(s string, header []string) (elems []string) {
	pos := 0
	for i, h := range header {
		if i == len(header)-1 {
			elems = append(elems, s[pos:])
			continue
		}

		lim := len(h)
		elems = append(elems, s[pos:pos+lim])
		pos += lim
	}
	return elems
}
