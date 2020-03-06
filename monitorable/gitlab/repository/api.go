package repository

import (
	"net/http"
	"time"

	"github.com/gregjones/httpcache"
	"github.com/monitoror/monitoror/models"

	"github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/monitorable/gitlab"
	gitlabModels "github.com/monitoror/monitoror/monitorable/gitlab/models"
	"github.com/monitoror/monitoror/pkg/gogitlab"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/gravatar"
	gitlabApi "github.com/xanzy/go-gitlab"
)

type (
	gitlabRepository struct {
		searchService       gogitlab.SearchService
		mergeRequestService gogitlab.MergeRequestsService
		commitsService      gogitlab.CommitsService
		pipelinesService    gogitlab.PipelinesService

		config *config.Gitlab
	}
)

func NewGitlabRepository(config *config.Gitlab) gitlab.Repository {
	httpClient := &http.Client{
		Transport: httpcache.NewMemoryCacheTransport(),
		Timeout:   time.Duration(config.Timeout) * time.Millisecond,
	}

	// Init Gitlab Client
	client := gitlabApi.NewClient(httpClient, config.Token)

	return &gitlabRepository{
		searchService:       client.Search,
		mergeRequestService: client.MergeRequests,
		commitsService:      client.Commits,
		pipelinesService:    client.Pipelines,
		config:              config,
	}
}

func (gr *gitlabRepository) GetCount(query string) (int, error) {
	_, response, err := gr.searchService.Issues(query, &gitlabApi.SearchOptions{})
	if err != nil {
		return 0, err
	}

	return response.TotalItems, err
}

func (gr *gitlabRepository) GetPipelines(repository, ref string) (*gitlabModels.Pipelines, error) {
	pipelines := &gitlabModels.Pipelines{
		Runs: []gitlabModels.Run{},
	}

	projectPipelines, _, err := gr.pipelinesService.ListProjectPipelines(
		repository, &gitlabApi.ListProjectPipelinesOptions{
			Ref:     gitlabApi.String(ref),
			OrderBy: gitlabApi.String("updated_at"),
			Sort:    gitlabApi.String("desc"),
		},
	)
	if err != nil {
		return nil, err
	}

	for _, info := range projectPipelines {
		p, _, err := gr.pipelinesService.GetPipeline(
			repository, info.ID,
		)
		if err != nil {
			return nil, err
		}

		run := gitlabModels.Run{
			ID:        p.ID,
			Status:    p.Status,
			Duration:  p.Duration,
			CreatedAt: *p.CreatedAt,
		}

		if p.StartedAt != nil {
			run.StartedAt = p.StartedAt
		}

		if p.FinishedAt != nil {
			run.FinishedAt = p.FinishedAt
		}

		pipelines.HeadCommit = p.SHA
		pipelines.Runs = append(pipelines.Runs, run)
	}

	return pipelines, nil
}

func (gr *gitlabRepository) GetMergeRequests(repository string) ([]gitlabModels.MergeRequest, error) {
	mergeRequests, _, err := gr.mergeRequestService.ListProjectMergeRequests(
		repository, &gitlabApi.ListProjectMergeRequestsOptions{
			State: gitlabApi.String("opened"),
		},
	)
	if err != nil {
		return nil, err
	}

	var result []gitlabModels.MergeRequest
	for _, mr := range mergeRequests {
		pr := gitlabModels.MergeRequest{
			ID:         mr.ID,
			Repository: repository,
			Ref:        mr.SourceBranch,
		}

		result = append(result, pr)
	}

	return result, nil
}

func (gr *gitlabRepository) GetCommit(repository, sha string) (*gitlabModels.Commit, error) {
	commit, _, err := gr.commitsService.GetCommit(repository, sha)
	if err != nil {
		return nil, err
	}

	result := &gitlabModels.Commit{
		SHA: sha,
		Author: &models.Author{
			Name:      commit.AuthorName,
			AvatarURL: gravatar.GetGravatarURL(commit.AuthorEmail),
		},
	}

	return result, nil
}
