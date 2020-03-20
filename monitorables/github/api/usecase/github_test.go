package usecase

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/monitoror/monitoror/monitorables/github/api"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/github/api/mocks"
	"github.com/monitoror/monitoror/monitorables/github/api/models"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/hash"

	. "github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func TestCount_Error(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetCount", AnythingOfType("string")).
		Return(0, errors.New("boom"))

	gu := NewGithubUsecase(mockRepository)

	tile, err := gu.Count(&models.CountParams{Query: "test"})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &coreModels.MonitororError{}, err)
		assert.Equal(t, "unable to find count or wrong query", err.Error())
		mockRepository.AssertNumberOfCalls(t, "GetCount", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestCount_Success(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetCount", AnythingOfType("string")).
		Return(10, nil)

	gu := NewGithubUsecase(mockRepository)

	expected := coreModels.NewTile(api.GithubCountTileType).WithValue(coreModels.NumberUnit)
	expected.Label = "test"
	expected.Status = coreModels.SuccessStatus
	expected.Value.Values = []string{"10"}

	tile, err := gu.Count(&models.CountParams{Query: "test"})
	if assert.NoError(t, err) {
		assert.NotNil(t, tile)
		assert.Equal(t, expected, tile)

		mockRepository.AssertNumberOfCalls(t, "GetCount", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestChecks_Error(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetChecks", AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string")).
		Return(nil, errors.New("boom"))

	gu := NewGithubUsecase(mockRepository)

	tile, err := gu.Checks(&models.ChecksParams{Owner: "test", Repository: "test", Ref: "master"})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &coreModels.MonitororError{}, err)
		assert.Equal(t, "unable to find ref checks", err.Error())
		mockRepository.AssertNumberOfCalls(t, "GetChecks", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestChecks_NoChecks(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetChecks", AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string")).
		Return(&models.Checks{}, nil)

	gu := NewGithubUsecase(mockRepository)

	tile, err := gu.Checks(&models.ChecksParams{Owner: "test", Repository: "test", Ref: "master"})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &coreModels.MonitororError{}, err)
		assert.Equal(t, "no ref checks found", err.Error())
		mockRepository.AssertNumberOfCalls(t, "GetChecks", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestChecks_Success(t *testing.T) {
	refTime := time.Now()
	startedAt := refTime.Add(-time.Second * 30)
	finishedAt := refTime.Add(-time.Second * 15)

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetChecks", AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string")).
		Return(&models.Checks{
			Runs: []models.Run{
				{
					ID:          10,
					Status:      "completed",
					Conclusion:  "success",
					StartedAt:   ToTime(startedAt),
					CompletedAt: ToTime(finishedAt),
				},
			},
		}, nil)

	gu := NewGithubUsecase(mockRepository)

	expected := coreModels.NewTile(api.GithubChecksTileType).WithBuild()
	expected.Label = "test"
	expected.Build.Branch = ToString("master")

	expected.Status = coreModels.SuccessStatus
	expected.Build.PreviousStatus = coreModels.UnknownStatus
	expected.Build.StartedAt = ToTime(startedAt)
	expected.Build.FinishedAt = ToTime(finishedAt)

	tile, err := gu.Checks(&models.ChecksParams{Owner: "test", Repository: "test", Ref: "master"})
	if assert.NoError(t, err) {
		assert.NotNil(t, tile)
		assert.Equal(t, expected, tile)

		mockRepository.AssertNumberOfCalls(t, "GetChecks", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestChecks_Failure(t *testing.T) {
	refTime := time.Now()

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetChecks", AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string")).
		Return(&models.Checks{
			HeadCommit: ToString("sha"),
			Runs: []models.Run{
				{
					ID:          10,
					Status:      "completed",
					Conclusion:  "failure",
					StartedAt:   ToTime(refTime.Add(-time.Second * 30)),
					CompletedAt: ToTime(refTime.Add(-time.Second * 15)),
				},
			},
		}, nil)
	mockRepository.On("GetCommit", AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string")).
		Return(&models.Commit{
			Author: &coreModels.Author{
				Name:      "test",
				AvatarURL: "https://test.example.com",
			},
		}, nil)

	gu := NewGithubUsecase(mockRepository)

	expected := coreModels.NewTile(api.GithubChecksTileType).WithBuild()
	expected.Label = "test"
	expected.Build.Branch = ToString("master")

	expected.Status = coreModels.FailedStatus
	expected.Build.PreviousStatus = coreModels.UnknownStatus
	expected.Build.StartedAt = ToTime(refTime.Add(-time.Second * 30))
	expected.Build.FinishedAt = ToTime(refTime.Add(-time.Second * 15))
	expected.Build.Author = &coreModels.Author{
		Name:      "test",
		AvatarURL: "https://test.example.com",
	}

	tile, err := gu.Checks(&models.ChecksParams{Owner: "test", Repository: "test", Ref: "master"})
	if assert.NoError(t, err) {
		assert.NotNil(t, tile)
		assert.Equal(t, expected, tile)

		mockRepository.AssertNumberOfCalls(t, "GetChecks", 1)
		mockRepository.AssertNumberOfCalls(t, "GetCommit", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestChecks_Queued(t *testing.T) {
	refTime := time.Now()

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetChecks", AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string")).
		Return(&models.Checks{
			HeadCommit: ToString("sha"),
			Runs: []models.Run{
				{
					ID:        10,
					Status:    "queued",
					StartedAt: ToTime(refTime.Add(-time.Second * 30)),
				},
			},
		}, nil)

	gu := NewGithubUsecase(mockRepository)

	expected := coreModels.NewTile(api.GithubChecksTileType).WithBuild()
	expected.Label = "test"
	expected.Build.Branch = ToString("master")

	expected.Status = coreModels.QueuedStatus
	expected.Build.PreviousStatus = coreModels.UnknownStatus
	expected.Build.StartedAt = ToTime(refTime.Add(-time.Second * 30))

	tile, err := gu.Checks(&models.ChecksParams{Owner: "test", Repository: "test", Ref: "master"})
	if assert.NoError(t, err) {
		assert.NotNil(t, tile)
		assert.Equal(t, expected, tile)

		mockRepository.AssertNumberOfCalls(t, "GetChecks", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestChecks_Running(t *testing.T) {
	refTime := time.Now()

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetChecks", AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string")).
		Return(&models.Checks{
			HeadCommit: ToString("sha"),
			Runs: []models.Run{
				{
					ID:        10,
					Status:    "in_progress",
					StartedAt: ToTime(refTime.Add(-time.Second * 30)),
				},
			},
		}, nil)

	gu := NewGithubUsecase(mockRepository)
	gUsecase, ok := gu.(*githubUsecase)
	if assert.True(t, ok) {
		expected := coreModels.NewTile(api.GithubChecksTileType).WithBuild()
		expected.Label = "test"
		expected.Build.Branch = ToString("master")

		expected.Status = coreModels.RunningStatus
		expected.Build.PreviousStatus = coreModels.UnknownStatus
		expected.Build.StartedAt = ToTime(refTime.Add(-time.Second * 30))
		expected.Build.Duration = ToInt64(int64(30))
		expected.Build.EstimatedDuration = ToInt64(int64(0))

		tile, err := gUsecase.Checks(&models.ChecksParams{Owner: "test", Repository: "test", Ref: "master"})
		if assert.NoError(t, err) {
			assert.NotNil(t, tile)
			assert.Equal(t, expected, tile)
		}

		params := &models.ChecksParams{Owner: "test", Repository: "test", Ref: "master"}
		gUsecase.buildsCache.Add(params, hash.GetMD5Hash("10"), coreModels.SuccessStatus, time.Second*120)

		expected.Build.PreviousStatus = coreModels.SuccessStatus
		expected.Build.EstimatedDuration = ToInt64(int64(120))

		tile, err = gUsecase.Checks(&models.ChecksParams{Owner: "test", Repository: "test", Ref: "master"})
		if assert.NoError(t, err) {
			assert.NotNil(t, tile)
			assert.Equal(t, expected, tile)

			mockRepository.AssertNumberOfCalls(t, "GetChecks", 2)
			mockRepository.AssertExpectations(t)
		}
	}
}

func TestPullRequests_Error(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetPullRequests", AnythingOfType("string"), AnythingOfType("string")).
		Return(nil, errors.New("boom"))

	gu := NewGithubUsecase(mockRepository)

	results, err := gu.PullRequests(&models.PullRequestParams{Owner: "test", Repository: "test"})
	if assert.Error(t, err) {
		assert.Nil(t, results)
		assert.IsType(t, &coreModels.MonitororError{}, err)
		assert.Equal(t, "unable to find pull request", err.Error())
		mockRepository.AssertNumberOfCalls(t, "GetPullRequests", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestPullRequests_Success(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetPullRequests", AnythingOfType("string"), AnythingOfType("string")).
		Return([]models.PullRequest{
			{
				ID:         2,
				Owner:      "test",
				Repository: "test",
				Ref:        "master",
			},
		}, nil)

	gu := NewGithubUsecase(mockRepository)

	results, err := gu.PullRequests(&models.PullRequestParams{Owner: "test", Repository: "test"})
	if assert.NoError(t, err) {
		assert.NotNil(t, results)
		assert.Len(t, results, 1)
		assert.Equal(t, api.GithubChecksTileType, results[0].TileType)
		assert.Equal(t, "PR#2 @ test", results[0].Label)
		assert.Equal(t, "test", results[0].Params["owner"])
		assert.Equal(t, "test", results[0].Params["repository"])
		assert.Equal(t, "master", results[0].Params["ref"])

		mockRepository.AssertNumberOfCalls(t, "GetPullRequests", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestComputeRefStatus_Status(t *testing.T) {
	for index, testcase := range []struct {
		runs           []models.Run
		statuses       []models.Status
		expectedStatus coreModels.TileStatus
	}{
		{
			runs:           []models.Run{},
			statuses:       []models.Status{},
			expectedStatus: coreModels.UnknownStatus,
		},
		{
			runs: []models.Run{
				{Status: "completed", Conclusion: "success"},
			},
			statuses: []models.Status{
				{State: "success"},
			},
			expectedStatus: coreModels.SuccessStatus,
		},
		{
			runs: []models.Run{
				{Status: "completed", Conclusion: "success"},
			},
			statuses: []models.Status{
				{Title: "test1", State: "success"},
				{Title: "test2", State: "failure", CreatedAt: time.Now()},
				{Title: "test2", State: "pending", CreatedAt: time.Now().Add(-time.Minute)}, // will be removed because title is duplicated
			},
			expectedStatus: coreModels.FailedStatus,
		},
		{
			runs: []models.Run{
				{Status: "completed", Conclusion: "success"},
				{Status: "in_progress"},
			},
			statuses: []models.Status{
				{State: "success"},
				{State: "success"},
			},
			expectedStatus: coreModels.RunningStatus,
		},
		{
			runs: []models.Run{
				{Status: "completed", Conclusion: "success"},
			},
			statuses: []models.Status{
				{Title: "test1", State: "success"},
				{Title: "test2", State: "pending"},
			},
			expectedStatus: coreModels.RunningStatus,
		},
		{
			runs: []models.Run{
				{Status: "queued"},
			},
			statuses: []models.Status{
				{State: "success"},
			},
			expectedStatus: coreModels.QueuedStatus,
		},
		{
			runs: []models.Run{
				{Status: "queued"},
			},
			statuses: []models.Status{
				{State: "error"},
			},
			expectedStatus: coreModels.FailedStatus,
		},
		{
			runs: []models.Run{
				{Status: "queued"},
				{Status: "completed", Conclusion: "success"},
			},
			statuses: []models.Status{
				{State: "pending"},
			},
			expectedStatus: coreModels.RunningStatus,
		},
		{
			runs: []models.Run{
				{Status: "completed", Conclusion: "timed_out"},
			},
			statuses: []models.Status{
				{State: "success"},
			},
			expectedStatus: coreModels.FailedStatus,
		},
		{
			runs: []models.Run{
				{Status: "completed", Conclusion: "failure"},
			},
			statuses: []models.Status{
				{Title: "test1", State: "success"},
				{Title: "test2", State: "pending"},
			},
			expectedStatus: coreModels.RunningStatus,
		},
		{
			runs: []models.Run{
				{Status: "completed", Conclusion: "cancelled"},
			},
			statuses: []models.Status{
				{State: "success"},
			},
			expectedStatus: coreModels.CanceledStatus,
		},
		{
			runs: []models.Run{
				{Status: "completed", Conclusion: "neutral"},
			},
			statuses: []models.Status{
				{State: "success"},
			},
			expectedStatus: coreModels.WarningStatus,
		},
		{
			runs: []models.Run{
				{Status: "completed", Conclusion: "action_required"},
			},
			statuses: []models.Status{
				{State: "success"},
			},
			expectedStatus: coreModels.ActionRequiredStatus,
		},
	} {
		status, _, _, _ := computeChecks(&models.Checks{Runs: testcase.runs, Statuses: testcase.statuses})
		assert.Equal(t, testcase.expectedStatus, status, fmt.Sprintf("test %d failed", index))
	}
}

func TestComputeRefStatus_Time(t *testing.T) {
	expectedStartedAt := time.Now().Add(-time.Minute * 2)
	expectedFinishedAt := time.Now().Add(+time.Minute * 2)

	refStatus := &models.Checks{
		Runs: []models.Run{
			{StartedAt: ToTime(time.Now()), CompletedAt: ToTime(time.Now())},
			{StartedAt: ToTime(time.Now().Add(-time.Minute * 1)), CompletedAt: ToTime(time.Now().Add(+time.Minute * 1))},
		},
		Statuses: []models.Status{
			{CreatedAt: expectedStartedAt, UpdatedAt: expectedFinishedAt},
		},
	}

	_, startedAt, finishedAt, _ := computeChecks(refStatus)
	assert.Equal(t, expectedStartedAt, *startedAt)
	assert.Equal(t, expectedFinishedAt, *finishedAt)
}

func TestComputeRefStatus_ID(t *testing.T) {

	refStatus := &models.Checks{
		Runs: []models.Run{
			{ID: 12},
			{ID: 13},
		},
		Statuses: []models.Status{
			{ID: 137},
		},
	}

	_, _, _, id := computeChecks(refStatus)
	assert.Equal(t, "b103a99f8ef3da68771355b76aa05ccf", id)
}
