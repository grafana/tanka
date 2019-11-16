package kubernetes

import "github.com/grafana/tanka/pkg/kubernetes/manifest"

// This file  contains data for testing

// testData holds data for tests
type testData struct {
	deep interface{}
	flat map[string]manifest.Manifest
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
		flat: map[string]manifest.Manifest{
			".deployment": {
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
			".service": {
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
		flat: map[string]manifest.Manifest{
			".": {
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
		flat: map[string]manifest.Manifest{
			".nginx.deployment": map[string]interface{}{
				"apiVersion": "apps/v1",
				"kind":       "Deployment",
				"metadata": map[string]interface{}{
					"name": "nginx",
				},
			},
		},
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
		flat: map[string]manifest.Manifest{
			".app.web.backend.server.nginx.deployment": {
				"kind":       "Deployment",
				"apiVersion": "apps/v1",
				"metadata": map[string]interface{}{
					"name": "nginx",
				},
			},
			".app.web.frontend.nodejs.express.service": {
				"kind":       "Service",
				"apiVersion": "v1",
				"metadata": map[string]interface{}{
					"name": "frontend",
				},
			},
			".app.web.frontend.nodejs.express.deployment": {
				"kind":       "Deployment",
				"apiVersion": "apps/v1",
				"metadata": map[string]interface{}{
					"name": "frontend",
				},
			},
			".app.namespace": {
				"kind":       "Namespace",
				"apiVersion": "v1",
				"metadata": map[string]interface{}{
					"name": "app",
				},
			},
		},
	}
}
