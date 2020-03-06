package repository

import (
	"errors"
	"testing"
	"time"

	"github.com/monitoror/monitoror/pkg/monitoror/utils/gravatar"

	"github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/pkg/gogitlab/mocks"

	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"

	"github.com/xanzy/go-gitlab"
)

func initRepository(t *testing.T) *gitlabRepository {
	conf := config.InitConfig()
	conf.Monitorable.Gitlab[config.DefaultVariant].Token = "test"

	repository := NewGitlabRepository(conf.Monitorable.Gitlab[config.DefaultVariant])

	apiGitlabRepository, ok := repository.(*gitlabRepository)
	if assert.True(t, ok) {
		return apiGitlabRepository
	}
	return nil
}

func TestRepository_GetSearchCount_Error(t *testing.T) {
	gitlabErr := errors.New("gitlab error")

	mocksSearchService := new(mocks.SearchService)
	mocksSearchService.On("Issues", AnythingOfType("string"), Anything).
		Return(nil, nil, gitlabErr)

	repository := initRepository(t)
	if repository != nil {
		repository.searchService = mocksSearchService

		_, err := repository.GetCount("test")
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "gitlab error")
			mocksSearchService.AssertNumberOfCalls(t, "Issues", 1)
			mocksSearchService.AssertExpectations(t)
		}
	}
}

func TestRepository_GetSearchCount_Success(t *testing.T) {
	mocksSearchService := new(mocks.SearchService)
	mocksSearchService.On("Issues", AnythingOfType("string"), Anything).
		Return(nil, &gitlab.Response{TotalItems: 42}, nil)

	repository := initRepository(t)
	if repository != nil {
		repository.searchService = mocksSearchService

		value, err := repository.GetCount("test")
		if assert.NoError(t, err) {
			assert.Equal(t, 42, value)
			mocksSearchService.AssertNumberOfCalls(t, "Issues", 1)
			mocksSearchService.AssertExpectations(t)
		}
	}
}

func TestRepository_GetPipelines_PipelineServiceListError(t *testing.T) {
	gitlabErr := errors.New("gitlab error")

	mocksPipelinesService := new(mocks.PipelinesService)
	mocksPipelinesService.
		On("ListProjectPipelines", AnythingOfType("string"), Anything).
		Return(nil, nil, gitlabErr)

	repository := initRepository(t)
	if repository != nil {
		repository.pipelinesService = mocksPipelinesService

		_, err := repository.GetPipelines("test", "master")
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "gitlab error")
			mocksPipelinesService.AssertNumberOfCalls(t, "ListProjectPipelines", 1)
			mocksPipelinesService.AssertExpectations(t)
		}
	}
}

func TestRepository_GetPipelines_PipelineServiceGetError(t *testing.T) {
	gitlabErr := errors.New("gitlab error")

	mocksPipelinesService := new(mocks.PipelinesService)
	mocksPipelinesService.
		On("ListProjectPipelines", AnythingOfType("string"), Anything).
		Return([]*gitlab.PipelineInfo{
			{ID: 1},
		}, nil, nil)
	mocksPipelinesService.
		On("GetPipeline", AnythingOfType("string"), 1, Anything).
		Return(nil, nil, gitlabErr)

	repository := initRepository(t)
	if repository != nil {
		repository.pipelinesService = mocksPipelinesService

		_, err := repository.GetPipelines("test", "master")
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "gitlab error")
			mocksPipelinesService.AssertNumberOfCalls(t, "ListProjectPipelines", 1)
			mocksPipelinesService.AssertNumberOfCalls(t, "GetPipeline", 1)
			mocksPipelinesService.AssertExpectations(t)
		}
	}
}

