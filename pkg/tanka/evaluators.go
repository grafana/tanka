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
func DefaultEvaluator(path string, opts jsonnet.Opts) (raw string, err error) {
	// evaluate Jsonnet
	if opts.EvalScript != "" {
		raw, err = EvalscriptEvaluator(path, opts)
		if err != nil {
			return "", errors.Wrap(err, "evaluating jsonnet")
		}
	} else {
		entrypoint, err := jpath.Entrypoint(path)
		if err != nil {
			return "", err
		}

		raw, err = jsonnet.EvaluateFile(entrypoint, opts)
		if err != nil {
			return "", errors.Wrap(err, "evaluating jsonnet")
		}
	}
	return raw, nil
}

// EvalscriptEvaluator finds the Environment object (without its .data object) at
// the given file system path intended for use by the `tk env` command
func EvalscriptEvaluator(path string, opts jsonnet.Opts) (string, error) {
	entrypoint, err := jpath.Entrypoint(path)
	if err != nil {
		return "", err
	}

	// evaluate Jsonnet with noData snippet
	var raw string
	evalScript := fmt.Sprintf(opts.EvalScript, entrypoint)
	raw, err = jsonnet.Evaluate(entrypoint, evalScript, opts)
	if err != nil {
		return "", errors.Wrap(err, "evaluating jsonnet")
	}
	return raw, nil
}

// EnvsOnlyEvaluator finds the Environment object (without its .data object) at
// the given file system path intended for use by the `tk env` command
func EnvsOnlyEvaluator(path string, opts jsonnet.Opts) (string, error) {
	// Snippet to find all Environment objects and remove the .data object for faster evaluation
	opts.EvalScript = `
local noDataEnv(object) =
  std.prune(
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
    else {}
  );


noDataEnv(import '%s')
`

	return EvalscriptEvaluator(path, opts)
}

// SingleEnvEvaluator finds the Environment object (without its .data object) at
// the given file system path intended for use by the `tk env` command
func SingleEnvEvaluator(path string, opts jsonnet.Opts) (string, error) {
	entrypoint, err := jpath.Entrypoint(path)
	if err != nil {
		return "", err
	}

	// Snippet to find all Environment objects and remove the .data object for faster evaluation
	singleEnv := `
local SingleEnv(object) =
  std.prune(
    if std.isObject(object)
    then
      if std.objectHas(object, 'apiVersion')
         && std.objectHas(object, 'kind')
      then
        if object.kind == 'Environment'
           && object.metadata.name = '%s'
        then object
        else {}
      else
        std.mapWithKey(
          function(key, obj)
            SingleEnv(obj),
          object
        )
    else if std.isArray(object)
    then
      std.map(
        function(obj)
          SingleEnv(obj),
        object
      )
    else {}
  );


SingleEnv(import '%s')
`

	// evaluate Jsonnet with singleEnv snippet
	var raw string
	evalScript := fmt.Sprintf(singleEnv, opts.EvalPattern, entrypoint)
	raw, err = jsonnet.Evaluate(entrypoint, evalScript, opts)
	if err != nil {
		return "", errors.Wrap(err, "evaluating jsonnet")
	}
	return raw, nil
}
