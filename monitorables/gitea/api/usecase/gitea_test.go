package usecase

import (
	"errors"
	"testing"
	"time"

	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/gitea/api"
	"github.com/monitoror/monitoror/monitorables/gitea/api/mocks"
	"github.com/monitoror/monitoror/monitorables/gitea/api/models"
	"github.com/monitoror/monitoror/pkg/hash"

	. "github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func TestCount_Error(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetCount", AnythingOfType("*models.CountParams")).
		Return(0, errors.New("boom"))

	gu := NewGiteaUsecase(mockRepository)

	tile, err := gu.Count(&models.CountParams{Owner: "test", Repository: "test"})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &coreModels.MonitororError{}, err)
		assert.Equal(t, "unable to load count", err.Error())
		mockRepository.AssertNumberOfCalls(t, "GetCount", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestCount_Success(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetCount", AnythingOfType("*models.CountParams")).
		Return(10, nil)

	gu := NewGiteaUsecase(mockRepository)

	expected := coreModels.NewTile(api.GiteaCountTileType).WithMetrics(coreModels.NumberUnit)
	expected.Label = "Gitea count"
	expected.Status = coreModels.SuccessStatus
	expected.Metrics.Values = []string{"10"}

	tile, err := gu.Count(&models.CountParams{Owner: "test", Repository: "test"})
	if assert.NoError(t, err) {
		assert.Equal(t, expected, tile)
		mockRepository.AssertNumberOfCalls(t, "GetCount", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestChecks_Error(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetChecks", AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string")).
		Return(nil, errors.New("boom"))

	gu := NewGiteaUsecase(mockRepository)

	tile, err := gu.Checks(&models.ChecksParams{Owner: "test", Repository: "test", Ref: "master"})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &coreModels.MonitororError{}, err)
		assert.Equal(t, "unable to load ref checks", err.Error())
		mockRepository.AssertExpectations(t)
	}
}

func TestChecks_NoChecks(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetChecks", AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string")).
		Return(&models.Checks{}, nil)

	gu := NewGiteaUsecase(mockRepository)

	tile, err := gu.Checks(&models.ChecksParams{Owner: "test", Repository: "test", Ref: "master"})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &coreModels.MonitororError{}, err)
		assert.Equal(t, "no ref checks found", err.Error())
		assert.Equal(t, coreModels.UnknownStatus, err.(*coreModels.MonitororError).ErrorStatus)
		mockRepository.AssertExpectations(t)
	}
}

func TestChecks_Success(t *testing.T) {
	refTime := time.Now().Add(-time.Hour)
	checks := &models.Checks{
		HeadCommit: ToString("sha"),
		Statuses: []models.Status{
			{ID: 1, Title: "build", State: "success", CreatedAt: refTime, UpdatedAt: refTime.Add(time.Minute)},
			{ID: 2, Title: "test", State: "success", CreatedAt: refTime.Add(time.Second), UpdatedAt: refTime.Add(time.Minute * 2)},
		},
	}

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetChecks", "test", "test", "master").Return(checks, nil)

	gu := NewGiteaUsecase(mockRepository)

	expected := coreModels.NewTile(api.GiteaChecksTileType).WithBuild()
	expected.Label = "test"
	expected.Status = coreModels.SuccessStatus
	expected.Build.Branch = ToString("master")
	expected.Build.PreviousStatus = coreModels.UnknownStatus
	expected.Build.StartedAt = ToTime(refTime)
	expected.Build.FinishedAt = ToTime(refTime.Add(time.Minute * 2))

	tile, err := gu.Checks(&models.ChecksParams{Owner: "test", Repository: "test", Ref: "master"})
	if assert.NoError(t, err) {
		assert.Equal(t, expected, tile)
		mockRepository.AssertNotCalled(t, "GetCommit", Anything, Anything, Anything)
		mockRepository.AssertExpectations(t)
	}

	// Second call: previous status comes from build cache
	checks.Statuses = append(checks.Statuses, models.Status{ID: 3, Title: "lint", State: "success", CreatedAt: refTime, UpdatedAt: refTime})
	tile, err = gu.Checks(&models.ChecksParams{Owner: "test", Repository: "test", Ref: "master"})
	if assert.NoError(t, err) {
		assert.Equal(t, coreModels.SuccessStatus, tile.Build.PreviousStatus)
	}
}

