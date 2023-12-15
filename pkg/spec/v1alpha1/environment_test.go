package v1alpha1

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironmentNameLabel(t *testing.T) {
	type testCase struct {
		name                 string
		inputEnvironment     *Environment
		expectedLabelPreHash string
		expectError          bool
	}

	testCases := []testCase{
		{
			name: "Default environment label hash",
			inputEnvironment: &Environment{
				Spec: Spec{
					Namespace: "default",
				},
				Metadata: Metadata{
					Name:      "environments/a-nice-go-test",
					Namespace: "main.jsonnet",
				},
			},
			expectedLabelPreHash: "environments/a-nice-go-test:main.jsonnet",
		},
		{
			name: "Overriden single nested field",
			inputEnvironment: &Environment{
				Spec: Spec{
					Namespace: "default",
					TankaEnvLabelFromFields: []string{
						".metadata.name",
					},
				},
				Metadata: Metadata{
					Name: "environments/another-nice-go-test",
				},
			},
			expectedLabelPreHash: "environments/another-nice-go-test",
		},
		{
			name: "Overriden multiple nested field",
			inputEnvironment: &Environment{
				Spec: Spec{
					Namespace: "default",
					TankaEnvLabelFromFields: []string{
						".metadata.name",
						".spec.namespace",
					},
				},
				Metadata: Metadata{
					Name: "environments/another-nice-go-test",
				},
			},
			expectedLabelPreHash: "environments/another-nice-go-test:default",
		},
		{
			name: "Override field of map type",
			inputEnvironment: &Environment{
				Spec: Spec{
					TankaEnvLabelFromFields: []string{
						".metadata.labels.project",
					},
				},
				Metadata: Metadata{
					Name: "environments/another-nice-go-test",
					Labels: map[string]string{
						"project": "an-equally-nice-project",
					},
				},
			},
			expectedLabelPreHash: "an-equally-nice-project",
		},
		{
			name: "Label value not primitive type",
			inputEnvironment: &Environment{
				Spec: Spec{
					TankaEnvLabelFromFields: []string{
						".metadata",
					},
				},
				Metadata: Metadata{
					Name: "environments/another-nice-go-test",
				},
			},
			expectError: true,
		},
		{
			name: "Attempted descent past non-object like type",
			inputEnvironment: &Environment{
				Spec: Spec{
					TankaEnvLabelFromFields: []string{
						".metadata.name.nonExistent",
					},
				},
				Metadata: Metadata{
					Name: "environments/not-an-object",
				},
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expectedLabelHashParts := sha256.Sum256([]byte(tc.expectedLabelPreHash))
			expectedLabelHashChars := []rune(hex.EncodeToString(expectedLabelHashParts[:]))
			expectedLabelHash := string(expectedLabelHashChars[:48])
			actualLabelHash, err := tc.inputEnvironment.NameLabel()

			if tc.expectedLabelPreHash != "" {
				assert.Equal(t, expectedLabelHash, actualLabelHash)
			} else {
				assert.Equal(t, "", actualLabelHash)
			}

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
