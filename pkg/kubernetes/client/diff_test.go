package client

import (
	"fmt"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/stretchr/testify/require"
)

func TestParseDiffError(t *testing.T) {
	tests := map[string]struct {
		err         error
		stderr      string
		version     *semver.Version
		expectedErr error
	}{
		// If this is not an exit error, then we just pass it through:
		"no-exiterr": {
			err:         fmt.Errorf("something else"),
			stderr:      "error-details",
			version:     semver.MustParse("1.17.0"),
			expectedErr: fmt.Errorf("something else"),
		},
		// If kubectl returns with an exit code other than 1 its an indicator
		// that it is not an internal error and so we return it as is:
		"return-internal-as-is": {
			err: &dummyExitError{
				exitCode: 123,
			},
			stderr:      "error-details",
			version:     semver.MustParse("1.17.0"),
			expectedErr: fmt.Errorf("ExitError"),
		},
		// If kubectl is is < 1.18.0, then the error should contain then stderr
		// content:
		"lt-1.18.0-contains-stderr": {
			err: &dummyExitError{
				exitCode: 1,
			},
			stderr:      "error-details",
			version:     semver.MustParse("1.17.0"),
			expectedErr: fmt.Errorf("diff failed: ExitError (error-details)"),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			err := parseDiffErr(test.err, test.stderr, test.version)
			if test.expectedErr == nil {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			require.Equal(t, test.expectedErr.Error(), err.Error())
		})
	}
}

type dummyExitError struct {
	exitCode int
}

func (e *dummyExitError) Error() string {
	return "ExitError"
}

func (e *dummyExitError) ExitCode() int {
	return e.exitCode
}
