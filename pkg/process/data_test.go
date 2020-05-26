package process

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/grafana/tanka/pkg/jsonnet"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
)

// testData holds data for tests
type testData struct {
	Deep interface{}                  `json:"deep"`
	Flat map[string]manifest.Manifest `json:"flat"`
}

func loadFixture(name string) testData {
	raw, err := ioutil.ReadFile("./testdata/td" + strings.Title(name) + ".jsonnet")
	if err != nil {
		panic(fmt.Sprint("loading fixture:", err))
	}

	data, err := jsonnet.Evaluate(string(raw), []string{"./testdata"})
	if err != nil {
		panic(fmt.Sprint("loading fixture:", err))
	}

	var d testData
	if err := json.Unmarshal([]byte(data), &d); err != nil {
		panic(fmt.Sprint("loading fixture:", err))
	}

	return d
}

// testDataRegular is a regular output of jsonnet without special things, but it
// is nested.
func testDataRegular() testData {
	return loadFixture("regular")
}

// testDataFlat is a flat manifest that does not need reconciliation
func testDataFlat() testData {
	return loadFixture("flat")
}

// testDataPrimitive is an invalid manifest, because it ends with a primitive
// without including required fields
func testDataPrimitive() testData {
	return loadFixture("invalidPrimitive")
}

// testDataDeep is super deeply nested on multiple levels
func testDataDeep() testData {
	return loadFixture("deep")
}

// testDataArray is an array of (deeply nested) dicts that should be fully
// flattened
func testDataArray() testData {
	return loadFixture("array")
}
