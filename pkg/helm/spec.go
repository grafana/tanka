package helm

import (
	"fmt"

	"github.com/Masterminds/semver"
)

const (
	// Version of the current Chartfile implementation
	Version = 1

	// Filename of the Chartfile
	Filename = "chartfile.yaml"

	// DefaultDir is the directory used for storing Charts if not specified
	// otherwise
	DefaultDir = "charts"
)

// Chartfile is the schema used to declaratively define locally required Helm
// Charts
type Chartfile struct {
	// Version of the Chartfile schema (for future use)
	Version uint `json:"version"`

	// Repositories to source from
	Repositories Repos `json:"repositories"`

	// Requires lists Charts expected to be present in the charts folder
	Requires Requirements `json:"requires"`

	// Folder to use for storing Charts. Defaults to 'charts'
	Directory string `json:"directory,omitempty"`
}

// Repo describes a single Helm repository
type Repo struct {
	Name     string `json:"name,omitempty"`
	URL      string `json:"url,omitempty"`
	CAFile   string `json:"caFile,omitempty"`
	CertFile string `json:"certFile,omitempty"`
	KeyFile  string `json:"keyFile,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type Repos []Repo

// Has reports whether 'repo' is already part of the repositories
func (r Repos) Has(repo Repo) bool {
	for _, x := range r {
		if x == repo {
			return true
		}
	}

	return false
}

// Requirement describes a single required Helm Chart.
// Both, Chart and Version are required
type Requirement struct {
	Chart     string         `json:"chart"`
	Version   semver.Version `json:"version"`
	Directory string         `json:"directory,omitempty"`
}

func (r Requirement) String() string {
	s := r.Chart + "@" + r.Version.String()
	if r.Directory != "" {
		s += fmt.Sprintf(" (dir: %s)", r.Directory)
	}
	return s
}

// Requirements is an aggregate of all required Charts
type Requirements []Requirement

// Has reports whether 'req' is already part of the requirements
func (r Requirements) Has(req Requirement) bool {
	for _, x := range r {
		if x == req {
			return true
		}
	}

	return false
}
