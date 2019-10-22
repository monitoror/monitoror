package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/monitorable/travisci"
	"github.com/monitoror/monitoror/monitorable/travisci/models"
	pkgTravis "github.com/monitoror/monitoror/pkg/gotravis"

	"github.com/shuheiktgw/go-travis"
)

type (
	travisCIRepository struct {
		config *config.TravisCI

		// Interfaces for Builds route
		travisBuildsApi pkgTravis.TravisCI
	}
)

func NewTravisCIRepository(config *config.TravisCI, githubConfig *config.Github) travisci.Repository {
	client := travis.NewClient(config.Url, config.Token)

	// Using Github token if exist
	// TODO: Change this to use Lazy load
	if githubConfig.Token != "" {
		_, _, err := client.Authentication.UsingGithubToken(context.Background(), githubConfig.Token)
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
func (r *travisCIRepository) GetLastBuildStatus(group, repository, branch string) (build *models.Build, err error) {
	// GetConfig
	repoSlug := fmt.Sprintf("%s/%s", group, repository)
	options := &travis.BuildsByRepoOption{
		BranchName: []string{branch},
		Limit:      1,
		Include:    []string{"build.commit"},
	}

	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, time.Duration(r.config.Timeout)*time.Millisecond)

	// Request
	builds, _, err := r.travisBuildsApi.ListByRepoSlug(ctx, repoSlug, options)
	if err != nil {
		return
	}

	// No build found
	if len(builds) == 0 {
		return
	}

	tBuild := builds[0]
	build = &models.Build{
		Id:     *tBuild.Id,
		Branch: *tBuild.Branch.Name,
		Author: models.Author{
			Name:      tBuild.Commit.Author.Name,
			AvatarUrl: tBuild.Commit.Author.AvatarURL,
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

	return
}

func parseDate(date string) time.Time {
	t, _ := time.Parse(time.RFC3339, date)
	return t
}

func parseDuration(duration uint) time.Duration {
	d := time.Duration(duration) * time.Second
	return d
}
