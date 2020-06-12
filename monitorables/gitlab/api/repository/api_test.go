package repository

import (
	"errors"
	"testing"
	"time"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/gitlab/api/models"
	"github.com/monitoror/monitoror/monitorables/gitlab/config"
	"github.com/monitoror/monitoror/pkg/gogitlab/mocks"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
	"github.com/xanzy/go-gitlab"
)

func initRepository(t *testing.T) *gitlabRepository {
	conf := &config.Gitlab{
		URL:     "https://gitlab.example.com",
		Token:   "xxx",
		Timeout: 1000,
	}

	repository := NewGitlabRepository(conf)

	apiGithubRepository, ok := repository.(*gitlabRepository)
	if assert.True(t, ok) {
		return apiGithubRepository
	}
	return nil
}

func TestNewGitlabRepository_Panic(t *testing.T) {
	conf := &config.Gitlab{
		URL: "test%test",
	}

	assert.Panics(t, func() {
		NewGitlabRepository(conf)
	})
}

func TestRepository_GetCountIssues_Error(t *testing.T) {
	gitlabErr := errors.New("gitlab error")

	mockIssueService := new(mocks.IssuesService)
	mockIssueService.On("ListIssues", Anything, Anything).
		Return(nil, nil, gitlabErr)
	mockIssueService.On("ListProjectIssues", Anything, Anything, Anything).
		Return(nil, nil, gitlabErr)

	repository := initRepository(t)
	if repository != nil {
		repository.issuesService = mockIssueService

		_, err := repository.GetCountIssues(&models.IssuesParams{})
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "gitlab error")
		}
		_, err = repository.GetCountIssues(&models.IssuesParams{ProjectID: pointer.ToInt(10)})
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "gitlab error")
		}

		mockIssueService.AssertNumberOfCalls(t, "ListIssues", 1)
		mockIssueService.AssertNumberOfCalls(t, "ListProjectIssues", 1)
		mockIssueService.AssertExpectations(t)
	}
}

func TestRepository_GetCountIssues_Success(t *testing.T) {
	mockIssueService := new(mocks.IssuesService)
	mockIssueService.On("ListIssues", Anything, Anything).
		Return(nil, &gitlab.Response{TotalItems: 42}, nil)
	mockIssueService.On("ListProjectIssues", Anything, Anything, Anything).
		Return(nil, &gitlab.Response{TotalItems: 42}, nil)

	repository := initRepository(t)
	if repository != nil {
		repository.issuesService = mockIssueService

		value, err := repository.GetCountIssues(&models.IssuesParams{})
		if assert.NoError(t, err) {
			assert.Equal(t, 42, value)
		}

		value, err = repository.GetCountIssues(&models.IssuesParams{ProjectID: pointer.ToInt(10)})
		if assert.NoError(t, err) {
			assert.Equal(t, 42, value)
		}

		mockIssueService.AssertNumberOfCalls(t, "ListIssues", 1)
		mockIssueService.AssertNumberOfCalls(t, "ListProjectIssues", 1)
		mockIssueService.AssertExpectations(t)
	}
}

func TestRepository_GetPipeline_Error(t *testing.T) {
	gitlabErr := errors.New("gitlab error")

	mockPipelineService := new(mocks.PipelinesService)
	mockPipelineService.On("GetPipeline", Anything, AnythingOfType("int"), Anything).
		Return(nil, nil, gitlabErr)

	repository := initRepository(t)
	if repository != nil {
		repository.pipelinesService = mockPipelineService

		_, err := repository.GetPipeline(10, 10)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "gitlab error")
			mockPipelineService.AssertNumberOfCalls(t, "GetPipeline", 1)
			mockPipelineService.AssertExpectations(t)
		}
	}
}

