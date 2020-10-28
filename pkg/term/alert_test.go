package term

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfirm(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected error
	}{
		{name: "linux yes", input: "yes\n", expected: nil},
		{name: "windows yes", input: "yes\r\n", expected: nil},
		{name: "linux no", input: "no\n", expected: ErrConfirmationFailed},
		{name: "windows no", input: "no\r\n", expected: ErrConfirmationFailed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := strings.NewReader(tt.input)
			out := &strings.Builder{}

			err := confirmFrom(in, out, "foo", "yes")

			assert.Equal(t, "foo\nPlease type 'yes' to confirm: ", out.String())

			if tt.expected != nil {
				assert.EqualError(t, err, tt.expected.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
