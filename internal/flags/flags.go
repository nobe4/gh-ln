/*
Package flags implement a CLI-based environment.Environment parser for local usage.
*/
package flags

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/nobe4/gh-ln/pkg/environment"
)

var ErrFlag = errors.New("invalid flag")

func Parse() (environment.Environment, error) {
	e := environment.Environment{
		RunID:    "running locally",
		OnAction: false,
		ExecURL:  "localhost",

		App: environment.App{},
	}

	//revive:disable:line-length-limit // For flags, it's ok
	flag.BoolVar(&e.Noop, "noop", false, "Execute in no-op mode")
	flag.BoolVar(&e.Debug, "debug", false, "Enable debug mode")

	flag.StringVar(&e.Token, "token", os.Getenv("GITHUB_TOKEN"), "GitHub token to use, defaults to GITHUB_TOKEN")
	flag.StringVar(&e.Server, "server", environment.DefaultServer, "GitHub server URL")
	flag.StringVar(&e.Endpoint, "endpoint", environment.DefaultEndpoint, "GitHub API endpoint")
	flag.StringVar(&e.Config, "config", environment.DefaultConfig, "Path to the config file on the specified repo")
	flag.StringVar(&e.LocalConfig, "local-config", "", "Path to the local config file")
	flag.StringVar(&e.App.ID, "app-id", os.Getenv("INPUT_APP_ID"), "GitHub App ID, defaults to INPUT_APP_ID")
	flag.StringVar(&e.App.PrivateKey, "app-private-key", os.Getenv("INPUT_APP_PRIVATE_KEY"), "GitHub App private key, defaults to INPUT_APP_PRIVATE_KEY")
	flag.StringVar(&e.App.InstallID, "app-install-id", os.Getenv("INPUT_APP_INSTALL_ID"), "GitHub App installation ID, defaults to INPUT_APP_INSTALL_ID")

	unsafeRepo := flag.String("repo", "", "GitHub repository where the config is stored")
	//revive:enable:line-length-limit

	flag.Parse()

	repo, err := environment.ParseRepo(*unsafeRepo)
	if err != nil {
		return e, fmt.Errorf("%w -repo: %w", ErrFlag, err)
	}

	e.Repo = repo

	return e, nil
}
