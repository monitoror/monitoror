package usecase

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/github"
	"github.com/monitoror/monitoror/monitorable/github/mocks"
	. "github.com/monitoror/monitoror/monitorable/github/models"
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

	gu := NewGithubUsecase(mockRepository)

	expected := models.NewTile(github.GithubCountTileType).WithValue(models.NumberUnit)
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

func TestChecks_Error(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetChecks", AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string")).
		Return(nil, errors.New("boom"))

	gu := NewGithubUsecase(mockRepository)

	tile, err := gu.Checks(&ChecksParams{Owner: "test", Repository: "test", Ref: "master"})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &models.MonitororError{}, err)
		assert.Equal(t, "unable to find ref checks", err.Error())
		mockRepository.AssertNumberOfCalls(t, "GetChecks", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestChecks_NoChecks(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetChecks", AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string")).
		Return(&Checks{}, nil)

	gu := NewGithubUsecase(mockRepository)

	tile, err := gu.Checks(&ChecksParams{Owner: "test", Repository: "test", Ref: "master"})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &models.MonitororError{}, err)
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
		Return(&Checks{
			Runs: []Run{
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

	expected := models.NewTile(github.GithubChecksTileType).WithBuild()
	expected.Label = "test"
	expected.Build.Branch = ToString("master")

	expected.Status = models.SuccessStatus
	expected.Build.PreviousStatus = models.UnknownStatus
	expected.Build.StartedAt = ToTime(startedAt)
	expected.Build.FinishedAt = ToTime(finishedAt)

	tile, err := gu.Checks(&ChecksParams{Owner: "test", Repository: "test", Ref: "master"})
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
		Return(&Checks{
			HeadCommit: ToString("sha"),
			Runs: []Run{
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
		Return(&Commit{
			Author: &models.Author{
				Name:      "test",
				AvatarURL: "https://test.example.com",
			},
		}, nil)

	gu := NewGithubUsecase(mockRepository)

	expected := models.NewTile(github.GithubChecksTileType).WithBuild()
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

	tile, err := gu.Checks(&ChecksParams{Owner: "test", Repository: "test", Ref: "master"})
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
		Return(&Checks{
			HeadCommit: ToString("sha"),
			Runs: []Run{
				{
					ID:        10,
					Status:    "queued",
					StartedAt: ToTime(refTime.Add(-time.Second * 30)),
				},
			},
		}, nil)

	gu := NewGithubUsecase(mockRepository)

	expected := models.NewTile(github.GithubChecksTileType).WithBuild()
	expected.Label = "test"
	expected.Build.Branch = ToString("master")

	expected.Status = models.QueuedStatus
	expected.Build.PreviousStatus = models.UnknownStatus
	expected.Build.StartedAt = ToTime(refTime.Add(-time.Second * 30))

	tile, err := gu.Checks(&ChecksParams{Owner: "test", Repository: "test", Ref: "master"})
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
		Return(&Checks{
			HeadCommit: ToString("sha"),
			Runs: []Run{
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
		expected := models.NewTile(github.GithubChecksTileType).WithBuild()
		expected.Label = "test"
		expected.Build.Branch = ToString("master")

		expected.Status = models.RunningStatus
		expected.Build.PreviousStatus = models.UnknownStatus
		expected.Build.StartedAt = ToTime(refTime.Add(-time.Second * 30))
		expected.Build.Duration = ToInt64(int64(30))
		expected.Build.EstimatedDuration = ToInt64(int64(0))

		tile, err := gUsecase.Checks(&ChecksParams{Owner: "test", Repository: "test", Ref: "master"})
		if assert.NoError(t, err) {
			assert.NotNil(t, tile)
			assert.Equal(t, expected, tile)
		}

		params := &ChecksParams{Owner: "test", Repository: "test", Ref: "master"}
		gUsecase.buildsCache.Add(params, hash.GetMD5Hash("10"), models.SuccessStatus, time.Second*120)

		expected.Build.PreviousStatus = models.SuccessStatus
		expected.Build.EstimatedDuration = ToInt64(int64(120))

		tile, err = gUsecase.Checks(&ChecksParams{Owner: "test", Repository: "test", Ref: "master"})
		if assert.NoError(t, err) {
			assert.NotNil(t, tile)
			assert.Equal(t, expected, tile)

			mockRepository.AssertNumberOfCalls(t, "GetChecks", 2)
			mockRepository.AssertExpectations(t)
		}
	}
}

func TestListDynamicTile_Error(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetPullRequests", AnythingOfType("string"), AnythingOfType("string")).
		Return(nil, errors.New("boom"))

	gu := NewGithubUsecase(mockRepository)

	results, err := gu.ListDynamicTile(&PullRequestParams{Owner: "test", Repository: "test"})
	if assert.Error(t, err) {
		assert.Nil(t, results)
		assert.IsType(t, &models.MonitororError{}, err)
		assert.Equal(t, "unable to find pull request", err.Error())
		mockRepository.AssertNumberOfCalls(t, "GetPullRequests", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestListDynamicTile_Success(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetPullRequests", AnythingOfType("string"), AnythingOfType("string")).
		Return([]PullRequest{
			{
				ID:         2,
				Owner:      "test",
				Repository: "test",
				Ref:        "master",
			},
		}, nil)

	gu := NewGithubUsecase(mockRepository)

	results, err := gu.ListDynamicTile(&PullRequestParams{Owner: "test", Repository: "test"})
	if assert.NoError(t, err) {
		assert.NotNil(t, results)
		assert.Len(t, results, 1)
		assert.Equal(t, github.GithubChecksTileType, results[0].TileType)
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
		runs           []Run
		statuses       []Status
		expectedStatus models.TileStatus
	}{
		{
			runs:           []Run{},
			statuses:       []Status{},
			expectedStatus: models.UnknownStatus,
		},
		{
			runs: []Run{
				{Status: "completed", Conclusion: "success"},
			},
			statuses: []Status{
				{State: "success"},
			},
			expectedStatus: models.SuccessStatus,
		},
		{
			runs: []Run{
				{Status: "completed", Conclusion: "success"},
			},
			statuses: []Status{
				{Title: "test1", State: "success"},
				{Title: "test2", State: "failure", CreatedAt: time.Now()},
				{Title: "test2", State: "pending", CreatedAt: time.Now().Add(-time.Minute)}, // will be removed because title is duplicated
			},
			expectedStatus: models.FailedStatus,
		},
		{
			runs: []Run{
				{Status: "completed", Conclusion: "success"},
				{Status: "in_progress"},
			},
			statuses: []Status{
				{State: "success"},
				{State: "success"},
			},
			expectedStatus: models.RunningStatus,
		},
		{
			runs: []Run{
				{Status: "completed", Conclusion: "success"},
			},
			statuses: []Status{
				{Title: "test1", State: "success"},
				{Title: "test2", State: "pending"},
			},
			expectedStatus: models.RunningStatus,
		},
		{
			runs: []Run{
				{Status: "queued"},
			},
			statuses: []Status{
				{State: "success"},
			},
			expectedStatus: models.QueuedStatus,
		},
		{
			runs: []Run{
				{Status: "queued"},
			},
			statuses: []Status{
				{State: "error"},
			},
			expectedStatus: models.FailedStatus,
		},
		{
			runs: []Run{
				{Status: "queued"},
				{Status: "completed", Conclusion: "success"},
			},
			statuses: []Status{
				{State: "pending"},
			},
			expectedStatus: models.RunningStatus,
		},
		{
			runs: []Run{
				{Status: "completed", Conclusion: "timed_out"},
			},
			statuses: []Status{
				{State: "success"},
			},
			expectedStatus: models.FailedStatus,
		},
		{
			runs: []Run{
				{Status: "completed", Conclusion: "failure"},
			},
			statuses: []Status{
				{Title: "test1", State: "success"},
				{Title: "test2", State: "pending"},
			},
			expectedStatus: models.RunningStatus,
		},
		{
			runs: []Run{
				{Status: "completed", Conclusion: "cancelled"},
			},
			statuses: []Status{
				{State: "success"},
			},
			expectedStatus: models.CanceledStatus,
		},
		{
			runs: []Run{
				{Status: "completed", Conclusion: "neutral"},
			},
			statuses: []Status{
				{State: "success"},
			},
			expectedStatus: models.WarningStatus,
		},
		{
			runs: []Run{
				{Status: "completed", Conclusion: "action_required"},
			},
			statuses: []Status{
				{State: "success"},
			},
			expectedStatus: models.ActionRequiredStatus,
		},
	} {
		status, _, _, _ := computeChecks(&Checks{Runs: testcase.runs, Statuses: testcase.statuses})
		assert.Equal(t, testcase.expectedStatus, status, fmt.Sprintf("test %d failed", index))
	}
}

func TestComputeRefStatus_Time(t *testing.T) {
	expectedStartedAt := time.Now().Add(-time.Minute * 2)
	expectedFinishedAt := time.Now().Add(+time.Minute * 2)

	refStatus := &Checks{
		Runs: []Run{
			{StartedAt: ToTime(time.Now()), CompletedAt: ToTime(time.Now())},
			{StartedAt: ToTime(time.Now().Add(-time.Minute * 1)), CompletedAt: ToTime(time.Now().Add(+time.Minute * 1))},
		},
		Statuses: []Status{
			{CreatedAt: expectedStartedAt, UpdatedAt: expectedFinishedAt},
		},
	}

	_, startedAt, finishedAt, _ := computeChecks(refStatus)
	assert.Equal(t, expectedStartedAt, *startedAt)
	assert.Equal(t, expectedFinishedAt, *finishedAt)
}

func TestComputeRefStatus_ID(t *testing.T) {

	refStatus := &Checks{
		Runs: []Run{
			{ID: 12},
			{ID: 13},
		},
		Statuses: []Status{
			{ID: 137},
		},
	}

	_, _, _, id := computeChecks(refStatus)
	assert.Equal(t, "b103a99f8ef3da68771355b76aa05ccf", id)
}
