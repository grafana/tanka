package kubernetes

// This file  contains data for testing

// testData holds data for tests
type testData struct {
	deep, flat interface{}
}

// testDataRegular is a regular output of jsonnet without special things, but it
// is nested.
func testDataRegular() testData {
	return (testData{
		deep: map[string]interface{}{
			"deployment": map[string]interface{}{
				"apiVersion": "apps/v1",
				"kind":       "Deployment",
				"metadata": map[string]interface{}{
					"name": "nginx",
				},
				"spec": map[string]interface{}{
					"containers": []interface{}{
						map[string]interface{}{
							"name":  "nginx",
							"image": "nginx",
						},
					},
				},
			},
			"service": map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "Service",
				"metadata": map[string]interface{}{
					"name": "nginx",
				},
				"spec": map[string]interface{}{
					"selector": map[string]interface{}{
						"app": "app",
					},
				},
			},
		},
		flat: []map[string]interface{}{
			{
				"apiVersion": "apps/v1",
				"kind":       "Deployment",
				"metadata": map[string]interface{}{
					"name": "nginx",
				},
				"spec": map[string]interface{}{
					"containers": []interface{}{
						map[string]interface{}{
							"name":  "nginx",
							"image": "nginx",
						},
					},
				},
			},
			{
				"apiVersion": "v1",
				"kind":       "Service",
				"metadata": map[string]interface{}{
					"name": "nginx",
				},
				"spec": map[string]interface{}{
					"selector": map[string]interface{}{
						"app": "app",
					},
				},
			},
		},
	})
}

// testDataFlat is a flat manifest that does not need reconciliation
func testDataFlat() testData {
	return testData{
		deep: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name": "nginx",
			},
			"spec": map[string]interface{}{
				"containers": []interface{}{
					map[string]interface{}{
						"name":  "nginx",
						"image": "nginx",
					},
				},
			},
		},
		flat: []map[string]interface{}{
			{
				"apiVersion": "apps/v1",
				"kind":       "Deployment",
				"metadata": map[string]interface{}{
					"name": "nginx",
				},
				"spec": map[string]interface{}{
					"containers": []interface{}{
						map[string]interface{}{
							"name":  "nginx",
							"image": "nginx",
						},
					},
				},
			},
		},
	}
}

// testDataPrimitive is an invalid manifest, because it ends with a primitive
// without including required fields
func testDataPrimitive() testData {
	return testData{
		deep: map[string]interface{}{
			"nginx": map[string]interface{}{
				"deployment": map[string]interface{}{
					"apiVersion": "apps/v1",
					"kind":       "Deployment",
					"metadata": map[string]interface{}{
						"name": "nginx",
					},
				},
				"service": map[string]interface{}{
					"note": "invalid because apiVersion and kind are missing",
				},
			},
		},
		flat: []map[string]interface{}(nil),
	}
}

// testDataDeep is super deeply nested on multiple levels
func testDataDeep() testData {
	return testData{
		deep: map[string]interface{}{
			"app": map[string]interface{}{
				"web": map[string]interface{}{
					"backend": map[string]interface{}{
						"server": map[string]interface{}{
							"nginx": map[string]interface{}{
								"deployment": map[string]interface{}{
									"kind":       "Deployment",
									"apiVersion": "apps/v1",
									"metadata": map[string]interface{}{
										"name": "nginx",
									},
								},
							},
						},
					},
					"frontend": map[string]interface{}{
						"nodejs": map[string]interface{}{
							"express": map[string]interface{}{
								"service": map[string]interface{}{
									"kind":       "Service",
									"apiVersion": "v1",
									"metadata": map[string]interface{}{
										"name": "frontend",
									},
								},
								"deployment": map[string]interface{}{
									"kind":       "Deployment",
									"apiVersion": "apps/v1",
									"metadata": map[string]interface{}{
										"name": "frontend",
									},
								},
							},
						},
					},
				},
				"namespace": map[string]interface{}{
					"kind":       "Namespace",
					"apiVersion": "v1",
					"metadata": map[string]interface{}{
						"name": "app",
					},
				},
			},
		},
		flat: []map[string]interface{}{
			{
				"kind":       "Deployment",
				"apiVersion": "apps/v1",
				"metadata": map[string]interface{}{
					"name": "nginx",
				},
			},
			{
				"kind":       "Service",
				"apiVersion": "v1",
				"metadata": map[string]interface{}{
					"name": "frontend",
				},
			},
			{
				"kind":       "Deployment",
				"apiVersion": "apps/v1",
				"metadata": map[string]interface{}{
					"name": "frontend",
				},
			},
			{
				"kind":       "Namespace",
				"apiVersion": "v1",
				"metadata": map[string]interface{}{
					"name": "app",
				},
			},
		},
	}
}

// testDataArray is an array of (deeply nested) dicts that should be fully
// flattened
func testDataArray() testData {
	return testData{
		deep: append([]map[string]interface{}{
			testDataDeep().deep.(map[string]interface{}),
		}, testDataFlat().deep.(map[string]interface{})),

		flat: append(testDataDeep().flat.([]map[string]interface{}), testDataFlat().flat.([]map[string]interface{})...),
	}
}
