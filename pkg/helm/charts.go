package helm

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/rs/zerolog/log"
	"sigs.k8s.io/yaml"
)

var (
	// https://regex101.com/r/7xFFtU/4
	chartExp = regexp.MustCompile(`^(?P<chart>\w+\/.+)@(?P<version>[^:\n\s]+)(?:\:(?P<path>[\w-. /\\]+))?$`)
	repoExp  = regexp.MustCompile(`^\w+$`)
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
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
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

	// Check that there are no output conflicts before vendoring
	if err := c.Manifest.Requires.Validate(); err != nil {
		return err
	}

	expectedDirs := make(map[string]bool)
	expectedDirs[c.Manifest.Directory] = true

	repositoriesUpdated := false
	log.Info().Msg("Vendoring...")
	for _, r := range c.Manifest.Requires {
		chartSubDir := parseReqName(r.Chart)
		if r.Directory != "" {
			chartSubDir = r.Directory
		}
		chartPath := filepath.Join(dir, chartSubDir)
		chartManifestPath := filepath.Join(chartPath, "Chart.yaml")
		for _, subDir := range strings.Split(chartSubDir, string(os.PathSeparator)) {
			expectedDirs[subDir] = true
		}
		expectedDirs[chartSubDir] = true

		chartDirExists, chartManifestExists := false, false
		if _, err := os.Stat(chartPath); err == nil {
			chartDirExists = true
			if _, err := os.Stat(chartManifestPath); err == nil {
				chartManifestExists = true
			} else if !os.IsNotExist(err) {
				return err
			}
		} else if !os.IsNotExist(err) {
			return err
		}

		if chartManifestExists {
			chartManifestBytes, err := os.ReadFile(chartManifestPath)
			if err != nil {
				return fmt.Errorf("reading chart manifest: %w", err)
			}
			var chartYAML chartManifest
			if err := yaml.Unmarshal(chartManifestBytes, &chartYAML); err != nil {
				return fmt.Errorf("unmarshalling chart manifest: %w", err)
			}

			if chartYAML.Version == r.Version {
				log.Info().Msgf("%s exists", r)
				continue
			}

			log.Info().Msgf("Removing %s", r)
			if err := os.RemoveAll(chartPath); err != nil {
				return err
			}
		} else if chartDirExists {
			// If the chart dir exists but the manifest doesn't, we'll clear it out and re-download the chart
			log.Info().Msgf("Removing %s", r)
			if err := os.RemoveAll(chartPath); err != nil {
				return err
			}
		}

		if !repositoriesUpdated {
			log.Info().Msg("Syncing Repositories ...")
			if err := c.Helm.RepoUpdate(Opts{Repositories: c.Manifest.Repositories}); err != nil {
				return err
			}
			repositoriesUpdated = true
		}
		log.Info().Msg("Pulling Charts ...")
		if repoName := parseReqRepo(r.Chart); !c.Manifest.Repositories.HasName(repoName) {
			return fmt.Errorf("repository %q not found for chart %q", repoName, r.Chart)
		}
		err := c.Helm.Pull(r.Chart, r.Version, PullOpts{
			Destination:      dir,
			ExtractDirectory: r.Directory,
			Opts:             Opts{Repositories: c.Manifest.Repositories},
		})
		if err != nil {
			return err
		}

		log.Info().Msgf("%s@%s downloaded", r.Chart, r.Version)
	}

	if prune {
		// Walk the charts directory looking for any unexpected directories or files to remove
		// Skips walking an expected directory that contains a Chart.yaml
		isChartFile := func(element fs.DirEntry) bool { return element.Name() == "Chart.yaml" }
		projectRootFs := os.DirFS(c.projectRoot)

		err := fs.WalkDir(projectRootFs, c.Manifest.Directory, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("error during prune: at path %s: %w", path, err)
			}
			itemType := "file"
			if !expectedDirs[d.Name()] {
				if d.IsDir() {
					itemType = "directory"
				}
				log.Info().Msgf("Pruning %s: %s", itemType, path)
				if localErr := os.RemoveAll(filepath.Join(c.projectRoot, path)); localErr != nil {
					return localErr
				}
				// If we just pruned a directory, the walk needs to skip it.
				if d.IsDir() {
					return filepath.SkipDir
				}
			} else {
				items, localErr := fs.ReadDir(projectRootFs, path)
				if localErr != nil {
					return fmt.Errorf("error listing content of dir %s: %w", path, localErr)
				}
				if slices.ContainsFunc(items, isChartFile) {
					return filepath.SkipDir
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// Add adds every Chart in reqs to the Manifest after validation, and runs
// Vendor afterwards
func (c *Charts) Add(reqs []string) error {
	log.Info().Msgf("Adding %v Charts ...", len(reqs))

	// parse new charts, append in memory
	requirements := c.Manifest.Requires
	for _, s := range reqs {
		r, err := parseReq(s)
		if err != nil {
			skip(s, err)
			continue
		}

		if requirements.Has(*r) {
			skip(s, fmt.Errorf("already exists"))
			continue
		}

		requirements = append(requirements, *r)
		log.Info().Msgf("OK: %s", s)
	}

	if err := requirements.Validate(); err != nil {
		return err
	}

	// write out
	added := len(requirements) - len(c.Manifest.Requires)
	c.Manifest.Requires = requirements
	if err := write(c.Manifest, c.ManifestFile()); err != nil {
		return err
	}

	// skipped some? fail then
	if added != len(reqs) {
		return fmt.Errorf("%v Chart(s) were skipped. Please check above logs for details", len(reqs)-added)
	}

	// worked fine? vendor it
	log.Info().Msgf("Added %v Charts to helmfile.yaml. Vendoring ...", added)
	return c.Vendor(false)
}

func (c *Charts) AddRepos(repos ...Repo) error {
	added := 0
	for _, r := range repos {
		if c.Manifest.Repositories.Has(r) {
			skip(r.Name, fmt.Errorf("already exists"))
			continue
		}

		if !repoExp.MatchString(r.Name) {
			skip(r.Name, fmt.Errorf("invalid name. cannot contain any special characters"))
			continue
		}

		c.Manifest.Repositories = append(c.Manifest.Repositories, r)
		added++
		log.Info().Msgf("OK: %s", r.Name)
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

// parseReq parses a requirement from a string of the format `repo/name@version`
func parseReq(s string) (*Requirement, error) {
	matches := chartExp.FindStringSubmatch(s)
	if matches == nil {
		return nil, fmt.Errorf("not of form 'repo/chart@version(:path)' where repo contains no special characters")
	}

	chart, ver := matches[1], matches[2]

	directory := ""
	if len(matches) == 4 {
		directory = matches[3]
	}

	return &Requirement{
		Chart:     chart,
		Version:   ver,
		Directory: directory,
	}, nil
}

// parseReqRepo parses a repo from a string of the format `repo/name`
func parseReqRepo(s string) string {
	elems := strings.SplitN(s, "/", 2)
	repo := elems[0]
	return repo
}

// parseReqName parses a name from a string of the format `repo/name`
func parseReqName(s string) string {
	elems := strings.SplitN(s, "/", 2)
	if len(elems) == 1 {
		return ""
	}
	name := elems[1]
	return name
}

func skip(s string, err error) {
	log.Info().Msgf("Skipping %s: %s.", s, err)
}
