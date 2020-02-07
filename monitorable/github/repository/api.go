package repository

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/monitoror/monitoror/pkg/monitoror/utils/gravatar"

	"golang.org/x/oauth2"

	"github.com/sourcegraph/httpcache"

	"github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/monitorable/github"
	"github.com/monitoror/monitoror/monitorable/github/models"
	"github.com/monitoror/monitoror/pkg/gogithub"

	githubApi "github.com/google/go-github/github"
)

type (
	githubRepository struct {
		searchService       gogithub.SearchService
		checksService       gogithub.ChecksService
		repositoriesService gogithub.RepositoriesService
		pullRequestService  gogithub.PullRequestService
		gitService          gogithub.GitService

		config *config.Github
	}
)

func NewGithubRepository(config *config.Github) github.Repository {
	httpClient := &http.Client{
		Transport: &oauth2.Transport{
			// Use NewMemoryCacheTransport to save github rate limit
			Base:   httpcache.NewMemoryCacheTransport(),
			Source: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: config.Token}),
		},
		Timeout: time.Duration(config.Timeout) * time.Millisecond,
	}

	// Init Github Client
	client := githubApi.NewClient(httpClient)

	return &githubRepository{
		searchService:       client.Search,
		checksService:       client.Checks,
		repositoriesService: client.Repositories,
		pullRequestService:  client.PullRequests,
		gitService:          client.Git,
		config:              config,
	}
}

func (gr *githubRepository) GetIssuesCount(query string) (int, error) {
	issuesResult, _, err := gr.searchService.Issues(context.TODO(), query, nil)
	if err != nil {
		return 0, err
	}

	return issuesResult.GetTotal(), err
}

func (gr *githubRepository) GetChecks(owner, repository, ref string) (*models.Checks, error) {
	checks := &models.Checks{Runs: []models.Run{}, Statuses: []models.Status{}}

	checkRuns, _, err := gr.checksService.ListCheckRunsForRef(context.TODO(), owner, repository, ref, nil)
	if err != nil {
		return nil, err
	}

	for _, checkRun := range checkRuns.CheckRuns {
		if checkRun.Name != nil && checkRun.Status != nil {
			run := models.Run{
				ID:         checkRun.GetID(),
				Title:      checkRun.GetName(),
				Status:     checkRun.GetStatus(),
				Conclusion: checkRun.GetConclusion(),
			}

			if checkRun.StartedAt != nil {
				run.StartedAt = &checkRun.StartedAt.Time
			}

			if checkRun.CompletedAt != nil {
				run.CompletedAt = &checkRun.CompletedAt.Time
			}

			checks.HeadCommit = checkRun.HeadSHA
			checks.Runs = append(checks.Runs, run)
		}
	}

	repoStatuses, _, err := gr.repositoriesService.ListStatuses(context.TODO(), owner, repository, ref, nil)
	if err != nil {
		return nil, err
	}

	for _, repoStatus := range repoStatuses {
		if repoStatus.Context != nil && repoStatus.State != nil && repoStatus.CreatedAt != nil && repoStatus.UpdatedAt != nil {
			status := models.Status{
				ID:        repoStatus.GetID(),
				Title:     repoStatus.GetContext(),
				State:     repoStatus.GetState(),
				CreatedAt: repoStatus.GetCreatedAt(),
				UpdatedAt: repoStatus.GetUpdatedAt(),
			}

			if repoStatus.URL != nil {
				headCommit := repoStatus.GetURL()
				headCommit = headCommit[strings.LastIndex(headCommit, "/")+1:]

				checks.HeadCommit = &headCommit
			}

			checks.Statuses = append(checks.Statuses, status)
		}
	}

	return checks, nil
}

func (gr *githubRepository) GetPullRequests(owner, repository string) ([]models.PullRequest, error) {
	pullRequests, _, err := gr.pullRequestService.List(context.TODO(), owner, repository, nil)
	if err != nil {
		return nil, err
	}

	var result []models.PullRequest
	for _, pullRequest := range pullRequests {
		pr := models.PullRequest{
			Title:      fmt.Sprintf("PR#%d - %s", pullRequest.GetNumber(), pullRequest.GetTitle()),
			Owner:      owner,
			Repository: repository,
			Ref:        pullRequest.Head.GetRef(),
		}

		result = append(result, pr)
	}

	return result, nil
}

func (gr *githubRepository) GetCommit(owner, repository, sha string) (*models.Commit, error) {
	commit, _, err := gr.gitService.GetCommit(context.TODO(), owner, repository, sha)
	if err != nil {
		return nil, err
	}

	result := &models.Commit{
		SHA: sha,
		Author: &models.Author{
			Name:      commit.Author.GetName(),
			AvatarURL: gravatar.GetGravatarURL(commit.Author.GetEmail()),
		},
	}

	return result, nil
}
