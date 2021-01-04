package tanka

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/grafana/tanka/pkg/jsonnet"
	"github.com/grafana/tanka/pkg/jsonnet/jpath"
)

// EvalJsonnet evaluates the jsonnet environment at the given file system path
func EvalJsonnet(path string, opts jsonnet.Opts) (raw string, err error) {
	entrypoint, err := jpath.Entrypoint(path)
	if err != nil {
		return "", err
	}

	// evaluate Jsonnet
	if opts.EvalScript != "" {
		evalScript := fmt.Sprintf(opts.EvalScript, entrypoint)
		raw, err = jsonnet.Evaluate(path, evalScript, opts)
		if err != nil {
			return "", errors.Wrap(err, "evaluating jsonnet")
		}
		return raw, nil
	}

	raw, err = jsonnet.EvaluateFile(entrypoint, opts)
	if err != nil {
		return "", errors.Wrap(err, "evaluating jsonnet")
	}
	return raw, nil
}

// EnvsOnlyEvalScript finds the Environment object (without its .data object)
const EnvsOnlyEvalScript = `
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
