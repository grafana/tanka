package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/grafana/tanka/pkg/spec/v1alpha1"
	"github.com/grafana/tanka/pkg/tanka"
)

type workflowFlagVars struct {
	targets []string
}

func workflowFlags(fs *pflag.FlagSet) *workflowFlagVars {
	v := workflowFlagVars{}
	fs.StringSliceVarP(&v.targets, "target", "t", nil, "only use the specified objects (Format: <type>/<name>)")
	return &v
}

func labelSelectorFlag(fs *pflag.FlagSet) func() labels.Selector {
	labelSelector := fs.StringP("selector", "l", "", "Label selector. Uses the same syntax as kubectl does")

	return func() labels.Selector {
		if *labelSelector != "" {
			selector, err := labels.Parse(*labelSelector)
			if err != nil {
				log.Fatalf("Could not parse selector (-l) %s", *labelSelector)
			}
			return selector
		}
		return nil
	}
}

func jsonnetFlags(fs *pflag.FlagSet) func() tanka.JsonnetOpts {
	getExtCode, getTLACode := cliCodeParser(fs)

	return func() tanka.JsonnetOpts {
		return tanka.JsonnetOpts{
			ExtCode: getExtCode(),
			TLACode: getTLACode(),
		}
	}
}

func cliCodeParser(fs *pflag.FlagSet) (func() map[string]string, func() map[string]string) {
	// need to use StringArray instead of StringSlice, because pflag attempts to
	// parse StringSlice using the csv parser, which breaks when passing objects
	extCode := fs.StringArray("ext-code", nil, "Set code value of extVar (Format: key=<code>)")
	extStr := fs.StringArrayP("ext-str", "V", nil, "Set string value of extVar (Format: key=value)")

	tlaCode := fs.StringArray("tla-code", nil, "Set code value of top level function (Format: key=<code>)")
	tlaStr := fs.StringArrayP("tla-str", "A", nil, "Set string value of top level function (Format: key=value)")

	newParser := func(kind string, code, str *[]string) func() map[string]string {
		return func() map[string]string {
			m := make(map[string]string)
			for _, s := range *code {
				split := strings.SplitN(s, "=", 2)
				if len(split) != 2 {
					log.Fatalf(kind+"-code argument has wrong format: `%s`. Expected `key=<code>`", s)
				}
				m[split[0]] = split[1]
			}

			for _, s := range *str {
				split := strings.SplitN(s, "=", 2)
				if len(split) != 2 {
					log.Fatalf(kind+"-str argument has wrong format: `%s`. Expected `key=<value>`", s)
				}
				m[split[0]] = fmt.Sprintf(`"%s"`, split[1])
			}
			return m
		}
	}

	return newParser("ext", extCode, extStr),
		newParser("tla", tlaCode, tlaStr)
}

func envSettingsFlags(env *v1alpha1.Environment, fs *pflag.FlagSet) {
	fs.StringVar(&env.Spec.APIServer, "server", env.Spec.APIServer, "endpoint of the Kubernetes API")
	fs.StringVar(&env.Spec.APIServer, "server-from-context", env.Spec.APIServer, "set the server to a known one from $KUBECONFIG")
	fs.StringVar(&env.Spec.Namespace, "namespace", env.Spec.Namespace, "namespace to create objects in")
	fs.StringVar(&env.Spec.DiffStrategy, "diff-strategy", env.Spec.DiffStrategy, "specify diff-strategy. Automatically detected otherwise.")
}