func TestRepository_GetPipeline_Success(t *testing.T) {
	now := time.Now()

	gitlabPipeline := &gitlab.Pipeline{
		ID:     10,
		Status: "failed",
		Ref:    "master",
		SHA:    "12345",
		User: &gitlab.BasicUser{
			Username:  "test",
			AvatarURL: "test.example.com",
		},
		CreatedAt:  pointer.ToTime(now),
		FinishedAt: pointer.ToTime(now.Add(time.Second * 30)),
	}

	mockPipelineService := new(mocks.PipelinesService)
	mockPipelineService.On("GetPipeline", Anything, AnythingOfType("int"), Anything).
		Return(gitlabPipeline, nil, nil)

	pipeline := &models.Pipeline{
		ID:     10,
		Branch: "master",
		Author: coreModels.Author{
			Name:      "test",
			AvatarURL: "test.example.com",
		},
		Status:     "failed",
		StartedAt:  pointer.ToTime(now),
		FinishedAt: pointer.ToTime(now.Add(time.Second * 30)),
	}

	repository := initRepository(t)
	if repository != nil {
		repository.pipelinesService = mockPipelineService

		result, err := repository.GetPipeline(10, 10)
		if assert.NoError(t, err) {
			assert.Equal(t, pipeline, result)
			mockPipelineService.AssertNumberOfCalls(t, "GetPipeline", 1)
			mockPipelineService.AssertExpectations(t)
		}
	}
}

func TestRepository_GetPipelines_Error(t *testing.T) {
	gitlabErr := errors.New("gitlab error")

	mockPipelineService := new(mocks.PipelinesService)
	mockPipelineService.On("ListProjectPipelines", Anything, Anything, Anything).
		Return(nil, nil, gitlabErr)

	repository := initRepository(t)
	if repository != nil {
		repository.pipelinesService = mockPipelineService

		_, err := repository.GetPipelines(10, "master")
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "gitlab error")
			mockPipelineService.AssertNumberOfCalls(t, "ListProjectPipelines", 1)
			mockPipelineService.AssertExpectations(t)
		}
	}
}

func TestRepository_GetPipelines_Success(t *testing.T) {
	mockPipelineService := new(mocks.PipelinesService)
	mockPipelineService.On("ListProjectPipelines", Anything, Anything, Anything).
		Return([]*gitlab.PipelineInfo{
			{ID: 10},
			{ID: 11},
			{ID: 12},
		}, nil, nil)

	repository := initRepository(t)
	if repository != nil {
		repository.pipelinesService = mockPipelineService

		pipelines, err := repository.GetPipelines(10, "master")
		if assert.NoError(t, err) {
			assert.Equal(t, []int{10, 11, 12}, pipelines)
			mockPipelineService.AssertNumberOfCalls(t, "ListProjectPipelines", 1)
			mockPipelineService.AssertExpectations(t)
		}
	}
}

func TestRepository_GetMergeRequest_Error(t *testing.T) {
	gitlabErr := errors.New("gitlab error")

	mockMergeRequestService := new(mocks.MergeRequestsService)
	mockMergeRequestService.On("GetMergeRequest", Anything, AnythingOfType("int"), Anything, Anything).
		Return(nil, nil, gitlabErr)

	repository := initRepository(t)
	if repository != nil {
		repository.mergeRequestsService = mockMergeRequestService

		_, err := repository.GetMergeRequest(10, 10)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "gitlab error")
			mockMergeRequestService.AssertNumberOfCalls(t, "GetMergeRequest", 1)
			mockMergeRequestService.AssertExpectations(t)
		}
	}
}

func TestRepository_GetMergeRequest_Success(t *testing.T) {
	gitlabMergeRequest := &gitlab.MergeRequest{
		ID:              10,
		IID:             20,
		Title:           "Test",
		SourceProjectID: 30,
		SourceBranch:    "master",
		SHA:             "12345",
		Author: &gitlab.BasicUser{
			Username:  "test",
			AvatarURL: "test.example.com",
		},
	}

	mockMergeRequestService := new(mocks.MergeRequestsService)
	mockMergeRequestService.On("GetMergeRequest", Anything, AnythingOfType("int"), Anything, Anything).
		Return(gitlabMergeRequest, nil, nil)

	mergeRequest := &models.MergeRequest{
		ID:    20,
		Title: "Test",
		Author: coreModels.Author{
			Name:      "test",
			AvatarURL: "test.example.com",
		},
		SourceProjectID: 30,
		SourceBranch:    "master",
		CommitSHA:       "12345",
	}

	repository := initRepository(t)
	if repository != nil {
		repository.mergeRequestsService = mockMergeRequestService

		result, err := repository.GetMergeRequest(10, 10)
		if assert.NoError(t, err) {
			assert.Equal(t, mergeRequest, result)
			mockMergeRequestService.AssertNumberOfCalls(t, "GetMergeRequest", 1)
			mockMergeRequestService.AssertExpectations(t)
		}
	}
}

