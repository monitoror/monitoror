package repository

import (
	"context"
	"fmt"
	"time"

	pkgTravis "github.com/monitoror/monitoror/pkg/gotravis"

	"github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/monitorable/travisci"
	"github.com/monitoror/monitoror/monitorable/travisci/models"

	"github.com/jsdidierlaurent/go-travis"
)

type (
	travisCIRepository struct {
		config *config.Config

		// Interfaces for Builds route
		travisBuildsApi pkgTravis.Builds
	}
)

func NewTravisCIRepository(conf *config.Config) travisci.Repository {
	client := travis.NewClient(conf.Monitorable.TravisCI.Url, conf.Monitorable.TravisCI.Token)

	// Using Github token if exist
	if conf.Monitorable.Github.Token != "" {
		_, _, err := client.Authentication.UsingGithubToken(context.Background(), conf.Monitorable.Github.Token)
		if err != nil {
			panic(fmt.Sprintf("Unable to connect to TravisCI Using Github Token\n. %v\n", err))
		}
	}

	return &travisCIRepository{
		conf,
		client.Builds,
	}
}

//Build fetch build information from travis-ci
func (r *travisCIRepository) Build(ctx context.Context, group, repository, branch string) (build *models.Build, err error) {
	// Config
	repoSlug := fmt.Sprintf("%s/%s", group, repository)
	options := &travis.BuildsByRepoOption{
		BranchName: []string{branch},
		EventType:  []string{travis.BuildEventTypePush},
		Limit:      1,
		State: []string{
			travis.BuildStatePassed,
			travis.BuildStateFailed,
			travis.BuildStateStarted,
			travis.BuildStateReceived,
			travis.BuildStateCreated,
			travis.BuildStateErrored,
		},
		Include: "build.commit",
	}

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
		Branch: tBuild.Branch.Name,
		Author: models.Author{
			Name:      tBuild.Commit.Author.Name,
			AvatarUrl: tBuild.Commit.Author.AvatarUrl,
		},
		State:         tBuild.State,
		PreviousState: tBuild.PreviousState,
		StartedAt:     parseDate(tBuild.StartedAt),
		FinishedAt:    parseDate(tBuild.FinishedAt),
		Duration:      parseDuration(tBuild.Duration),
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
