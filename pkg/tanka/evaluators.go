package tanka

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/grafana/tanka/pkg/jsonnet"
	"github.com/grafana/tanka/pkg/jsonnet/jpath"
)

// Evaluator signature for implementing arbitrary Jsonnet evaluators
type Evaluator func(path string, opts jsonnet.Opts) (string, error)

// DefaultEvaluator evaluates the jsonnet environment at the given file system path
func DefaultEvaluator(path string, opts jsonnet.Opts) (string, error) {
	entrypoint, err := jpath.Entrypoint(path)
	if err != nil {
		return "", err
	}

	// evaluate Jsonnet
	var raw string
	if opts.EvalPattern != "" {
		evalScript := fmt.Sprintf("(import '%s').%s", entrypoint, opts.EvalPattern)
		raw, err = jsonnet.Evaluate(entrypoint, evalScript, opts)
		if err != nil {
			return "", errors.Wrap(err, "evaluating jsonnet")
		}
	} else {
		raw, err = jsonnet.EvaluateFile(entrypoint, opts)
		if err != nil {
			return "", errors.Wrap(err, "evaluating jsonnet")
		}
	}
	return raw, nil
}

// EnvsOnlyEvaluator finds the Environment object (without its .data object) at
// the given file system path intended for use by the `tk env` command
func EnvsOnlyEvaluator(path string, opts jsonnet.Opts) (string, error) {
	entrypoint, err := jpath.Entrypoint(path)
	if err != nil {
		return "", err
	}

	// Snippet to find all Environment objects and remove the .data object for faster evaluation
	noData := `
local noDataEnv(object) =
  if std.isObject(object)
  then
    if std.objectHas(object, 'apiVersion')
       && std.objectHas(object, 'kind')
    then
      if object.kind == 'Environment'
      then object { data:: {} }
      else {}
    else
      std.mapWithKey(
        function(key, obj)
          noDataEnv(obj),
        object
      )
  else if std.isArray(object)
  then
    std.map(
      function(obj)
        noDataEnv(obj),
      object
    )
  else {};

noDataEnv(import '%s')
`

	// evaluate Jsonnet with noData snippet
	var raw string
	evalScript := fmt.Sprintf(noData, entrypoint)
	raw, err = jsonnet.Evaluate(entrypoint, evalScript, opts)
	if err != nil {
		return "", errors.Wrap(err, "evaluating jsonnet")
	}
	return raw, nil
}
