package telemetry

import "testing"

func TestHasOTELConfig(t *testing.T) {
	tests := []struct {
		name string
		env  []string
		want bool
	}{
		{
			name: "no OTEL vars",
			env:  []string{"HOME=/home/user", "PATH=/usr/bin"},
			want: false,
		},
		{
			name: "OTLP endpoint set",
			env:  []string{"OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318"},
			want: true,
		},
		{
			name: "OTLP traces endpoint set",
			env:  []string{"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT=http://localhost:4318"},
			want: true,
		},
		{
			name: "TRACEPARENT set",
			env:  []string{"TRACEPARENT=00-abc-def-01"},
			want: true,
		},
		{
			name: "unrelated OTEL var should not activate tracing",
			env:  []string{"OTEL_EXPORTER_OTLP_METRICS_TEMPORALITY_PREFERENCE=delta"},
			want: false,
		},
		{
			name: "SDK disabled",
			env:  []string{"OTEL_SDK_DISABLED=true", "OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318"},
			want: false,
		},
		{
			name: "traces exporter set to none",
			env:  []string{"OTEL_TRACES_EXPORTER=none", "OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasOTELConfig(tt.env)
			if got != tt.want {
				t.Errorf("hasOTELConfig(%v) = %v, want %v", tt.env, got, tt.want)
			}
		})
	}
}
