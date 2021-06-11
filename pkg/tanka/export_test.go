package tanka

import "testing"

func Test_replaceTmplText(t *testing.T) {
	type args struct {
		s   string
		old string
		new string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"text only", args{"a", "a", "b"}, "b"},
		{"action blocks", args{"{{a}}{{.}}", "a", "b"}, "{{a}}{{.}}"},
		{"mixed", args{"a{{a}}a{{a}}a", "a", "b"}, "b{{a}}b{{a}}b"},
		{"invalid template format handled as text", args{"a}}a{{a", "a", "b"}, "b}}b{{b"},
		{
			name: "keep path separator in action block",
			args: args{`{{index .metadata.labels "app.kubernetes.io/name"}}/{{.metadata.name}}`, "/", BelRune},
			want: "{{index .metadata.labels \"app.kubernetes.io/name\"}}\u0007{{.metadata.name}}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := replaceTmplText(tt.args.s, tt.args.old, tt.args.new); got != tt.want {
				t.Errorf("replaceInTmplText() = %v, want %v", got, tt.want)
			}
		})
	}
}
