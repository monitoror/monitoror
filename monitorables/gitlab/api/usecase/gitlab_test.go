package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/gitlab/api"
	"github.com/monitoror/monitoror/monitorables/gitlab/api/mocks"
	"github.com/monitoror/monitoror/monitorables/gitlab/api/models"

	"github.com/jsdidierlaurent/echo-middleware/cache"
	"github.com/stretchr/testify/mock"
)

func initUsecase(mockRepository api.Repository) *gitlabUsecase {
	store := cache.NewGoCacheStore(time.Minute*5, time.Second)
	gu := NewGitlabUsecase(mockRepository, store)
	castedGu := gu.(*gitlabUsecase)
	return castedGu
}

func TestUsecase_CountIssues_Error(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetCountIssues", mock.Anything).
		Return(0, errors.New("boom"))

	gu := initUsecase(mockRepository)

	tile, err := gu.CountIssues(&models.IssuesParams{})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &coreModels.MonitororError{}, err)
		assert.Equal(t, "unable to load issues", err.Error())
		mockRepository.AssertNumberOfCalls(t, "GetCountIssues", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestUsecase_CountIssues_Success(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetCountIssues", mock.Anything).
		Return(42, nil)

	gu := initUsecase(mockRepository)

	expected := coreModels.NewTile(api.GitlabCountIssuesTileType).WithValue(coreModels.NumberUnit)
	expected.Label = "GitLab count"
	expected.Status = coreModels.SuccessStatus
	expected.Value.Values = []string{"42"}

	tile, err := gu.CountIssues(&models.IssuesParams{})
	if assert.NoError(t, err) {
		assert.Equal(t, expected, tile)
		mockRepository.AssertNumberOfCalls(t, "GetCountIssues", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestUsecase_Pipeline_ErrorProject(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetProject", mock.Anything).
		Return(nil, errors.New("boom"))

	gu := initUsecase(mockRepository)

	tile, err := gu.Pipeline(&models.PipelineParams{ProjectID: pointer.ToInt(10), Ref: "master"})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &coreModels.MonitororError{}, err)
		assert.Equal(t, "unable to load project", err.Error())
		mockRepository.AssertNumberOfCalls(t, "GetProject", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestUsecase_Pipeline_ErrorPipelines(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetProject", mock.Anything).
		Return(&models.Project{Repository: "Test"}, nil)
	mockRepository.On("GetPipelines", mock.Anything, mock.Anything).
		Return(nil, errors.New("boom"))

	gu := initUsecase(mockRepository)

	tile, err := gu.Pipeline(&models.PipelineParams{ProjectID: pointer.ToInt(10), Ref: "master"})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &coreModels.MonitororError{}, err)
		assert.Equal(t, "unable to load pipelines", err.Error())
		mockRepository.AssertNumberOfCalls(t, "GetProject", 1)
		mockRepository.AssertNumberOfCalls(t, "GetPipelines", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestUsecase_Pipeline_NoPipelines(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetProject", mock.Anything).
		Return(&models.Project{Repository: "Test"}, nil)
	mockRepository.On("GetPipelines", mock.Anything, mock.Anything).
		Return([]int{}, nil)

	gu := initUsecase(mockRepository)

	tile, err := gu.Pipeline(&models.PipelineParams{ProjectID: pointer.ToInt(10), Ref: "master"})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &coreModels.MonitororError{}, err)
		assert.Equal(t, "no pipelines found", err.Error())
		mockRepository.AssertNumberOfCalls(t, "GetProject", 1)
		mockRepository.AssertNumberOfCalls(t, "GetPipelines", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestUsecase_Pipeline_ErrorPipeline(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetProject", mock.Anything).
		Return(&models.Project{Repository: "Test"}, nil)
	mockRepository.On("GetPipelines", mock.Anything, mock.Anything).
		Return([]int{10}, nil)
	mockRepository.On("GetPipeline", mock.Anything, mock.Anything).
		Return(nil, errors.New("boom"))

	gu := initUsecase(mockRepository)

	tile, err := gu.Pipeline(&models.PipelineParams{ProjectID: pointer.ToInt(10), Ref: "master"})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &coreModels.MonitororError{}, err)
		assert.Equal(t, "unable to load pipeline", err.Error())
		mockRepository.AssertNumberOfCalls(t, "GetProject", 1)
		mockRepository.AssertNumberOfCalls(t, "GetPipelines", 1)
		mockRepository.AssertNumberOfCalls(t, "GetPipeline", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestUsecase_Pipeline_Success(t *testing.T) {
	refTime := time.Now()
	startedAt := refTime.Add(-time.Second * 30)
	finishedAt := refTime.Add(-time.Second * 15)

	pipeline := &models.Pipeline{
		ID:         10,
		Branch:     "master",
		Status:     "success",
		StartedAt:  &startedAt,
		FinishedAt: &finishedAt,
	}

	expected := coreModels.NewTile(api.GitlabPipelineTileType).WithBuild()
	expected.Label = "project"
	expected.Build.Branch = pointer.ToString("master")
	expected.Status = coreModels.SuccessStatus
	expected.Build.PreviousStatus = coreModels.UnknownStatus
	expected.Build.StartedAt = pointer.ToTime(startedAt)
	expected.Build.FinishedAt = pointer.ToTime(finishedAt)

	testPipeline(t, pipeline, expected)
}

func TestUsecase_Pipeline_Failed(t *testing.T) {
	refTime := time.Now()
	startedAt := refTime.Add(-time.Second * 30)
	finishedAt := refTime.Add(-time.Second * 15)

	pipeline := &models.Pipeline{
		ID:     10,
		Branch: "master",
		Author: coreModels.Author{
			Name:      "author",
			AvatarURL: "author.exemple.com",
		},
		Status:     "failed",
		StartedAt:  &startedAt,
		FinishedAt: &finishedAt,
	}

	expected := coreModels.NewTile(api.GitlabPipelineTileType).WithBuild()
	expected.Label = "project"
	expected.Build.Branch = pointer.ToString("master")
	expected.Status = coreModels.FailedStatus
	expected.Build.PreviousStatus = coreModels.UnknownStatus
	expected.Build.StartedAt = pointer.ToTime(startedAt)
	expected.Build.FinishedAt = pointer.ToTime(finishedAt)
	expected.Build.Author = &coreModels.Author{
		Name:      "author",
		AvatarURL: "author.exemple.com",
	}

	testPipeline(t, pipeline, expected)
}

func TestUsecase_Pipeline_Running(t *testing.T) {
	refTime := time.Now()
	startedAt := refTime.Add(-time.Second * 30)

	pipeline := &models.Pipeline{
		ID:        10,
		Branch:    "master",
		Status:    "running",
		StartedAt: &startedAt,
	}

	expected := coreModels.NewTile(api.GitlabPipelineTileType).WithBuild()
	expected.Label = "project"
	expected.Build.Branch = pointer.ToString("master")
	expected.Status = coreModels.RunningStatus
	expected.Build.PreviousStatus = coreModels.UnknownStatus
	expected.Build.StartedAt = pointer.ToTime(startedAt)
	expected.Build.Duration = pointer.ToInt64(30)
	expected.Build.EstimatedDuration = pointer.ToInt64(0)

	testPipeline(t, pipeline, expected)
}

func TestUsecase_Pipeline_Queued(t *testing.T) {
	refTime := time.Now()
	startedAt := refTime.Add(-time.Second * 30)

	pipeline := &models.Pipeline{
		ID:        10,
		Branch:    "master",
		Status:    "pending",
		StartedAt: &startedAt,
	}

	expected := coreModels.NewTile(api.GitlabPipelineTileType).WithBuild()
	expected.Label = "project"
	expected.Build.Branch = pointer.ToString("master")
	expected.Status = coreModels.QueuedStatus
	expected.Build.PreviousStatus = coreModels.UnknownStatus
	expected.Build.StartedAt = pointer.ToTime(startedAt)

	testPipeline(t, pipeline, expected)
}

func testPipeline(t *testing.T, pipeline *models.Pipeline, expected *coreModels.Tile) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetProject", mock.Anything).
		Return(&models.Project{Repository: "project"}, nil)
	mockRepository.On("GetPipelines", mock.Anything, mock.Anything).
		Return([]int{10}, nil)
	mockRepository.On("GetPipeline", mock.Anything, mock.Anything).
		Return(pipeline, nil)

	gu := initUsecase(mockRepository)

	tile, err := gu.Pipeline(&models.PipelineParams{ProjectID: pointer.ToInt(10), Ref: "master"})
	if assert.NoError(t, err) {
		assert.Equal(t, expected, tile)
		mockRepository.AssertNumberOfCalls(t, "GetProject", 1)
		mockRepository.AssertNumberOfCalls(t, "GetPipelines", 1)
		mockRepository.AssertNumberOfCalls(t, "GetPipeline", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestUsecase_Pipeline_WithPrevious(t *testing.T) {
	refTime := time.Now()
	startedAt := refTime.Add(-time.Second * 30)
	finishedAt := refTime.Add(-time.Second * 15)

	pipeline := &models.Pipeline{
		ID:         10,
		Branch:     "master",
		Status:     "success",
		StartedAt:  &startedAt,
		FinishedAt: &finishedAt,
	}

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetProject", mock.Anything).
		Return(&models.Project{Repository: "project"}, nil)
	mockRepository.On("GetPipelines", mock.Anything, mock.Anything).
		Return([]int{10}, nil)
	mockRepository.On("GetPipeline", mock.Anything, mock.Anything).
		Return(pipeline, nil)

	gu := initUsecase(mockRepository)

	expected := coreModels.NewTile(api.GitlabPipelineTileType).WithBuild()
	expected.Label = "project"
	expected.Build.Branch = pointer.ToString("master")
	expected.Status = coreModels.SuccessStatus
	expected.Build.PreviousStatus = coreModels.UnknownStatus
	expected.Build.StartedAt = pointer.ToTime(startedAt)
	expected.Build.FinishedAt = pointer.ToTime(finishedAt)

	tile, err := gu.Pipeline(&models.PipelineParams{ProjectID: pointer.ToInt(10), Ref: "master"})
	if assert.NoError(t, err) {
		assert.Equal(t, expected, tile)
	}

	pipeline.ID = 20
	pipeline.Status = "running"
	pipeline.FinishedAt = nil

	expected.Status = coreModels.RunningStatus
	expected.Build.PreviousStatus = coreModels.SuccessStatus
	expected.Build.Duration = pointer.ToInt64(30)
	expected.Build.EstimatedDuration = pointer.ToInt64(15)
	expected.Build.FinishedAt = nil

	tile, err = gu.Pipeline(&models.PipelineParams{ProjectID: pointer.ToInt(10), Ref: "master"})
	if assert.NoError(t, err) {
		assert.Equal(t, expected, tile)
	}

	mockRepository.AssertNumberOfCalls(t, "GetProject", 1)
	mockRepository.AssertNumberOfCalls(t, "GetPipelines", 2)
	mockRepository.AssertNumberOfCalls(t, "GetPipeline", 2)
	mockRepository.AssertExpectations(t)
}

func TestUsecase_MergeRequest_ErrorProject(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetProject", mock.Anything).
		Return(nil, errors.New("boom"))

	gu := initUsecase(mockRepository)

	tile, err := gu.MergeRequest(&models.MergeRequestParams{ProjectID: pointer.ToInt(10), ID: pointer.ToInt(10)})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &coreModels.MonitororError{}, err)
		assert.Equal(t, "unable to load project", err.Error())
		mockRepository.AssertNumberOfCalls(t, "GetProject", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestUsecase_MergeRequest_ErrorMergeRequest(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetProject", mock.Anything).
		Return(&models.Project{Repository: "Test"}, nil)
	mockRepository.On("GetMergeRequest", mock.Anything, mock.Anything).
		Return(nil, errors.New("boom"))

	gu := initUsecase(mockRepository)

	tile, err := gu.MergeRequest(&models.MergeRequestParams{ProjectID: pointer.ToInt(10), ID: pointer.ToInt(10)})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &coreModels.MonitororError{}, err)
		assert.Equal(t, "unable to load merge request", err.Error())
		mockRepository.AssertNumberOfCalls(t, "GetProject", 1)
		mockRepository.AssertNumberOfCalls(t, "GetMergeRequest", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestUsecase_MergeRequest_ErrorSourceProject(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetProject", mock.Anything).
		Return(&models.Project{Repository: "Test"}, nil).Once()
	mockRepository.On("GetProject", mock.Anything).
		Return(nil, errors.New("boom"))
	mockRepository.On("GetMergeRequest", mock.Anything, mock.Anything).
		Return(&models.MergeRequest{SourceProjectID: 20}, nil)

	gu := initUsecase(mockRepository)

	tile, err := gu.MergeRequest(&models.MergeRequestParams{ProjectID: pointer.ToInt(10), ID: pointer.ToInt(10)})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &coreModels.MonitororError{}, err)
		assert.Equal(t, "unable to load project", err.Error())
		mockRepository.AssertNumberOfCalls(t, "GetProject", 2)
		mockRepository.AssertNumberOfCalls(t, "GetMergeRequest", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestUsecase_MergeRequest_ErrorPipelines(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetProject", mock.Anything).
		Return(&models.Project{Repository: "Test"}, nil)
	mockRepository.On("GetMergeRequest", mock.Anything, mock.Anything).
		Return(&models.MergeRequest{SourceProjectID: 20}, nil)
	mockRepository.On("GetMergeRequestPipelines", mock.Anything, mock.Anything).
		Return(nil, errors.New("boom"))

	gu := initUsecase(mockRepository)

	tile, err := gu.MergeRequest(&models.MergeRequestParams{ProjectID: pointer.ToInt(10), ID: pointer.ToInt(10)})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &coreModels.MonitororError{}, err)
		assert.Equal(t, "unable to load pipelines", err.Error())
		mockRepository.AssertNumberOfCalls(t, "GetProject", 2)
		mockRepository.AssertNumberOfCalls(t, "GetMergeRequest", 1)
		mockRepository.AssertNumberOfCalls(t, "GetMergeRequestPipelines", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestUsecase_MergeRequest_ErrorPipeline(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetProject", mock.Anything).
		Return(&models.Project{Repository: "Test"}, nil)
	mockRepository.On("GetMergeRequest", mock.Anything, mock.Anything).
		Return(&models.MergeRequest{SourceProjectID: 20}, nil)
	mockRepository.On("GetMergeRequestPipelines", mock.Anything, mock.Anything).
		Return([]int{30}, nil)
	mockRepository.On("GetPipeline", mock.Anything, mock.Anything).
		Return(nil, errors.New("boom"))

	gu := initUsecase(mockRepository)

	tile, err := gu.MergeRequest(&models.MergeRequestParams{ProjectID: pointer.ToInt(10), ID: pointer.ToInt(10)})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &coreModels.MonitororError{}, err)
		assert.Equal(t, "unable to load pipeline", err.Error())
		mockRepository.AssertNumberOfCalls(t, "GetProject", 2)
		mockRepository.AssertNumberOfCalls(t, "GetMergeRequest", 1)
		mockRepository.AssertNumberOfCalls(t, "GetMergeRequestPipelines", 1)
		mockRepository.AssertNumberOfCalls(t, "GetPipeline", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestUsecase_MergeRequest_NoPipelines(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetProject", mock.Anything).
		Return(&models.Project{Repository: "Test"}, nil)
	mockRepository.On("GetMergeRequest", mock.Anything, mock.Anything).
		Return(&models.MergeRequest{SourceProjectID: 20}, nil)
	mockRepository.On("GetMergeRequestPipelines", mock.Anything, mock.Anything).
		Return([]int{}, nil)

	gu := initUsecase(mockRepository)

	tile, err := gu.MergeRequest(&models.MergeRequestParams{ProjectID: pointer.ToInt(10), ID: pointer.ToInt(10)})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &coreModels.MonitororError{}, err)
		assert.Equal(t, "no pipelines found", err.Error())
		mockRepository.AssertNumberOfCalls(t, "GetProject", 2)
		mockRepository.AssertNumberOfCalls(t, "GetMergeRequest", 1)
		mockRepository.AssertNumberOfCalls(t, "GetMergeRequestPipelines", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestUsecase_MergeRequest_Success(t *testing.T) {
	refTime := time.Now()
	startedAt := refTime.Add(-time.Second * 30)
	finishedAt := refTime.Add(-time.Second * 15)

	mergeRequest := &models.MergeRequest{
		ID:    10,
		Title: "Test MR",
		Author: coreModels.Author{
			Name:      "author",
			AvatarURL: "author.example.com",
		},
		SourceProjectID: 20,
		SourceBranch:    "master",
		CommitSHA:       "12345",
	}

	pipeline := &models.Pipeline{
		ID:         10,
		Branch:     "master",
		Status:     "failed",
		StartedAt:  &startedAt,
		FinishedAt: &finishedAt,
	}

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetProject", mock.Anything).
		Return(&models.Project{Repository: "project"}, nil).Once()
	mockRepository.On("GetProject", mock.Anything).
		Return(&models.Project{Owner: "faker", Repository: "project"}, nil)
	mockRepository.On("GetMergeRequest", mock.Anything, mock.Anything).
		Return(mergeRequest, nil)
	mockRepository.On("GetMergeRequestPipelines", mock.Anything, mock.Anything).
		Return([]int{30}, nil)
	mockRepository.On("GetPipeline", mock.Anything, mock.Anything).
		Return(pipeline, nil)

	gu := initUsecase(mockRepository)

	expected := coreModels.NewTile(api.GitlabMergeRequestTileType).WithBuild()
	expected.Label = "project"
	expected.Build.Branch = pointer.ToString("faker:master")
	expected.Build.MergeRequest = &coreModels.TileMergeRequest{
		ID:    10,
		Title: "Test MR",
	}
	expected.Build.Author = &coreModels.Author{
		Name:      "author",
		AvatarURL: "author.example.com",
	}
	expected.Status = coreModels.FailedStatus
	expected.Build.PreviousStatus = coreModels.UnknownStatus
	expected.Build.StartedAt = pointer.ToTime(startedAt)
	expected.Build.FinishedAt = pointer.ToTime(finishedAt)

	tile, err := gu.MergeRequest(&models.MergeRequestParams{ProjectID: pointer.ToInt(10), ID: pointer.ToInt(10)})
	if assert.NoError(t, err) {
		assert.Equal(t, expected, tile)
		mockRepository.AssertNumberOfCalls(t, "GetProject", 2)
		mockRepository.AssertNumberOfCalls(t, "GetMergeRequest", 1)
		mockRepository.AssertNumberOfCalls(t, "GetMergeRequestPipelines", 1)
		mockRepository.AssertNumberOfCalls(t, "GetPipeline", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestUsecase_MergeRequests_Error(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetMergeRequests", mock.Anything).
		Return(nil, errors.New("boom"))

	gu := initUsecase(mockRepository)

	generated, err := gu.MergeRequestsGenerator(&models.MergeRequestGeneratorParams{ProjectID: pointer.ToInt(10)})
	if assert.Error(t, err) {
		assert.Nil(t, generated)
		assert.IsType(t, &coreModels.MonitororError{}, err)
		assert.Equal(t, "unable to load merge requests", err.Error())
		mockRepository.AssertNumberOfCalls(t, "GetMergeRequests", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestUsecase_MergeRequests_Success(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetMergeRequests", mock.Anything).
		Return([]models.MergeRequest{{ID: 10}}, nil)

	gu := initUsecase(mockRepository)

	generated, err := gu.MergeRequestsGenerator(&models.MergeRequestGeneratorParams{ProjectID: pointer.ToInt(10)})
	if assert.NoError(t, err) {
		assert.Len(t, generated, 1)
		assert.Equal(t, 10, *generated[0].Params.(*models.MergeRequestParams).ID)
		mockRepository.AssertNumberOfCalls(t, "GetMergeRequests", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestUsecase_parseStatus(t *testing.T) {
	assert.Equal(t, coreModels.RunningStatus, parseStatus("running"))
	assert.Equal(t, coreModels.QueuedStatus, parseStatus("pending"))
	assert.Equal(t, coreModels.SuccessStatus, parseStatus("success"))
	assert.Equal(t, coreModels.FailedStatus, parseStatus("failed"))
	assert.Equal(t, coreModels.CanceledStatus, parseStatus("canceled"))
	assert.Equal(t, coreModels.CanceledStatus, parseStatus("skipped"))
	assert.Equal(t, coreModels.QueuedStatus, parseStatus("created"))
	assert.Equal(t, coreModels.ActionRequiredStatus, parseStatus("manual"))
	assert.Equal(t, coreModels.UnknownStatus, parseStatus(""))
}

func TestUsecase_getProject(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetProject", mock.AnythingOfType("int")).
		Return(&models.Project{Repository: "TEST"}, nil)

	gu := initUsecase(mockRepository)

	project, err := gu.getProject(10)
	assert.NoError(t, err)
	assert.Equal(t, "TEST", project.Repository)

	project, err = gu.getProject(10)
	assert.NoError(t, err)
	assert.Equal(t, "TEST", project.Repository)

	mockRepository.AssertNumberOfCalls(t, "GetProject", 1)
	mockRepository.AssertExpectations(t)
}