func TestChecks_Failure(t *testing.T) {
	refTime := time.Now().Add(-time.Hour)
	checks := &models.Checks{
		HeadCommit: ToString("sha"),
		Statuses: []models.Status{
			{ID: 1, Title: "build", State: "success", CreatedAt: refTime, UpdatedAt: refTime.Add(time.Minute)},
			{ID: 2, Title: "test", State: "failure", CreatedAt: refTime, UpdatedAt: refTime.Add(time.Minute)},
		},
	}

	commit := &models.Commit{SHA: "sha"}
	commit.Author.Name = "John Doe"
	commit.Author.AvatarURL = "https://avatar"

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetChecks", "test", "test", "master").Return(checks, nil)
	mockRepository.On("GetCommit", "test", "test", "sha").Return(commit, nil)

	gu := NewGiteaUsecase(mockRepository)

	tile, err := gu.Checks(&models.ChecksParams{Owner: "test", Repository: "test", Ref: "master"})
	if assert.NoError(t, err) {
		assert.Equal(t, coreModels.FailedStatus, tile.Status)
		assert.Equal(t, &coreModels.Author{Name: "John Doe", AvatarURL: "https://avatar"}, tile.Build.Author)
		mockRepository.AssertExpectations(t)
	}
}

func TestChecks_FailureWithCommitError(t *testing.T) {
	refTime := time.Now().Add(-time.Hour)
	checks := &models.Checks{
		HeadCommit: ToString("sha"),
		Statuses:   []models.Status{{ID: 1, Title: "build", State: "error", CreatedAt: refTime, UpdatedAt: refTime}},
	}

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetChecks", "test", "test", "master").Return(checks, nil)
	mockRepository.On("GetCommit", "test", "test", "sha").Return(nil, errors.New("boom"))

	gu := NewGiteaUsecase(mockRepository)

	tile, err := gu.Checks(&models.ChecksParams{Owner: "test", Repository: "test", Ref: "master"})
	if assert.NoError(t, err) {
		assert.Equal(t, coreModels.FailedStatus, tile.Status)
		assert.Nil(t, tile.Build.Author)
	}
}

func TestChecks_Running(t *testing.T) {
	refTime := time.Now().Add(-time.Minute)
	checks := &models.Checks{
		Statuses: []models.Status{
			{ID: 1, Title: "build", State: "success", CreatedAt: refTime, UpdatedAt: refTime},
			{ID: 2, Title: "test", State: "pending", CreatedAt: refTime, UpdatedAt: refTime},
		},
	}

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetChecks", "test", "test", "master").Return(checks, nil)

	gu := NewGiteaUsecase(mockRepository)

	tile, err := gu.Checks(&models.ChecksParams{Owner: "test", Repository: "test", Ref: "master"})
	if assert.NoError(t, err) {
		assert.Equal(t, coreModels.RunningStatus, tile.Status)
		assert.Nil(t, tile.Build.FinishedAt)
		assert.NotNil(t, tile.Build.Duration)
		assert.Equal(t, ToInt64(0), tile.Build.EstimatedDuration)
	}
}

func TestPullRequest_ErrorOnPullRequest(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetPullRequest", "test", "test", 10).Return(nil, errors.New("boom"))

	gu := NewGiteaUsecase(mockRepository)

	tile, err := gu.PullRequest(&models.PullRequestParams{Owner: "test", Repository: "test", ID: ToInt(10)})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.Equal(t, "unable to load pull request", err.Error())
	}
}

