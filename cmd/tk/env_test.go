package main

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"

	"github.com/grafana/tanka/pkg/spec/v1alpha1"
)

func TestApplyEnvSettingsInjectLabels(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want bool
	}{
		{"unrelated flag keeps injectLabels", []string{"--namespace=other"}, true},
		{"no flags keeps injectLabels", nil, true},
		{"explicit false disables injectLabels", []string{"--inject-labels=false"}, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dst := v1alpha1.New()
			dst.Spec.InjectLabels = true

			src := v1alpha1.Environment{}
			fs := pflag.NewFlagSet("env-set", pflag.ContinueOnError)
			envSettingsFlags(&src, fs)
			require.NoError(t, fs.Parse(tc.args))

			applyEnvSettings(dst, &src, fs)
			require.Equal(t, tc.want, dst.Spec.InjectLabels)
		})
	}
}

func TestApplyEnvSettingsNamespace(t *testing.T) {
	dst := v1alpha1.New()
	dst.Spec.Namespace = "default"

	src := v1alpha1.Environment{}
	fs := pflag.NewFlagSet("env-set", pflag.ContinueOnError)
	envSettingsFlags(&src, fs)
	require.NoError(t, fs.Parse([]string{"--namespace=other"}))

	applyEnvSettings(dst, &src, fs)
	require.Equal(t, "other", dst.Spec.Namespace)
}
