package tkrc

import (
	"os"

	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"sigs.k8s.io/yaml"
)

type Config struct {
	AdditionalJPaths []AdditionalJPath `json:"additionalJPaths"`
}

type AdditionalJPath struct {
	RawName          string                            `json:"name"`
	RawPath          string                            `json:"path"`
	RawWeight        int                               `json:"weight"`
	MatchExpressions []metav1.LabelSelectorRequirement `json:"matchExpressions"`
}

func (jp *AdditionalJPath) Weight() int {
	return jp.RawWeight
}

func (jp *AdditionalJPath) Name() string {
	return jp.RawName
}

func (jp *AdditionalJPath) Path() string {
	return jp.RawPath
}

func (jp *AdditionalJPath) Matches(set labels.Labels) bool {
	if len(jp.MatchExpressions) == 0 {
		return true
	}
	selector := labels.NewSelector()
	for _, req := range jp.MatchExpressions {
		r, err := labels.NewRequirement(req.Key, selection.Operator(req.Operator), req.Values)
		if err != nil {
			log.Warn().Err(err).Str("rule", jp.Name()).Msg("invalid requirement, skipping whole rule")
			return false
		}
		selector = selector.Add(*r)
	}
	return selector.Matches(set)
}

func Load(path string) (*Config, error) {
	config := Config{}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := yaml.UnmarshalStrict(raw, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
