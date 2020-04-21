package repository

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/github/api"
	"github.com/monitoror/monitoror/monitorables/github/api/models"
	"github.com/monitoror/monitoror/monitorables/github/config"
	"github.com/monitoror/monitoror/pkg/gogithub"
	"github.com/monitoror/monitoror/pkg/gravatar"

	githubApi "github.com/google/go-github/github"
	"github.com/sourcegraph/httpcache"
	"golang.org/x/oauth2"
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

func NewGithubRepository(config *config.Github) api.Repository {
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

	// Add / if missing
	if !strings.HasSuffix(config.URL, "/") {
		config.URL += "/"
	}

	// Set baseURL
	baseURL, _ := url.Parse(config.URL) // Err check already done in Monitorable.Validate
	client.BaseURL = baseURL
	client.UploadURL = nil

	return &githubRepository{
		searchService:       client.Search,
		checksService:       client.Checks,
		repositoriesService: client.Repositories,
		pullRequestService:  client.PullRequests,
		gitService:          client.Git,
		config:              config,
	}
}

func (gr *githubRepository) GetCount(query string) (int, error) {
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

func (gr *githubRepository) GetPullRequest(owner, repository string, id int) (*models.PullRequest, error) {
	pullRequest, _, err := gr.pullRequestService.Get(context.TODO(), owner, repository, id)
	if err != nil {
		return nil, err
	}

	return fillPullRequest(pullRequest), nil
}

func (gr *githubRepository) GetPullRequests(owner, repository string) ([]models.PullRequest, error) {
	pullRequests, _, err := gr.pullRequestService.List(context.TODO(), owner, repository, nil)
	if err != nil {
		return nil, err
	}

	var result []models.PullRequest
	for _, pullRequest := range pullRequests {
		result = append(result, *fillPullRequest(pullRequest))
	}

	return result, nil
}

func fillPullRequest(pr *githubApi.PullRequest) *models.PullRequest {
	return &models.PullRequest{
		ID:    pr.GetNumber(),
		Title: pr.GetTitle(),
		Author: coreModels.Author{
			Name:      pr.GetUser().GetLogin(),
			AvatarURL: pr.GetUser().GetAvatarURL(),
		},
		Owner:      pr.GetHead().GetUser().GetLogin(),
		Repository: pr.GetHead().GetRepo().GetName(),
		Branch:     pr.GetHead().GetRef(),
		CommitSHA:  pr.GetHead().GetSHA(),
	}
}

func (gr *githubRepository) GetCommit(owner, repository, sha string) (*models.Commit, error) {
	commit, _, err := gr.gitService.GetCommit(context.TODO(), owner, repository, sha)
	if err != nil {
		return nil, err
	}

	result := &models.Commit{
		SHA: sha,
		Author: &coreModels.Author{
			Name:      commit.Author.GetLogin(),
			AvatarURL: gravatar.GetGravatarURL(commit.Author.GetEmail()),
		},
	}

	return result, nil
}
