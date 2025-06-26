package tkrc

import (
	"testing"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

func TestMatching(t *testing.T) {
	tests := map[string]struct {
		Labels           map[string]string
		Matches          bool
		MatchExpressions []metav1.LabelSelectorRequirement
	}{
		"match-due-to-no-requirements": {
			MatchExpressions: []metav1.LabelSelectorRequirement{},
			Labels:           map[string]string{"cluster_name": "test"},
			Matches:          true,
		},
		"single-req-match": {
			MatchExpressions: []metav1.LabelSelectorRequirement{
				{
					Key:      "cluster_name",
					Operator: metav1.LabelSelectorOperator(selection.In),
					Values:   []string{"test"},
				},
			},
			Labels:  map[string]string{"cluster_name": "test"},
			Matches: true,
		},
	}

	for testname, test := range tests {
		t.Run(testname, func(t *testing.T) {
			rule := AdditionalJPath{
				MatchExpressions: test.MatchExpressions,
			}
			require.Equal(t, test.Matches, rule.Matches(labels.Set(test.Labels)))
		})
	}
}
