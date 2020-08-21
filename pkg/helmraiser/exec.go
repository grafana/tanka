package helmraiser

import (
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const HELM_REPO_ENVVAR = "HELM_REPOSITORY_CONFIG"

type Helm struct {
	// Repos is the list of Helm Repositories this Helm can pull Charts from
	Repos Repos

	// internal fields
	reposFile string
}

func NewHelm(repos Repos) (*Helm, error) {
	h := Helm{
		Repos: repos,
	}

	tmp, err := writeTmpFile("repositories.yml", []byte(h.Repos.String()))
	if err != nil {
		return nil, errors.Wrap(err, "Writing Helm repositories.yml")
	}
	h.reposFile = tmp

	upd := h.run("repo", "update")
	if err := upd.Run(); err != nil {
		return nil, errors.Wrap(err, "Updating Helm repositories")
	}

	return &h, nil
}

func (h Helm) Close() error {
	return os.RemoveAll(h.reposFile)
}

type Repos []Repo
type Repo struct {
	Name     string `json:"name,omitempty"`
	URL      string `json:"url,omitempty"`
	CaFile   string `json:"caFile,omitempty"`
	CertFile string `json:"certFile,omitempty"`
	KeyFile  string `json:"keyFile,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func (r Repos) String() string {
	m := map[string]interface{}{
		"repositories": r,
	}

	data, err := yaml.Marshal(m)
	if err != nil {
		panic(err)
	}
	return string(data)
}

type Values map[string]interface{}

func (v Values) String() string {
	data, err := yaml.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func (h Helm) run(action string, args ...string) *exec.Cmd {
	if h.reposFile == "" {
		panic("Helm.reposFile is unset. This helmraiser.Helm was not properly constructed. Please raise an issue")
	}

	argv := []string{action}
	argv = append(argv, args...)

	cmd := helmCmd(argv...)

	env := os.Environ()
	env = append(env, HELM_REPO_ENVVAR+"="+h.reposFile)
	cmd.Env = env

	return cmd
}

func helmCmd(args ...string) *exec.Cmd {
	binary := "helm"
	if env := os.Getenv("TANKA_HELM_PATH"); env != "" {
		binary = env
	}
	return exec.Command(binary, args...)
}

func writeTmpFile(name string, contents []byte) (string, error) {
	tmp, err := ioutil.TempFile("", "helmraiser-"+name)
	if err != nil {
		return "", err
	}

	if err := ioutil.WriteFile(tmp.Name(), contents, 0644); err != nil {
		return "", err
	}

	return tmp.Name(), nil
}
