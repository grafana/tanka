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

func TestDiffExitCodeMapping(t *testing.T) {
	cases := []struct {
		name      string
		err       error
		expect    bool
		expectErr bool
	}{
		{name: "nilErrorNoChanges", err: nil, expect: false, expectErr: false},
		{name: "exit0NoChanges", err: &dummyExitError{exitCode: 0}, expect: false, expectErr: false},
		{name: "exit1HasChanges", err: &dummyExitError{exitCode: 1}, expect: true, expectErr: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got bool
			var err error

			if exitErr, ok := tc.err.(exitError); !ok {
				if tc.err != nil {
					err = tc.err
				} else {
					got = false
				}
			} else {
				switch exitErr.ExitCode() {
				case 0:
					got = false
				case 1:
					got = true
				default:
					err = tc.err
				}
			}

			if tc.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expect, got)
		})
	}
}
