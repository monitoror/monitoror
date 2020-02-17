package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/monitoror/monitoror/models"

	"github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/monitorable/travisci"
	travisModels "github.com/monitoror/monitoror/monitorable/travisci/models"
	pkgTravis "github.com/monitoror/monitoror/pkg/gotravis"

	"github.com/shuheiktgw/go-travis"
)

type (
	travisCIRepository struct {
		config *config.TravisCI

		// Interfaces for Builds route
		travisBuildsAPI pkgTravis.TravisCI
	}
)

func NewTravisCIRepository(config *config.TravisCI) travisci.Repository {
	client := travis.NewClient(config.URL, config.Token)

	// Using Github token if exist
	// TODO: Change this to use Lazy load
	if config.GithubToken != "" {
		_, _, err := client.Authentication.UsingGithubToken(context.Background(), config.GithubToken)
		if err != nil {
			panic(fmt.Sprintf("unable to connect to TravisCI Using Github Token\n. %v\n", err))
		}
	}

	return &travisCIRepository{
		config,
		client.Builds,
	}
}

// GetBuildStatus fetch build information from travis-ci
func (r *travisCIRepository) GetLastBuildStatus(owner, repository, branch string) (*travisModels.Build, error) {
	// GetConfig
	repoSlug := fmt.Sprintf("%s/%s", owner, repository)
	options := &travis.BuildsByRepoOption{
		BranchName: []string{branch},
		Limit:      1,
		Include:    []string{"build.commit"},
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(r.config.Timeout)*time.Millisecond)
	defer cancel()

	// Request
	builds, _, err := r.travisBuildsAPI.ListByRepoSlug(ctx, repoSlug, options)
	if err != nil {
		return nil, err
	}

	// No build found
	if len(builds) == 0 {
		return nil, nil
	}

	tBuild := builds[0]
	build := &travisModels.Build{
		ID:     *tBuild.Id,
		Branch: *tBuild.Branch.Name,
		Author: models.Author{
			Name:      tBuild.Commit.Author.Name,
			AvatarURL: tBuild.Commit.Author.AvatarURL,
		},
		State: *tBuild.State,
	}

	if tBuild.StartedAt != nil {
		build.StartedAt = parseDate(*tBuild.StartedAt)
	}

	if tBuild.FinishedAt != nil {
		build.FinishedAt = parseDate(*tBuild.FinishedAt)
	}

	if tBuild.Duration != nil {
		build.Duration = parseDuration(*tBuild.Duration)
	}

	return build, nil
}

func parseDate(date string) time.Time {
	t, _ := time.Parse(time.RFC3339, date)
	return t
}

func parseDuration(duration uint) time.Duration {
	d := time.Duration(duration) * time.Second
	return d
}
