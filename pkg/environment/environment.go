/*
Package environment implements helpers to get inputs and environment from the
GitHub action's environment variables.
Called `environment` to avoid conflict with the `context` package.

https://docs.github.com/en/actions/writing-workflows/choosing-what-your-workflow-does/accessing-contextual-information-about-workflow-runs
*/
package environment

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/nobe4/gh-ln/pkg/github"
)

var (
	ErrInvalidEnvironment = errors.New("error parsing environment")
	ErrNoToken            = errors.New("github token not found")
	ErrNoRepo             = errors.New("github repository not found")
	ErrInvalidRepo        = errors.New("github repository invalid: want owner/repo")
)

const (
	DefaultEndpoint = "https://api.github.com"
	DefaultServer   = "https://github.com"
	DefaultConfig   = ".ln-config.yaml"
	DefaultRunID    = ""
	Redacted        = "[redacted]"
	Missing         = "[missing]"
)

type App struct {
	ID         string `json:"app_id"`          // INPUT_APP_ID
	PrivateKey string `json:"app_private_key"` // INPUT_APP_PRIVATE_KEY
	InstallID  string `json:"app_install_id"`  // INPUT_APP_INSTALL_ID
}

type Environment struct {
	Noop        bool        `json:"noop"`         // INPUT_NOOP
	Token       string      `json:"token"`        // GITHUB_TOKEN / INPUT_TOKEN
	App         App         `json:"app"`          // For Github-App authentication
	Repo        github.Repo `json:"repo"`         // GITHUB_REPOSITORY
	Server      string      `json:"server"`       // GITHUB_SERVER_URL
	Endpoint    string      `json:"endpoint"`     // GITHUB_API_URL
	RunID       string      `json:"run_id"`       // GITHUB_RUN_ID
	Config      string      `json:"config"`       // INPUT_CONFIG
	LocalConfig string      `json:"local_config"` // Read config from the filesystem.
	OnAction    bool        `json:"on_action"`
	ExecURL     string      `json:"exec_url"`
	Debug       bool        `json:"debug"` // RUNNER_DEBUG
}

//nolint:revive // No, I don't want to leak secrets.
func (e Environment) String() string {
	e.Token = missingOrRedacted(e.Token)
	e.App.ID = missingOrRedacted(e.App.ID)
	e.App.PrivateKey = missingOrRedacted(e.App.PrivateKey)
	e.App.InstallID = missingOrRedacted(e.App.InstallID)

	out, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return err.Error()
	}

	return string(out)
}

// TODO: this should be in github.Repo.Parse
func ParseRepo(repoName string) (github.Repo, error) {
	repo := github.Repo{}

	if repoName == "" {
		return repo, ErrNoRepo
	}

	var found bool
	repo.Owner.Login, repo.Repo, found = strings.Cut(repoName, "/")

	if !found {
		return repo, ErrInvalidRepo
	}

	return repo, nil
}

func missingOrRedacted(s string) string {
	if s == "" {
		return Missing
	}

	return Redacted
}
