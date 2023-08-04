package process

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/grafana/tanka/pkg/jsonnet"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
)

// testData holds data for tests
type testData struct {
	Deep interface{}                  `json:"deep"`
	Flat map[string]manifest.Manifest `json:"flat"`
}

func loadFixture(name string) testData {
	filename := filepath.Join("./testdata", name)

	vm := jsonnet.VMPool.Get(jsonnet.Opts{
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
	return loadFixture("tdRegular.jsonnet")
}

// testDataFlat is a flat manifest that does not need reconciliation
func testDataFlat() testData {
	return loadFixture("tdFlat.jsonnet")
}

// testDataPrimitive is an invalid manifest, because it ends with a primitive
// without including required fields
func testDataPrimitive() testData {
	return loadFixture("tdInvalidPrimitive.jsonnet")
}

// testBadKindType is an invalid manifest, because it has an invalid `kind` value
func testBadKindType() testData {
	return loadFixture("tdBadKindType.jsonnet")
}

// testMissingAttribute is an invalid manifest, because it is missing the `kind`
func testMissingAttribute() testData {
	return loadFixture("tdMissingAttribute.jsonnet")
}

// testDataDeep is super deeply nested on multiple levels
func testDataDeep() testData {
	return loadFixture("tdDeep.jsonnet")
}

// testDataArray is an array of (deeply nested) dicts that should be fully
// flattened
func testDataArray() testData {
	return loadFixture("tdArray.jsonnet")
}
