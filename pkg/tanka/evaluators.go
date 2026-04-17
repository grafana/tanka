package tanka

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/grafana/tanka/pkg/jsonnet"
	"github.com/grafana/tanka/pkg/jsonnet/implementations/types"
	"github.com/grafana/tanka/pkg/jsonnet/jpath"
)

// buildEvalScript constructs the jsonnet snippet that wraps the user's
// EvalScript around an `import` of the entrypoint. The entrypoint is
// normalised to forward slashes so Windows paths (e.g. `C:\proj\main.jsonnet`)
// do not produce invalid jsonnet escape sequences when embedded in a
// single-quoted string literal. See grafana/tanka#551.
func buildEvalScript(entrypoint, evalScript string, tlas []string, isFunction bool) string {
	entrypoint = filepath.ToSlash(entrypoint)
	if isFunction {
		tlaParams := strings.Join(tlas, ", ")
		tlaArgs := make([]string, len(tlas))
		for i, k := range tlas {
			tlaArgs[i] = k + "=" + k
		}
		tlaArgsJoin := strings.Join(tlaArgs, ", ")
		return fmt.Sprintf(`
function(%s)
  local main = (import '%s')(%s);
  %s
`, tlaParams, entrypoint, tlaArgsJoin, evalScript)
	}
	return fmt.Sprintf(`
  local main = (import '%s');
  %s
`, entrypoint, evalScript)
}

// EvalJsonnet evaluates the jsonnet environment at the given file system path
func evalJsonnet(ctx context.Context, path string, impl types.JsonnetImplementation, opts jsonnet.Opts) (raw string, err error) {
	entrypoint, err := jpath.Entrypoint(path)
	if err != nil {
		return "", err
	}

	// evaluate Jsonnet
	if opts.EvalScript != "" {
		// Determine if the entrypoint is a function.
		isFunctionProbe := fmt.Sprintf("std.isFunction(import '%s')", filepath.ToSlash(entrypoint))
		isFunction, err := jsonnet.Evaluate(ctx, path, impl, isFunctionProbe, opts)
		if err != nil {
			return "", fmt.Errorf("evaluating jsonnet in path '%s': %w", path, err)
		}
		var tlas []string
		for k := range opts.TLACode {
			tlas = append(tlas, k)
		}
		evalScript := buildEvalScript(entrypoint, opts.EvalScript, tlas, isFunction == "true\n")

		raw, err = jsonnet.Evaluate(ctx, path, impl, evalScript, opts)
		if err != nil {
			return "", fmt.Errorf("evaluating jsonnet in path '%s': %w", path, err)
		}
		return raw, nil
	}

	raw, err = jsonnet.EvaluateFile(ctx, impl, entrypoint, opts)
	if err != nil {
		return "", errors.Wrap(err, "evaluating jsonnet")
	}
	return raw, nil
}

func PatternEvalScript(expr string) string {
	if strings.HasPrefix(expr, "[") {
		return fmt.Sprintf("main%s", expr)
	}
	return fmt.Sprintf("main.%s", expr)
}

// MetadataEvalScript finds the Environment object (without its .data object)
const MetadataEvalScript = `
local noDataEnv(object) =
  std.prune(
    if std.isObject(object)
    then
      if std.objectHas(object, 'apiVersion')
         && std.objectHas(object, 'kind')
      then
        if object.kind == 'Environment'
        then object { data+:: {} }
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

noDataEnv(main)
`

// MetadataSingleEnvEvalScript returns a Single Environment object
const MetadataSingleEnvEvalScript = `
local singleEnv(object) =
  std.prune(
    if std.isObject(object)
    then
      if std.objectHas(object, 'apiVersion')
         && std.objectHas(object, 'kind')
      then
        if object.kind == 'Environment'
        && std.member(object.metadata.name, '%s')
        then object { data:: super.data }
        else {}
      else
        std.mapWithKey(
          function(key, obj)
            singleEnv(obj),
          object
        )
    else if std.isArray(object)
    then
      std.map(
        function(obj)
          singleEnv(obj),
        object
      )
    else {}
  );

singleEnv(main)
`

// SingleEnvEvalScript returns a Single Environment object
const SingleEnvEvalScript = `
local singleEnv(object) =
  if std.isObject(object)
  then
    if std.objectHas(object, 'apiVersion')
       && std.objectHas(object, 'kind')
    then
      if object.kind == 'Environment'
      && std.member(object.metadata.name, '%s')
      then object
      else {}
    else
      std.mapWithKey(
        function(key, obj)
          singleEnv(obj),
        object
      )
  else if std.isArray(object)
  then
    std.map(
      function(obj)
        singleEnv(obj),
      object
    )
  else {};

singleEnv(main)
`