func TestRepository_GetMergeRequests_Error(t *testing.T) {
	gitlabErr := errors.New("gitlab error")

	mockMergeRequestService := new(mocks.MergeRequestsService)
	mockMergeRequestService.On("ListProjectMergeRequests", Anything, Anything, Anything).
		Return(nil, nil, gitlabErr)

	repository := initRepository(t)
	if repository != nil {
		repository.mergeRequestsService = mockMergeRequestService

		_, err := repository.GetMergeRequests(10)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "gitlab error")
			mockMergeRequestService.AssertNumberOfCalls(t, "ListProjectMergeRequests", 1)
			mockMergeRequestService.AssertExpectations(t)
		}
	}
}

func TestRepository_GetMergeRequests_Success(t *testing.T) {
	gitlabMergeRequest := &gitlab.MergeRequest{
		ID:              10,
		IID:             20,
		Title:           "Test",
		SourceProjectID: 30,
		SourceBranch:    "master",
		SHA:             "12345",
		Author: &gitlab.BasicUser{
			Username:  "test",
			AvatarURL: "test.example.com",
		},
	}

	mockMergeRequestService := new(mocks.MergeRequestsService)
	mockMergeRequestService.On("ListProjectMergeRequests", Anything, Anything, Anything).
		Return([]*gitlab.MergeRequest{gitlabMergeRequest}, nil, nil)

	mergeRequest := models.MergeRequest{
		ID:    20,
		Title: "Test",
		Author: coreModels.Author{
			Name:      "test",
			AvatarURL: "test.example.com",
		},
		SourceProjectID: 30,
		SourceBranch:    "master",
		CommitSHA:       "12345",
	}

	repository := initRepository(t)
	if repository != nil {
		repository.mergeRequestsService = mockMergeRequestService

		result, err := repository.GetMergeRequests(10)
		if assert.NoError(t, err) {
			assert.Len(t, result, 1)
			assert.Equal(t, mergeRequest, result[0])
			mockMergeRequestService.AssertNumberOfCalls(t, "ListProjectMergeRequests", 1)
			mockMergeRequestService.AssertExpectations(t)
		}
	}
}

func TestRepository_GetProject_Error(t *testing.T) {
	gitlabErr := errors.New("gitlab error")

	mockProjectService := new(mocks.ProjectService)
	mockProjectService.On("GetProject", Anything, Anything, Anything).
		Return(nil, nil, gitlabErr)

	repository := initRepository(t)
	if repository != nil {
		repository.projectService = mockProjectService

		_, err := repository.GetProject(10)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "gitlab error")
			mockProjectService.AssertNumberOfCalls(t, "GetProject", 1)
			mockProjectService.AssertExpectations(t)
		}
	}
}

func TestRepository_GetProject_Success(t *testing.T) {
	gitlabProject := &gitlab.Project{
		ID:        10,
		Path:      "test2",
		Namespace: &gitlab.ProjectNamespace{Path: "test1"},
	}

	mockProjectService := new(mocks.ProjectService)
	mockProjectService.On("GetProject", Anything, Anything, Anything).
		Return(gitlabProject, nil, nil)

	project := &models.Project{
		ID:         10,
		Owner:      "test1",
		Repository: "test2",
	}

	repository := initRepository(t)
	if repository != nil {
		repository.projectService = mockProjectService

		result, err := repository.GetProject(10)
		if assert.NoError(t, err) {
			assert.Equal(t, project, result)
			mockProjectService.AssertNumberOfCalls(t, "GetProject", 1)
			mockProjectService.AssertExpectations(t)
		}
	}
}
