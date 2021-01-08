// Package tanka allows to use most of Tanka's features available on the
// command line programmatically as a Golang library. Keep in mind that the API
// is still experimental and may change without and signs of warnings while
// Tanka is still in alpha. Nevertheless, we try to avoid breaking changes.
package tanka

import (
	"fmt"

	"github.com/Masterminds/semver"

	"github.com/grafana/tanka/pkg/jsonnet"
	"github.com/grafana/tanka/pkg/process"
)

type JsonnetOpts = jsonnet.Opts

// Opts specify general, optional properties that apply to all actions
type Opts struct {
	JsonnetOpts

	// Filters are used to optionally select a subset of the resources
	Filters process.Matchers
}

// DEFAULT_DEV_VERSION is the placeholder version used when no actual semver is
// provided using ldflags
const DEFAULT_DEV_VERSION = "dev"

// CURRENT_VERSION is the current version of the running Tanka code
var CURRENT_VERSION = DEFAULT_DEV_VERSION

func checkVersion(constraint string) error {
	if constraint == "" {
		return nil
	}
	if CURRENT_VERSION == DEFAULT_DEV_VERSION {
		return nil
	}

	c, err := semver.NewConstraint(constraint)
	if err != nil {
		return fmt.Errorf("Parsing version constraint: '%w'. Please check 'spec.expectVersions.tanka'", err)
	}

	v, err := semver.NewVersion(CURRENT_VERSION)
	if err != nil {
		return fmt.Errorf("'%s' is not a valid semantic version: '%w'.\nThis likely means your build of Tanka is broken, as this is a compile-time value. When in doubt, please raise an issue", CURRENT_VERSION, err)
	}

	if !c.Check(v) {
		return fmt.Errorf("Current version '%s' does not satisfy the version required by the environment: '%s'. You likely need to use another version of Tanka", CURRENT_VERSION, constraint)
	}

	return nil
}
