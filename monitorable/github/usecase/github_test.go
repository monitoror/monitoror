package usecase

import (
	"errors"
	"fmt"
	"testing"
	"time"

	. "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/github"
	"github.com/monitoror/monitoror/monitorable/github/mocks"
	"github.com/monitoror/monitoror/monitorable/github/models"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/hash"

	. "github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func TestIssues_Error(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetIssuesCount", AnythingOfType("string")).
		Return(0, errors.New("boom"))

	gu := NewGithubUsecase(mockRepository)

	tile, err := gu.Issues(&models.IssuesParams{Query: "test"})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &MonitororError{}, err)
		assert.Equal(t, "unable to find issues count or wrong query", err.Error())
		mockRepository.AssertNumberOfCalls(t, "GetIssuesCount", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestIssues_Success(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetIssuesCount", AnythingOfType("string")).
		Return(10, nil)

	gu := NewGithubUsecase(mockRepository)

	expected := NewTile(github.GithubIssuesTileType)
	expected.Label = "test"
	expected.Status = SuccessStatus
	expected.Values = []float64{10}

	tile, err := gu.Issues(&models.IssuesParams{Query: "test"})
	if assert.NoError(t, err) {
		assert.NotNil(t, tile)
		assert.Equal(t, expected, tile)

		mockRepository.AssertNumberOfCalls(t, "GetIssuesCount", 1)
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
		assert.IsType(t, &MonitororError{}, err)
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
		assert.IsType(t, &MonitororError{}, err)
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

	expected := NewTile(github.GithubChecksTileType)
	expected.Label = "test\n@master"
	expected.Status = SuccessStatus
	expected.PreviousStatus = UnknownStatus
	expected.StartedAt = ToTime(startedAt)
	expected.FinishedAt = ToTime(finishedAt)

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
			Author: &models.Author{
				Name:      "test",
				AvatarURL: "https://test.example.com",
			},
		}, nil)

	gu := NewGithubUsecase(mockRepository)

	expected := NewTile(github.GithubChecksTileType)
	expected.Label = "test\n@master"
	expected.Status = FailedStatus
	expected.PreviousStatus = UnknownStatus
	expected.StartedAt = ToTime(refTime.Add(-time.Second * 30))
	expected.FinishedAt = ToTime(refTime.Add(-time.Second * 15))
	expected.Author = &Author{
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

	expected := NewTile(github.GithubChecksTileType)
	expected.Label = "test\n@master"
	expected.Status = QueuedStatus
	expected.PreviousStatus = UnknownStatus
	expected.StartedAt = ToTime(refTime.Add(-time.Second * 30))

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
		expected := NewTile(github.GithubChecksTileType)
		expected.Label = "test\n@master"
		expected.Status = RunningStatus
		expected.PreviousStatus = UnknownStatus
		expected.StartedAt = ToTime(refTime.Add(-time.Second * 30))
		expected.Duration = ToInt64(int64(30))
		expected.EstimatedDuration = ToInt64(int64(0))

		tile, err := gUsecase.Checks(&models.ChecksParams{Owner: "test", Repository: "test", Ref: "master"})
		if assert.NoError(t, err) {
			assert.NotNil(t, tile)
			assert.Equal(t, expected, tile)
		}

		params := &models.ChecksParams{Owner: "test", Repository: "test", Ref: "master"}
		gUsecase.buildsCache.Add(params, hash.GetMD5Hash("10"), SuccessStatus, time.Second*120)

		expected.PreviousStatus = SuccessStatus
		expected.EstimatedDuration = ToInt64(int64(120))

		tile, err = gUsecase.Checks(&models.ChecksParams{Owner: "test", Repository: "test", Ref: "master"})
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

	results, err := gu.ListDynamicTile(&models.PullRequestParams{Owner: "test", Repository: "test"})
	if assert.Error(t, err) {
		assert.Nil(t, results)
		assert.IsType(t, &MonitororError{}, err)
		assert.Equal(t, "unable to find pull request", err.Error())
		mockRepository.AssertNumberOfCalls(t, "GetPullRequests", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestListDynamicTile_Success(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetPullRequests", AnythingOfType("string"), AnythingOfType("string")).
		Return([]models.PullRequest{
			{
				Title:      "PR#2 - TEST",
				Owner:      "test",
				Repository: "test",
				Ref:        "master",
			},
		}, nil)

	gu := NewGithubUsecase(mockRepository)

	results, err := gu.ListDynamicTile(&models.PullRequestParams{Owner: "test", Repository: "test"})
	if assert.NoError(t, err) {
		assert.NotNil(t, results)
		assert.Len(t, results, 1)
		assert.Equal(t, github.GithubChecksTileType, results[0].TileType)
		assert.Equal(t, "test\nPR#2 - TEST", results[0].Label)
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
		expectedStatus TileStatus
	}{
		{
			runs:           []models.Run{},
			statuses:       []models.Status{},
			expectedStatus: UnknownStatus,
		},
		{
			runs: []models.Run{
				{Status: "completed", Conclusion: "success"},
			},
			statuses: []models.Status{
				{State: "success"},
			},
			expectedStatus: SuccessStatus,
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
			expectedStatus: FailedStatus,
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
			expectedStatus: RunningStatus,
		},
		{
			runs: []models.Run{
				{Status: "completed", Conclusion: "success"},
			},
			statuses: []models.Status{
				{Title: "test1", State: "success"},
				{Title: "test2", State: "pending"},
			},
			expectedStatus: RunningStatus,
		},
		{
			runs: []models.Run{
				{Status: "queued"},
			},
			statuses: []models.Status{
				{State: "success"},
			},
			expectedStatus: QueuedStatus,
		},
		{
			runs: []models.Run{
				{Status: "queued"},
			},
			statuses: []models.Status{
				{State: "error"},
			},
			expectedStatus: FailedStatus,
		},
		{
			runs: []models.Run{
				{Status: "queued"},
				{Status: "completed", Conclusion: "success"},
			},
			statuses: []models.Status{
				{State: "pending"},
			},
			expectedStatus: RunningStatus,
		},
		{
			runs: []models.Run{
				{Status: "completed", Conclusion: "timed_out"},
			},
			statuses: []models.Status{
				{State: "success"},
			},
			expectedStatus: FailedStatus,
		},
		{
			runs: []models.Run{
				{Status: "completed", Conclusion: "failure"},
			},
			statuses: []models.Status{
				{Title: "test1", State: "success"},
				{Title: "test2", State: "pending"},
			},
			expectedStatus: RunningStatus,
		},
		{
			runs: []models.Run{
				{Status: "completed", Conclusion: "cancelled"},
			},
			statuses: []models.Status{
				{State: "success"},
			},
			expectedStatus: CanceledStatus,
		},
		{
			runs: []models.Run{
				{Status: "completed", Conclusion: "neutral"},
			},
			statuses: []models.Status{
				{State: "success"},
			},
			expectedStatus: WarningStatus,
		},
		{
			runs: []models.Run{
				{Status: "completed", Conclusion: "action_required"},
			},
			statuses: []models.Status{
				{State: "success"},
			},
			expectedStatus: ActionRequiredStatus,
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
