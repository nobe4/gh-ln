package ln

import (
	"context"
	"fmt"

	"github.com/nobe4/gh-ln/internal/config"
	"github.com/nobe4/gh-ln/internal/format"
	"github.com/nobe4/gh-ln/pkg/github"
	"github.com/nobe4/gh-ln/pkg/log"
)

const (
	headName         = "auto-action-ln"
	pullTitle        = "auto(ln): update links"
	pullBodyTemplate = `
{{/* This defines a backtick character to use in the markdown. */}}
{{- $b := "` + "`" + `" -}}
This automated PR updates the following files:

| From | To  | Status |
| ---  | --- | ---    |
{{ range .Data -}}
| [{{ $b }}{{ .From }}{{ $b }}]({{ .From.HTMLURL }}) | {{ $b }}{{ .To.Path }}{{ $b }} | {{ .Status }} |
{{ end }}

---

| Quick links | [execution]({{ .Environment.ExecURL }}) | [configuration]({{ .Environment.Server }}{{ .Config.Source.HTMLPath }}) | [gh-ln](https://github.com/nobe4/gh-ln) |
| --- | --- | --- | --- |
`
)

func processGroups(ctx context.Context, g *github.GitHub, f format.Formatter, groups config.Groups) error {
	for _, l := range groups {
		err := processLinks(ctx, g, f, l)
		if err != nil {
			return err
		}
	}

	return nil
}

func processLinks(ctx context.Context, g *github.GitHub, f format.Formatter, l config.Links) error {
	toRepo := l[0].To.Repo

	log.Group("Processing links for " + toRepo.String())
	defer log.GroupEnd()

	base, head, err := g.GetBaseAndHeadBranches(ctx, toRepo, headName)
	if err != nil {
		return fmt.Errorf("failed to prepare branches: %w", err)
	}

	log.Debug("Parsed branches", "head", head, "base", base)

	updated := l.Update(ctx, g, f, head)
	if !updated && head.New {
		log.Info("No link was updated, cleaning up.", "repo", toRepo, "branch", head.Name)

		err = g.DeleteBranch(ctx, toRepo, head.Name)
		if err != nil {
			return fmt.Errorf("failed to delete non-updated branch: %w", err)
		}

		return nil
	}

	pullBody, err := f.Format(pullBodyTemplate, l)
	if err != nil {
		return fmt.Errorf("failed to create pull request body: %w", err)
	}

	log.Debug("Pull body", "body", pullBody)

	pull, err := g.GetOrCreatePull(ctx, toRepo, base.Name, head.Name, pullTitle, pullBody)
	if err != nil {
		return fmt.Errorf("failed to get pull request: %w", err)
	}

	log.Info("Result pull request", "pull", pull, "new", pull.New)

	return nil
}
