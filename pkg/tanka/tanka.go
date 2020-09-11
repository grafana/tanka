// Package tanka allows to use most of Tanka's features available on the
// command line programmatically as a Golang library. Keep in mind that the API
// is still experimental and may change without and signs of warnings while
// Tanka is still in alpha. Nevertheless, we try to avoid breaking changes.
package tanka

import (
	"github.com/grafana/tanka/pkg/jsonnet"
	"github.com/grafana/tanka/pkg/process"
)

// Opts specify general, optional properties that apply to all actions
type Opts struct {
	JsonnetOpts

	// Filters are used to optionally select a subset of the resources
	Filters process.Matchers
}

type JsonnetOpts = jsonnet.Opts