func TestRepository_GetPipelines_Success(t *testing.T) {
	now := time.Now()
	listPipelinesResult := []*gitlab.PipelineInfo{
		{ID: 1},
	}

	pipelineResult := &gitlab.Pipeline{
		Status:     "success",
		Duration:   1,
		CreatedAt:  &now,
		StartedAt:  &now,
		FinishedAt: &now,
		SHA:        "sha",
	}

	mocksPipelinesService := new(mocks.PipelinesService)
	mocksPipelinesService.
		On("ListProjectPipelines", AnythingOfType("string"), Anything).
		Return(listPipelinesResult, nil, nil)
	mocksPipelinesService.
		On("GetPipeline", AnythingOfType("string"), 1).
		Return(pipelineResult, nil, nil)

	repository := initRepository(t)
	if repository != nil {
		repository.pipelinesService = mocksPipelinesService

		pipelines, err := repository.GetPipelines("test", "test")
		if assert.NoError(t, err) {
			assert.Len(t, pipelines.Runs, 1)

			assert.Equal(t, pipelineResult.ID, pipelines.Runs[0].ID)
			assert.Equal(t, pipelineResult.Status, pipelines.Runs[0].Status)
			assert.Equal(t, pipelineResult.Duration, pipelines.Runs[0].Duration)
			assert.Equal(t, pipelineResult.StartedAt, pipelines.Runs[0].StartedAt)
			assert.Equal(t, pipelineResult.FinishedAt, pipelines.Runs[0].FinishedAt)
			assert.Equal(t, pipelineResult.CreatedAt, &pipelines.Runs[0].CreatedAt)

			assert.Equal(t, pipelineResult.SHA, pipelines.HeadCommit)

			mocksPipelinesService.AssertNumberOfCalls(t, "ListProjectPipelines", 1)
			mocksPipelinesService.AssertNumberOfCalls(t, "GetPipeline", 1)
			mocksPipelinesService.AssertExpectations(t)
		}
	}
}

func TestRepository_GetMergeRequest_Error(t *testing.T) {
	gitlabErr := errors.New("gitlab error")

	mocksMergeRequestsService := new(mocks.MergeRequestsService)
	mocksMergeRequestsService.On("ListProjectMergeRequests", AnythingOfType("string"), Anything).
		Return(nil, nil, gitlabErr)

	repository := initRepository(t)
	if repository != nil {
		repository.mergeRequestService = mocksMergeRequestsService

		_, err := repository.GetMergeRequests("test")
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "gitlab error")
			mocksMergeRequestsService.AssertNumberOfCalls(t, "ListProjectMergeRequests", 1)
			mocksMergeRequestsService.AssertExpectations(t)
		}
	}
}

func TestRepository_GetMergeRequest_Success(t *testing.T) {
	mocksMergeRequestService := new(mocks.MergeRequestsService)
	mocksMergeRequestService.On("ListProjectMergeRequests", AnythingOfType("string"), Anything).
		Return([]*gitlab.MergeRequest{
			{
				ID:           10,
				SourceBranch: "master",
			},
		}, nil, nil)

	repository := initRepository(t)
	if repository != nil {
		repository.mergeRequestService = mocksMergeRequestService

		mergeRequests, err := repository.GetMergeRequests("test")
		if assert.NoError(t, err) {
			assert.Len(t, mergeRequests, 1)
			assert.Equal(t, 10, mergeRequests[0].ID)
			assert.Equal(t, "test", mergeRequests[0].Repository)
			assert.Equal(t, "master", mergeRequests[0].Ref)

			mocksMergeRequestService.AssertNumberOfCalls(t, "ListProjectMergeRequests", 1)
			mocksMergeRequestService.AssertExpectations(t)
		}
	}
}

func TestRepository_GetCommit_Error(t *testing.T) {
	gitlabErr := errors.New("gitlab error")

	mocksCommitsService := new(mocks.CommitsService)
	mocksCommitsService.On("GetCommit", AnythingOfType("string"), AnythingOfType("string")).
		Return(nil, nil, gitlabErr)

	repository := initRepository(t)
	if repository != nil {
		repository.commitsService = mocksCommitsService

		_, err := repository.GetCommit("test", "sha")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "gitlab error")
		mocksCommitsService.AssertNumberOfCalls(t, "GetCommit", 1)
		mocksCommitsService.AssertExpectations(t)
	}
}

func TestRepository_GetCommit_Success(t *testing.T) {
	mocksCommitsService := new(mocks.CommitsService)
	mocksCommitsService.On("GetCommit", AnythingOfType("string"), AnythingOfType("string")).
		Return(&gitlab.Commit{
			AuthorName:  "Test",
			AuthorEmail: "test@example.com",
		}, nil, nil)

	repository := initRepository(t)
	if repository != nil {
		repository.commitsService = mocksCommitsService

		commit, err := repository.GetCommit("test", "sha")
		if assert.NoError(t, err) {
			assert.Equal(t, "sha", commit.SHA)
			assert.Equal(t, "Test", commit.Author.Name)
			assert.Equal(t, gravatar.GetGravatarURL("test@example.com"), commit.Author.AvatarURL)
			mocksCommitsService.AssertNumberOfCalls(t, "GetCommit", 1)
			mocksCommitsService.AssertExpectations(t)
		}
	}
}
