package kubernetes

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func loadFixture(t *testing.T, fixtureName string) map[string]interface{} {
	data, err := ioutil.ReadFile(path.Join("testdata", fixtureName))
	if err != nil {
		t.Errorf("failed to read fixture: %v", err)
	}

	var fixture map[string]interface{}
	if err := json.Unmarshal(data, &fixture); err != nil {
		t.Errorf("failed to parse fixture: %v", err)
	}

	return fixture
}

func TestWalkJSON(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]interface{}
		expected []Manifest
	}{
		{
			name: "regular dictionary of kubernetes resources",
			data: map[string]interface{}{
				"deployment": loadFixture(t, "deployment.json"),
				"service":    loadFixture(t, "service.json"),
			},
			expected: []Manifest{
				Manifest(loadFixture(t, "deployment.json")),
				Manifest(loadFixture(t, "service.json")),
			},
		},
		{
			name: "deeply nested dictionaries of kubernetes resources",
			data: map[string]interface{}{
				"top": map[string]interface{}{
					"deployment": loadFixture(t, "deployment.json"),
					"service":    loadFixture(t, "service.json"),
				},
			},
			expected: []Manifest{
				Manifest(loadFixture(t, "deployment.json")),
				Manifest(loadFixture(t, "service.json")),
			},
		},
		{
			name: "with lists of kubernetes resources within the original dictionary",
			data: map[string]interface{}{
				"resources": []interface{}{
					loadFixture(t, "deployment.json"),
					loadFixture(t, "service.json"),
				},
			},
			expected: []Manifest{
				Manifest(loadFixture(t, "deployment.json")),
				Manifest(loadFixture(t, "service.json")),
			},
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			got, _ := walkJSON(c.data)
			assert.ElementsMatch(t, c.expected, got)
		})
	}
}
