package helm

import (
	"errors"
	"fmt"
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
	data, err := os.ReadFile(chartfile)
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

// chartManifest represents a Helm chart's Chart.yaml
type chartManifest struct {
	Name    string         `yaml:"name"`
	Version semver.Version `yaml:"version"`
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
func (c Charts) Vendor(prune bool) error {
	dir := c.ChartDir()
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	expectedDirs := make(map[string]bool)

	repositoriesUpdated := false
	log.Println("Pulling Charts ...")
	for _, r := range c.Manifest.Requires {
		chartName := parseReqName(r.Chart)
		chartPath := filepath.Join(dir, chartName)
		expectedDirs[chartName] = true

		if r.Directory != "" {
			chartPath = filepath.Join(dir, r.Directory)
		}

		_, err := os.Stat(chartPath)
		if err == nil {
			chartManifestPath := filepath.Join(chartPath, "Chart.yaml")
			chartManifestBytes, err := os.ReadFile(chartManifestPath)
			if err != nil {
				return fmt.Errorf("reading chart manifest: %w", err)
			}
			var chartYAML chartManifest
			if err := yaml.Unmarshal(chartManifestBytes, &chartYAML); err != nil {
				return fmt.Errorf("unmarshalling chart manifest: %w", err)
			}

			if chartYAML.Version.String() == r.Version.String() {
				log.Printf(" %s exists", r)
				continue
			} else {
				log.Printf("Removing %s", r)
				if err := os.RemoveAll(chartPath); err != nil {
					return err
				}
			}
		} else if !os.IsNotExist(err) {
			return err
		}

		if !repositoriesUpdated {
			log.Println("Syncing Repositories ...")
			if err := c.Helm.RepoUpdate(Opts{Repositories: c.Manifest.Repositories}); err != nil {
				return err
			}
			repositoriesUpdated = true
		}
		err = c.Helm.Pull(r.Chart, r.Version.String(), PullOpts{
			Destination:      dir,
			ExtractDirectory: r.Directory,
			Opts:             Opts{Repositories: c.Manifest.Repositories},
		})
		if err != nil {
			return err
		}

		log.Printf(" %s@%s downloaded", r.Chart, r.Version.String())
	}

	if prune {
		items, err := os.ReadDir(dir)
		if err != nil {
			return fmt.Errorf("error listing the content of the charts dir: %w", err)
		}
		for _, i := range items {
			if !expectedDirs[i.Name()] {
				itemType := "file"
				if i.IsDir() {
					itemType = "directory"
				}
				log.Printf("Pruning %s: %s", itemType, i.Name())
				if err := os.RemoveAll(filepath.Join(dir, i.Name())); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// Add adds every Chart in reqs to the Manifest after validation, and runs
// Vendor afterwards
func (c *Charts) Add(reqs []string) error {
	log.Printf("Adding %v Charts ...", len(reqs))

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
		return fmt.Errorf("%v Chart(s) were skipped. Please check above logs for details", len(reqs)-added)
	}

	// worked fine? vendor it
	log.Printf("Added %v Charts to helmfile.yaml. Vendoring ...", added)
	return c.Vendor(false)
}

func (c *Charts) AddRepos(repos ...Repo) error {
	added := 0
	for _, r := range repos {
		if c.Manifest.Repositories.Has(r) {
			skip(r.Name, fmt.Errorf("already exists"))
			continue
		}

		c.Manifest.Repositories = append(c.Manifest.Repositories, r)
		added++
		log.Println(" OK:", r.Name)
	}

	// write out
	if err := write(c.Manifest, c.ManifestFile()); err != nil {
		return err
	}

	if added != len(repos) {
		return fmt.Errorf("%v Repo(s) were skipped. Please check above logs for details", len(repos)-added)
	}

	return nil
}

func InitChartfile(path string) (*Charts, error) {
	c := Chartfile{
		Version: Version,
		Repositories: []Repo{{
			Name: "stable",
			URL:  "https://charts.helm.sh/stable",
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

	return os.WriteFile(dest, data, 0644)
}

// https://regex101.com/r/VAklNg/1
var chartExp = regexp.MustCompile(`\w+\/.+@.+(\:.+)?`)

// parseReq parses a requirement from a string of the format `repo/name@version`
func parseReq(s string) (*Requirement, error) {
	if !chartExp.MatchString(s) {
		return nil, fmt.Errorf("not of form 'repo/chart@version(:path)'")
	}

	elems := strings.Split(s, ":")
	directory := ""
	if len(elems) > 1 {
		s = elems[0]
		directory = elems[1]
	}

	elems = strings.Split(s, "@")
	chart := elems[0]
	ver, err := semver.NewVersion(elems[1])
	if errors.Is(err, semver.ErrInvalidSemVer) {
		return nil, fmt.Errorf("version is invalid")
	} else if err != nil {
		return nil, fmt.Errorf("version is invalid: %s", err)
	}

	return &Requirement{
		Chart:     chart,
		Version:   *ver,
		Directory: directory,
	}, nil
}

// parseReqName parses a name from a string of the format `repo/name`
func parseReqName(s string) string {
	elems := strings.Split(s, "/")
	name := elems[1]
	return name
}

func skip(s string, err error) {
	log.Printf(" Skipping %s: %s.", s, err)
}