func TestPullRequest_ErrorOnChecks(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetPullRequest", "test", "test", 10).Return(&models.PullRequest{ID: 10, CommitSHA: "sha"}, nil)
	mockRepository.On("GetChecks", "test", "test", "sha").Return(nil, errors.New("boom"))

	gu := NewGiteaUsecase(mockRepository)

	tile, err := gu.PullRequest(&models.PullRequestParams{Owner: "test", Repository: "test", ID: ToInt(10)})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.Equal(t, "unable to load ref checks", err.Error())
	}
}

func TestPullRequest_NoChecks(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetPullRequest", "test", "test", 10).
		Return(&models.PullRequest{ID: 10, Title: "PR", SourceOwner: "test", SourceBranch: "feature", CommitSHA: "sha"}, nil)
	mockRepository.On("GetChecks", "test", "test", "sha").Return(&models.Checks{}, nil)

	gu := NewGiteaUsecase(mockRepository)

	expected := coreModels.NewTile(api.GiteaPullRequestTileType).WithBuild()
	expected.Label = "test"
	expected.Status = coreModels.SuccessStatus
	expected.Build.Branch = ToString("feature")
	expected.Build.PreviousStatus = coreModels.UnknownStatus
	expected.Build.MergeRequest = &coreModels.TileMergeRequest{ID: 10, Title: "PR"}

	tile, err := gu.PullRequest(&models.PullRequestParams{Owner: "test", Repository: "test", ID: ToInt(10)})
	if assert.NoError(t, err) {
		assert.Equal(t, expected, tile)
	}
}

func TestPullRequest_FailureFromFork(t *testing.T) {
	refTime := time.Now().Add(-time.Hour)
	pullRequest := &models.PullRequest{ID: 10, Title: "PR", SourceOwner: "fork", SourceBranch: "feature", CommitSHA: "sha"}
	pullRequest.Author.Name = "John Doe"

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetPullRequest", "test", "test", 10).Return(pullRequest, nil)
	mockRepository.On("GetChecks", "test", "test", "sha").Return(&models.Checks{
		Statuses: []models.Status{{ID: 1, Title: "build", State: "failure", CreatedAt: refTime, UpdatedAt: refTime}},
	}, nil)

	gu := NewGiteaUsecase(mockRepository)

	tile, err := gu.PullRequest(&models.PullRequestParams{Owner: "test", Repository: "test", ID: ToInt(10)})
	if assert.NoError(t, err) {
		assert.Equal(t, coreModels.FailedStatus, tile.Status)
		assert.Equal(t, ToString("fork:feature"), tile.Build.Branch)
		assert.Equal(t, "John Doe", tile.Build.Author.Name)
	}
}

func TestPullRequest_DeletedFork(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetPullRequest", "test", "test", 10).
		Return(&models.PullRequest{ID: 10, SourceOwner: "", SourceBranch: "feature", CommitSHA: "sha"}, nil)
	mockRepository.On("GetChecks", "test", "test", "sha").Return(&models.Checks{}, nil)

	gu := NewGiteaUsecase(mockRepository)

	tile, err := gu.PullRequest(&models.PullRequestParams{Owner: "test", Repository: "test", ID: ToInt(10)})
	if assert.NoError(t, err) {
		assert.Equal(t, ToString("feature"), tile.Build.Branch)
	}
}

func TestPullRequestsGenerator_Error(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetPullRequests", "test", "test").Return(nil, errors.New("boom"))

	gu := NewGiteaUsecase(mockRepository)

	tiles, err := gu.PullRequestsGenerator(&models.PullRequestGeneratorParams{Owner: "test", Repository: "test"})
	if assert.Error(t, err) {
		assert.Nil(t, tiles)
		assert.Equal(t, "unable to load pull request", err.Error())
	}
}

