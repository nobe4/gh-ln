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
	e := environment.New()

	//revive:disable:line-length-limit // For flags, it's ok
	noop := flag.Bool("noop", false, "Execute in no-op mode")
	token := flag.String("token", os.Getenv("GITHUB_TOKEN"), "GitHub token to use, defaults to GITHUB_TOKEN")
	repo := flag.String("repo", "", "GitHub repository where the config is stored")
	server := flag.String("server", environment.DefaultServer, "GitHub server URL")
	endpoint := flag.String("endpoint", environment.DefaultEndpoint, "GitHub API endpoint")
	config := flag.String("config", environment.DefaultConfig, "Path to the config file on the specified repo")
	localConfig := flag.String("local-config", "", "Path to the local config file")
	appID := flag.String("app-id", os.Getenv("INPUT_APP_ID"), "GitHub App ID, defaults to INPUT_APP_ID")
	appPrivateKey := flag.String("app-private-key", os.Getenv("INPUT_APP_PRIVATE_KEY"), "GitHub App private key, defaults to INPUT_APP_PRIVATE_KEY")
	appInstallID := flag.String("app-install-id", os.Getenv("INPUT_APP_INSTALL_ID"), "GitHub App installation ID, defaults to INPUT_APP_INSTALL_ID")
	debug := flag.Bool("debug", false, "Enable debug mode")
	//revive:enable:line-length-limit

	flag.Parse()

	e.Noop = *noop
	e.Token = *token

	r, err := environment.ParseRepo(*repo)
	if err != nil {
		return e, fmt.Errorf("%w -repo: %w", ErrFlag, err)
	}

	e.Repo = r

	e.Server = *server
	e.Endpoint = *endpoint
	e.RunID = "running locally"
	e.Config = *config
	e.LocalConfig = *localConfig
	e.App = environment.App{
		ID:         *appID,
		PrivateKey: *appPrivateKey,
		InstallID:  *appInstallID,
	}
	e.OnAction = false
	e.ExecURL = "localhost"
	e.Debug = *debug

	return e, nil
}
