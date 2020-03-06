package usecase

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/gitlab"
	"github.com/monitoror/monitoror/monitorable/gitlab/mocks"
	. "github.com/monitoror/monitoror/monitorable/gitlab/models"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/hash"

	. "github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func TestCount_Error(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetCount", AnythingOfType("string")).
		Return(0, errors.New("boom"))

	gu := NewGitlabUsecase(mockRepository)

	tile, err := gu.Count(&CountParams{Query: "test"})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &models.MonitororError{}, err)
		assert.Equal(t, "unable to find count or wrong query", err.Error())
		mockRepository.AssertNumberOfCalls(t, "GetCount", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestCount_Success(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetCount", AnythingOfType("string")).
		Return(10, nil)

	gu := NewGitlabUsecase(mockRepository)

	expected := models.NewTile(gitlab.GitlabCountTileType).WithValue(models.NumberUnit)
	expected.Label = "test"
	expected.Status = models.SuccessStatus
	expected.Value.Values = []string{"10"}

	tile, err := gu.Count(&CountParams{Query: "test"})
	if assert.NoError(t, err) {
		assert.NotNil(t, tile)
		assert.Equal(t, expected, tile)

		mockRepository.AssertNumberOfCalls(t, "GetCount", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestPipelines_Error(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetPipelines", AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string")).
		Return(nil, errors.New("boom"))

	gu := NewGitlabUsecase(mockRepository)

	tile, err := gu.Pipelines(&PipelinesParams{Repository: "test", Ref: "master"})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &models.MonitororError{}, err)
		assert.Equal(t, "unable to find ref pipelines", err.Error())
		mockRepository.AssertNumberOfCalls(t, "GetPipelines", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestPipelines_NoPipelines(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetPipelines", AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string")).
		Return(&Pipelines{}, nil)

	gu := NewGitlabUsecase(mockRepository)

	tile, err := gu.Pipelines(&PipelinesParams{Repository: "test", Ref: "master"})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &models.MonitororError{}, err)
		assert.Equal(t, "no ref pipelines found", err.Error())
		mockRepository.AssertNumberOfCalls(t, "GetPipelines", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestPipelines_Success(t *testing.T) {
	refTime := time.Now()
	startedAt := refTime.Add(-time.Second * 30)
	finishedAt := refTime.Add(-time.Second * 15)

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetPipelines", AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string")).
		Return(&Pipelines{
			Runs: []Run{
				{
					ID:         10,
					Status:     "success",
					StartedAt:  ToTime(startedAt),
					FinishedAt: ToTime(finishedAt),
				},
			},
		}, nil)

	gu := NewGitlabUsecase(mockRepository)

	expected := models.NewTile(gitlab.GitlabPipelinesTileType).WithBuild()
	expected.Label = "test"
	expected.Build.Branch = ToString("master")

	expected.Status = models.SuccessStatus
	expected.Build.PreviousStatus = models.UnknownStatus
	expected.Build.StartedAt = ToTime(startedAt)
	expected.Build.FinishedAt = ToTime(finishedAt)

	tile, err := gu.Pipelines(&PipelinesParams{Repository: "test", Ref: "master"})
	if assert.NoError(t, err) {
		assert.NotNil(t, tile)
		assert.Equal(t, expected, tile)

		mockRepository.AssertNumberOfCalls(t, "GetPipelines", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestPipelines_Failure(t *testing.T) {
	refTime := time.Now()

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetPipelines", AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string")).
		Return(&Pipelines{
			HeadCommit: "sha",
			Runs: []Run{
				{
					ID:         10,
					Status:     "failed",
					StartedAt:  ToTime(refTime.Add(-time.Second * 30)),
					FinishedAt: ToTime(refTime.Add(-time.Second * 15)),
				},
			},
		}, nil)
	mockRepository.On("GetCommit", AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string")).
		Return(&Commit{
			Author: &models.Author{
				Name:      "test",
				AvatarURL: "https://test.example.com",
			},
		}, nil)

	gu := NewGitlabUsecase(mockRepository)

	expected := models.NewTile(gitlab.GitlabPipelinesTileType).WithBuild()
	expected.Label = "test"
	expected.Build.Branch = ToString("master")

	expected.Status = models.FailedStatus
	expected.Build.PreviousStatus = models.UnknownStatus
	expected.Build.StartedAt = ToTime(refTime.Add(-time.Second * 30))
	expected.Build.FinishedAt = ToTime(refTime.Add(-time.Second * 15))
	expected.Build.Author = &models.Author{
		Name:      "test",
		AvatarURL: "https://test.example.com",
	}

	tile, err := gu.Pipelines(&PipelinesParams{Repository: "test", Ref: "master"})
	if assert.NoError(t, err) {
		assert.NotNil(t, tile)
		assert.Equal(t, expected, tile)

		mockRepository.AssertNumberOfCalls(t, "GetPipelines", 1)
		mockRepository.AssertNumberOfCalls(t, "GetCommit", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestPipelines_Queued(t *testing.T) {
	refTime := time.Now()

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetPipelines", AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string")).
		Return(&Pipelines{
			HeadCommit: "sha",
			Runs: []Run{
				{
					ID:        10,
					Status:    "pending",
					StartedAt: ToTime(refTime.Add(-time.Second * 30)),
				},
			},
		}, nil)

	gu := NewGitlabUsecase(mockRepository)

	expected := models.NewTile(gitlab.GitlabPipelinesTileType).WithBuild()
	expected.Label = "test"
	expected.Build.Branch = ToString("master")

	expected.Status = models.QueuedStatus
	expected.Build.PreviousStatus = models.UnknownStatus
	expected.Build.StartedAt = ToTime(refTime.Add(-time.Second * 30))

	tile, err := gu.Pipelines(&PipelinesParams{Repository: "test", Ref: "master"})
	if assert.NoError(t, err) {
		assert.NotNil(t, tile)
		assert.Equal(t, expected, tile)

		mockRepository.AssertNumberOfCalls(t, "GetPipelines", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestPipelines_Running(t *testing.T) {
	refTime := time.Now()

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetPipelines", AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string")).
		Return(&Pipelines{
			HeadCommit: "sha",
			Runs: []Run{
				{
					ID:        10,
					Status:    "running",
					StartedAt: ToTime(refTime.Add(-time.Second * 30)),
					Duration:  30,
				},
			},
		}, nil)

	gu := NewGitlabUsecase(mockRepository)
	gUsecase, ok := gu.(*gitlabUsecase)
	if assert.True(t, ok) {
		expected := models.NewTile(gitlab.GitlabPipelinesTileType).WithBuild()
		expected.Label = "test"
		expected.Build.Branch = ToString("master")

		expected.Status = models.RunningStatus
		expected.Build.PreviousStatus = models.UnknownStatus
		expected.Build.StartedAt = ToTime(refTime.Add(-time.Second * 30))
		expected.Build.Duration = ToInt64(int64(30))
		expected.Build.EstimatedDuration = ToInt64(int64(0))

		tile, err := gUsecase.Pipelines(&PipelinesParams{Repository: "test", Ref: "master"})
		if assert.NoError(t, err) {
			assert.NotNil(t, tile)
			assert.Equal(t, expected, tile)
		}

		params := &PipelinesParams{Repository: "test", Ref: "master"}
		gUsecase.buildsCache.Add(params, hash.GetMD5Hash("10"), models.SuccessStatus, time.Second*120)

		expected.Build.PreviousStatus = models.SuccessStatus
		expected.Build.EstimatedDuration = ToInt64(int64(120))

		tile, err = gUsecase.Pipelines(&PipelinesParams{Repository: "test", Ref: "master"})
		if assert.NoError(t, err) {
			assert.NotNil(t, tile)
			assert.Equal(t, expected, tile)

			mockRepository.AssertNumberOfCalls(t, "GetPipelines", 2)
			mockRepository.AssertExpectations(t)
		}
	}
}

func TestListDynamicTile_Error(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetMergeRequests", AnythingOfType("string"), AnythingOfType("string")).
		Return(nil, errors.New("boom"))

	gu := NewGitlabUsecase(mockRepository)

	results, err := gu.ListDynamicTile(&MergeRequestParams{Repository: "test"})
	if assert.Error(t, err) {
		assert.Nil(t, results)
		assert.IsType(t, &models.MonitororError{}, err)
		assert.Equal(t, "unable to find merge request", err.Error())
		mockRepository.AssertNumberOfCalls(t, "GetMergeRequests", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestListDynamicTile_Success(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetMergeRequests", AnythingOfType("string"), AnythingOfType("string")).
		Return([]MergeRequest{
			{
				ID:         2,
				Repository: "test",
				Ref:        "master",
			},
		}, nil)

	gu := NewGitlabUsecase(mockRepository)

	results, err := gu.ListDynamicTile(&MergeRequestParams{Repository: "test"})
	if assert.NoError(t, err) {
		assert.NotNil(t, results)
		assert.Len(t, results, 1)
		assert.Equal(t, gitlab.GitlabPipelinesTileType, results[0].TileType)
		assert.Equal(t, "MR#2 @ test", results[0].Label)
		assert.Equal(t, "test", results[0].Params["repository"])
		assert.Equal(t, "master", results[0].Params["ref"])

		mockRepository.AssertNumberOfCalls(t, "GetMergeRequests", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestComputePipelines_Status(t *testing.T) {
	for index, testcase := range []struct {
		runs           []Run
		expectedStatus models.TileStatus
	}{
		{
			runs:           []Run{},
			expectedStatus: models.UnknownStatus,
		},
		{
			runs: []Run{
				{Status: "success"},
			},
			expectedStatus: models.SuccessStatus,
		},
		{
			runs: []Run{
				{Status: "failed"},
			},
			expectedStatus: models.FailedStatus,
		},
		{
			runs: []Run{
				{Status: "success"},
				{Status: "running", CreatedAt: time.Now()},
			},
			expectedStatus: models.RunningStatus,
		},
		{
			runs: []Run{
				{Status: "pending"},
			},
			expectedStatus: models.QueuedStatus,
		},
		{
			runs: []Run{
				{Status: "running"},
				{Status: "success"},
			},
			expectedStatus: models.RunningStatus,
		},
		{
			runs: []Run{
				{Status: "canceled"},
			},
			expectedStatus: models.CanceledStatus,
		},
		{
			runs: []Run{
				{Status: "skipped"},
			},
			expectedStatus: models.DisabledStatus,
		},
	} {
		status, _, _, _, _ := computePipelines(&Pipelines{Runs: testcase.runs})
		assert.Equal(t, testcase.expectedStatus, status, fmt.Sprintf("test %d failed", index))
	}
}

func TestComputePipelines_Time(t *testing.T) {
	expectedStartedAt := time.Now().Add(-time.Minute).Truncate(time.Minute)
	expectedFinishedAt := time.Now().Add(+time.Minute).Truncate(time.Minute)

	pipelines := &Pipelines{
		Runs: []Run{
			{StartedAt: ToTime(time.Now()), FinishedAt: ToTime(time.Now())},
			{StartedAt: ToTime(time.Now().Add(-time.Minute * 1)), FinishedAt: ToTime(time.Now().Add(+time.Minute * 1))},
		},
	}

	_, startedAt, finishedAt, _, _ := computePipelines(pipelines)
	assert.Equal(t, expectedStartedAt, (*startedAt).Truncate(time.Minute))
	assert.Equal(t, expectedFinishedAt, (*finishedAt).Truncate(time.Minute))
}

func TestComputePipelines_ID(t *testing.T) {
	pipelines := &Pipelines{
		Runs: []Run{
			{ID: 12},
			{ID: 13},
		},
	}

	_, _, _, _, id := computePipelines(pipelines)
	assert.Equal(t, "580335ece448f965fb1e254ee96d4cff", id)
}