func TestPullRequestsGenerator_Success(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetPullRequests", "test", "test").
		Return([]models.PullRequest{{ID: 1}, {ID: 2}}, nil)

	gu := NewGiteaUsecase(mockRepository)

	tiles, err := gu.PullRequestsGenerator(&models.PullRequestGeneratorParams{Owner: "test", Repository: "test"})
	if assert.NoError(t, err) {
		assert.Equal(t, []uiConfigModels.GeneratedTile{
			{Params: &models.PullRequestParams{Owner: "test", Repository: "test", ID: ToInt(1)}},
			{Params: &models.PullRequestParams{Owner: "test", Repository: "test", ID: ToInt(2)}},
		}, tiles)
	}
}

func TestParseStatus(t *testing.T) {
	for state, expected := range map[string]coreModels.TileStatus{
		"success": coreModels.SuccessStatus,
		"failure": coreModels.FailedStatus,
		"error":   coreModels.FailedStatus,
		"warning": coreModels.WarningStatus,
		"pending": coreModels.RunningStatus,
		"skipped": coreModels.DisabledStatus,
		"unknown": coreModels.UnknownStatus,
	} {
		assert.Equal(t, expected, parseStatus(&models.Status{State: state}), state)
	}
}

func TestConvertChecks_SameSecondHighestIDWins(t *testing.T) {
	// pending then success posted within the same second: the newest (highest ID) must win
	refTime := time.Now().Truncate(time.Second)
	checks := &models.Checks{
		Statuses: []models.Status{
			{ID: 7, Title: "build", State: "pending", CreatedAt: refTime, UpdatedAt: refTime},
			{ID: 8, Title: "build", State: "success", CreatedAt: refTime, UpdatedAt: refTime},
		},
	}

	statuses, _, _, _ := convertChecks(checks)
	assert.Equal(t, []coreModels.TileStatus{coreModels.SuccessStatus}, statuses)
}

func TestChecks_DurationFromStatusHistory(t *testing.T) {
	// Gitea returns one row per state change: the pending row gives the start of the build
	refTime := time.Now().Add(-time.Hour)
	checks := &models.Checks{
		Statuses: []models.Status{
			{ID: 2, Title: "build", State: "success", CreatedAt: refTime.Add(time.Minute * 3), UpdatedAt: refTime.Add(time.Minute * 3)},
			{ID: 1, Title: "build", State: "pending", CreatedAt: refTime, UpdatedAt: refTime},
		},
	}

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetChecks", "test", "test", "master").Return(checks, nil)

	gu := NewGiteaUsecase(mockRepository)

	tile, err := gu.Checks(&models.ChecksParams{Owner: "test", Repository: "test", Ref: "master"})
	if assert.NoError(t, err) {
		assert.Equal(t, coreModels.SuccessStatus, tile.Status)
		assert.Equal(t, time.Minute*3, tile.Build.FinishedAt.Sub(*tile.Build.StartedAt))
	}
}

func TestConvertChecks_DeduplicateAndOrder(t *testing.T) {
	refTime := time.Now()
	checks := &models.Checks{
		Statuses: []models.Status{
			// Older "build" status is ignored, the newest one wins
			{ID: 1, Title: "build", State: "failure", CreatedAt: refTime.Add(-time.Hour), UpdatedAt: refTime.Add(-time.Hour)},
			{ID: 2, Title: "build", State: "success", CreatedAt: refTime, UpdatedAt: refTime.Add(time.Minute)},
			{ID: 3, Title: "lint", State: "warning", CreatedAt: refTime.Add(-time.Minute), UpdatedAt: refTime},
		},
	}

	statuses, startedAt, finishedAt, id := convertChecks(checks)
	assert.Equal(t, []coreModels.TileStatus{coreModels.WarningStatus, coreModels.SuccessStatus}, statuses)
	assert.Equal(t, refTime.Add(-time.Hour), *startedAt)
	assert.Equal(t, refTime.Add(time.Minute), *finishedAt)
	assert.Equal(t, hash.GetMD5Hash("-2-3-1"), id)
}
