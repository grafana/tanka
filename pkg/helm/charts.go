package helm

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Masterminds/semver"
	"sigs.k8s.io/yaml"
)

// LoadChartfile opens a Chartfile tree
func LoadChartfile(projectRoot string) (*Charts, error) {
	// make sure project root is valid
	abs, err := filepath.Abs(projectRoot)
	if err != nil {
		return nil, err
	}

	// open chartfile
	chartfile := filepath.Join(abs, Filename)
	data, err := ioutil.ReadFile(chartfile)
	if err != nil {
		return nil, err
	}

	// parse it
	c := Chartfile{
		Version:   Version,
		Directory: DefaultDir,
	}
	if err := yaml.UnmarshalStrict(data, &c); err != nil {
		return nil, err
	}

	for i, r := range c.Requires {
		if r.Chart == "" {
			return nil, fmt.Errorf("requirements[%v]: 'chart' must be set", i)
		}
	}

	// return Charts handle
	charts := &Charts{
		Manifest:    c,
		projectRoot: abs,

		// default to ExecHelm, but allow injecting from the outside
		Helm: ExecHelm{},
	}
	return charts, nil
}

// Charts exposes the central Chartfile management functions
type Charts struct {
	// Manifest are the chartfile.yaml contents. It holds data about the developers intentions
	Manifest Chartfile

	// projectRoot is the enclosing directory of chartfile.yaml
	projectRoot string

	// Helm is the helm implementation underneath. ExecHelm is the default, but
	// any implementation of the Helm interface may be used
	Helm Helm
}

// ChartDir returns the directory pulled charts are saved in
func (c Charts) ChartDir() string {
	return filepath.Join(c.projectRoot, c.Manifest.Directory)
}

// ManifestFile returns the full path to the chartfile.yaml
func (c Charts) ManifestFile() string {
	return filepath.Join(c.projectRoot, Filename)
}

// Vendor pulls all Charts specified in the manifest into the local charts
// directory. It fetches the repository index before doing so.
func (c Charts) Vendor() error {
	dir := c.ChartDir()
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	log.Println("Syncing Repositories ...")
	if err := c.Helm.RepoUpdate(Opts{Repositories: c.Manifest.Repositories}); err != nil {
		return err
	}

	log.Println("Pulling Charts ...")
	for _, r := range c.Manifest.Requires {
		err := c.Helm.Pull(r.Chart, r.Version.String(), PullOpts{
			Destination: dir,
			Opts:        Opts{Repositories: c.Manifest.Repositories},
		})
		if err != nil {
			return err
		}

		log.Printf(" %s@%s", r.Chart, r.Version.String())
	}

	return nil
}

// Add adds every Chart in reqs to the Manifest after validation, and runs
// Vendor afterwards
func (c Charts) Add(reqs []string) error {
	log.Printf("Adding %v Charts ...", len(reqs))

	skip := func(s string, err error) {
		log.Printf(" Skipping %s: %s.", s, err)
	}

	// parse new charts, append in memory
	added := 0
	for _, s := range reqs {
		r, err := parseReq(s)
		if err != nil {
			skip(s, err)
			continue
		}

		if c.Manifest.Requires.Has(*r) {
			skip(s, fmt.Errorf("already exists"))
			continue
		}

		c.Manifest.Requires = append(c.Manifest.Requires, *r)
		added++
		log.Println(" OK:", s)
	}

	// write out
	if err := write(c.Manifest, c.ManifestFile()); err != nil {
		return err
	}

	// skipped some? fail then
	if added != len(reqs) {
		return fmt.Errorf("%v Charts were skipped. Please check above logs for details", len(reqs)-added)
	}

	// worked fine? vendor it
	log.Printf("Added %v Charts to helmfile.yaml. Vendoring ...", added)
	return c.Vendor()
}

func InitChartfile(path string) (*Charts, error) {
	c := Chartfile{
		Version: Version,
		Repositories: []Repo{{
			Name: "stable",
			URL:  "https://kubernetes-charts.storage.googleapis.com",
		}},
		Requires: make(Requirements, 0),
	}

	if err := write(c, path); err != nil {
		return nil, err
	}

	return LoadChartfile(filepath.Dir(path))
}

// write saves a Chartfile to dest
func write(c Chartfile, dest string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(dest, data, 0644)
}

var chartExp = regexp.MustCompile(`\w+\/.+@.+`)

// parseReq parses a requirement from a string of the format `repo/name@version`
func parseReq(s string) (*Requirement, error) {
	if !chartExp.MatchString(s) {
		return nil, fmt.Errorf("not of form 'repo/chart@version'")
	}

	elems := strings.Split(s, "@")
	chart := elems[0]
	ver, err := semver.NewVersion(elems[1])
	if errors.Is(err, semver.ErrInvalidSemVer) {
		return nil, fmt.Errorf("version is invalid")
	} else if err != nil {
		return nil, fmt.Errorf("version is invalid: %s", err)
	}

	return &Requirement{
		Chart:   chart,
		Version: *ver,
	}, nil
}
