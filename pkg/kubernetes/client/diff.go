package client

import (
	"bytes"
	"log"
	"os/exec"
	"regexp"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/grafana/tanka/pkg/kubernetes/util"
)

// DiffServerSide takes the desired state and computes the differences on the
// server, returning them in `diff(1)` format
func (k Kubectl) DiffServerSide(data manifest.List) (*string, error) {
	supportedKinds := make(map[string]bool)
	supportedResources, err := k.APIResources()

	if err != nil {
		log.Fatalf("Cannot fetch list of supported resources : %v", err)
	}

	for i := 0; i < len(supportedResources); i++ {
		res := supportedResources[i]
		supportedKinds[res.Kind] = true
	}

	noUnkRes, unknownResources := separateUnknownResources(data, supportedKinds)

	ns, err := k.Namespaces()
	if err != nil {
		return nil, err
	}

	ready, missingNamespaces := separateMissingNamespace(noUnkRes, ns)
	cmd := k.ctl("diff", "-f", "-")

	raw := bytes.Buffer{}
	cmd.Stdout = &raw

	fw := FilterWriter{filters: []*regexp.Regexp{regexp.MustCompile(`exit status \d`)}}
	cmd.Stderr = &fw

	cmd.Stdin = strings.NewReader(ready.String())

	err = cmd.Run()
	if diffErr := parseDiffErr(err, fw.buf, k.Info().ClientVersion); diffErr != nil {
		return nil, diffErr
	}

	s := raw.String()
	missing := append(missingNamespaces, unknownResources...)
	for _, r := range missing {
		d, err := util.DiffStr(util.DiffName(r), "", r.String())
		if err != nil {
			return nil, err
		}
		s += d
	}

	if s != "" {
		return &s, nil
	}

	// no diff -> nil
	return nil, nil
}

// parseDiffErr handles the exit status code of `kubectl diff`. It returns err
// when an error happened, nil otherwise.
// "Differences found (exit status 1)" is not an error.
//
// kubectl >= 1.18:
// 0: no error, no differences
// 1: differences found
// >1: error
//
// kubectl < 1.18:
// 0: no error, no differences
// 1: error OR differences found
func parseDiffErr(err error, stderr string, version *semver.Version) error {
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		// this error is not kubectl related
		return err
	}

	// internal kubectl error
	if exitErr.ExitCode() != 1 {
		return err
	}

	// before 1.18 "exit status 1" meant error as well ... so we need to check stderr
	if version.LessThan(semver.MustParse("1.18.0")) && stderr != "" {
		return err
	}

	// differences found is not an error
	return nil
}

func separateMissingNamespace(in manifest.List, exists map[string]bool) (ready, missingNamespace manifest.List) {
	for _, r := range in {
		// namespace does not exist, also ignore implicit default ("")
		if ns := r.Metadata().Namespace(); ns != "" && !exists[ns] {
			missingNamespace = append(missingNamespace, r)
			continue
		}
		ready = append(ready, r)
	}
	return
}

func separateUnknownResources(in manifest.List, kinds map[string]bool) (ready, unknown manifest.List) {
	// Note: Matching kind only is the simplest thing that can be done. A correct solution would
	// be much more complex, since:
	// - this would need to take care of api groups and their versions for all supported kinds
	// - even if a particular Kind version is not supported, it might be possible to convert the resource
	//
	// Since these add a lot of complexity for a little benefit, the code currently implements only the
	// straightforward check, risking that API versions won't cause too much trouble.
	for i := 0; i < len(in); i++ {
		kind := in[i].Kind()

		// Validate whether `Kind` is known by the server
		_, ok := kinds[kind]
		if !ok {
			unknown = append(unknown, in[i])
		} else {
			ready = append(ready, in[i])
		}

	}

	return
}
