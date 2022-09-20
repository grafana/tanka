package process

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/tanka/pkg/jsonnet"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// testData holds data for tests
type testData struct {
	Deep interface{}                  `json:"deep"`
	Flat map[string]manifest.Manifest `json:"flat"`
}

func loadFixture(name string) testData {
	caser := cases.Title(language.English)
	filename := "./testdata/td" + caser.String(name) + ".jsonnet"

	vm := jsonnet.MakeVM(jsonnet.Opts{
		ImportPaths: []string{"./testdata"},
	})

	data, err := vm.EvaluateFile(filename)
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
