package tanka

import (
	"fmt"
	"strings"

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
		var tla []string
		for k := range opts.TLACode {
			tla = append(tla, k)
		}
		evalScript := fmt.Sprintf(`
  local main = (import '%s');
  %s
`, entrypoint, opts.EvalScript)

		if len(tla) != 0 {
			tlaJoin := strings.Join(tla, ", ")
			evalScript = fmt.Sprintf(`
function(%s)
  local main = (import '%s')(%s);
  %s
`, tlaJoin, entrypoint, tlaJoin, opts.EvalScript)
		}

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

const PatternEvalScript = "main.%s"

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

noDataEnv(main)
`

const MetaOnlyEvalScript = `
local isKube(x) = std.isObject(x)
                  && std.objectHas(x, 'apiVersion')
                  && std.objectHas(x, 'metadata')
                  && std.objectHas(x.metadata, 'name');

local isEnv(x) = isKube(x)
                 && x.apiVersion == 'tanka.dev/v1alpha1'
                 && x.kind == 'Environment';

local isList(x) = isKube(x)
                  && std.objectHas(x, 'items');

local N(o) =
  if isEnv(o) then
    o { data: N(o.data) }
  else if isList(o) then
    o { items: N(o.items) }
  else if isKube(o) then
    { apiVersion: o.apiVersion, kind: o.kind, metadata: o.metadata }
  else if std.isObject(o) then
    std.mapWithKey(function(k, v) N(v), o)
  else if std.isArray(o) then
    std.map(N, o)
  else
    o;

N(main)
`
