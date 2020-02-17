package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/travisci"
	"github.com/monitoror/monitoror/monitorable/travisci/mocks"
	travisModels "github.com/monitoror/monitoror/monitorable/travisci/models"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/git"

	. "github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

var owner, repo, branch = "test", "test", "master"

func TestBuild_Error(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetLastBuildStatus", AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string")).
		Return(nil, errors.New("boom"))

	tu := NewTravisCIUsecase(mockRepository)

	tile, err := tu.Build(&travisModels.BuildParams{Owner: owner, Repository: repo, Branch: branch})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &models.MonitororError{}, err)
		assert.Equal(t, "unable to find build", err.Error())
		mockRepository.AssertNumberOfCalls(t, "GetLastBuildStatus", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestBuild_Error_NoBuild(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetLastBuildStatus", AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string")).
		Return(nil, nil)

	tu := NewTravisCIUsecase(mockRepository)

	tile, err := tu.Build(&travisModels.BuildParams{Owner: owner, Repository: repo, Branch: branch})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &models.MonitororError{}, err)
		assert.Equal(t, "no build found", err.Error())
		assert.Equal(t, models.UnknownStatus, err.(*models.MonitororError).ErrorStatus)
		mockRepository.AssertNumberOfCalls(t, "GetLastBuildStatus", 1)
		mockRepository.AssertExpectations(t)
	}
}

//nolint:dupl
func TestBuild_Success(t *testing.T) {
	build := buildResponse(branch, "passed", time.Now(), time.Now(), time.Second*100)

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetLastBuildStatus", AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string")).
		Return(build, nil)

	tu := NewTravisCIUsecase(mockRepository)
	tUsecase, ok := tu.(*travisCIUsecase)
	if assert.True(t, ok, "enable to case tu into travisCIUsecase") {
		// Expected
		expected := models.NewTile(travisci.TravisCIBuildTileType).WithBuild()
		expected.Label = repo
		expected.Build.Branch = ToString(git.HumanizeBranch(branch))
		expected.Build.ID = ToString("1")

		expected.Status = parseState(build.State)
		expected.Build.PreviousStatus = models.SuccessStatus
		expected.Build.StartedAt = ToTime(build.StartedAt)
		expected.Build.FinishedAt = ToTime(build.FinishedAt)

		// Tests
		params := &travisModels.BuildParams{Owner: owner, Repository: repo, Branch: branch}
		tUsecase.buildsCache.Add(params, "0", models.SuccessStatus, time.Second*120)
		tile, err := tu.Build(params)
		if assert.NotNil(t, tile) {
			assert.NoError(t, err)
			assert.Equal(t, expected, tile)

			// Check if duration is added into cache
			previousDuration := tUsecase.buildsCache.GetEstimatedDuration(params)
			assert.NotNil(t, previousDuration)
			assert.Equal(t, time.Second*110, *previousDuration)

			mockRepository.AssertNumberOfCalls(t, "GetLastBuildStatus", 1)
			mockRepository.AssertExpectations(t)
		}
	}
}

//nolint:dupl
func TestBuild_Failed(t *testing.T) {
	build := buildResponse(branch, "failed", time.Now(), time.Now(), time.Second*100)

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetLastBuildStatus", AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string")).
		Return(build, nil)

	tu := NewTravisCIUsecase(mockRepository)
	tUsecase, ok := tu.(*travisCIUsecase)
	if assert.True(t, ok, "enable to case tu into travisCIUsecase") {
		// Expected
		expected := models.NewTile(travisci.TravisCIBuildTileType).WithBuild()
		expected.Label = repo
		expected.Build.Branch = ToString(git.HumanizeBranch(branch))
		expected.Build.ID = ToString("1")

		expected.Status = parseState(build.State)
		expected.Build.PreviousStatus = models.SuccessStatus
		expected.Build.StartedAt = ToTime(build.StartedAt)
		expected.Build.FinishedAt = ToTime(build.FinishedAt)
		expected.Build.Author = &models.Author{
			Name:      build.Author.Name,
			AvatarURL: build.Author.AvatarURL,
		}

		params := &travisModels.BuildParams{Owner: owner, Repository: repo, Branch: branch}
		tUsecase.buildsCache.Add(params, "0", models.SuccessStatus, time.Second*120)
		tile, err := tu.Build(params)
		if assert.NotNil(t, tile) {
			assert.NoError(t, err)
			assert.Equal(t, expected, tile)

			// Check if duration is added into cache
			previousDuration := tUsecase.buildsCache.GetEstimatedDuration(params)
			assert.NotNil(t, previousDuration)
			assert.Equal(t, time.Second*110, *previousDuration)

			mockRepository.AssertNumberOfCalls(t, "GetLastBuildStatus", 1)
			mockRepository.AssertExpectations(t)
		}
	}
}

//nolint:dupl
func TestBuild_Queued(t *testing.T) {
	build := buildResponse(branch, "received", time.Now(), time.Time{}, 100)

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetLastBuildStatus", AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string")).
		Return(build, nil)

	tu := NewTravisCIUsecase(mockRepository)
	tUsecase, ok := tu.(*travisCIUsecase)
	if assert.True(t, ok) {
		// Expected
		expected := models.NewTile(travisci.TravisCIBuildTileType).WithBuild()
		expected.Label = repo
		expected.Build.Branch = ToString(git.HumanizeBranch(branch))
		expected.Build.ID = ToString("1")

		expected.Status = parseState(build.State)
		expected.Build.PreviousStatus = models.SuccessStatus
		expected.Build.StartedAt = ToTime(build.StartedAt)

		// Without Estimated Duration
		params := &travisModels.BuildParams{Owner: owner, Repository: repo, Branch: branch}
		tUsecase.buildsCache.Add(params, "0", models.SuccessStatus, time.Second*10)
		tile, err := tu.Build(params)
		if assert.NotNil(t, tile) {
			assert.NoError(t, err)
			assert.Equal(t, expected, tile)
		}
	}
}

func TestBuild_Running(t *testing.T) {
	build := buildResponse(branch, "started", time.Now(), time.Time{}, 100)

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetLastBuildStatus", AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string")).
		Return(build, nil)

	tu := NewTravisCIUsecase(mockRepository)
	tUsecase, ok := tu.(*travisCIUsecase)
	if assert.True(t, ok, "enable to case tu into travisCIUsecase") {
		// Expected
		expected := models.NewTile(travisci.TravisCIBuildTileType).WithBuild()
		expected.Label = repo
		expected.Build.Branch = ToString(git.HumanizeBranch(branch))
		expected.Build.ID = ToString("1")

		expected.Status = parseState(build.State)
		expected.Build.PreviousStatus = models.UnknownStatus
		expected.Build.Duration = ToInt64(int64(build.Duration / time.Second))
		expected.Build.EstimatedDuration = ToInt64(int64(0))
		expected.Build.StartedAt = ToTime(build.StartedAt)

		// Without Previous Build
		params := &travisModels.BuildParams{Owner: owner, Repository: repo, Branch: branch}
		tile, err := tu.Build(params)
		if assert.NotNil(t, tile) {
			assert.NoError(t, err)
			assert.Equal(t, expected, tile)
		}

		// With Previous Build
		expected.Build.PreviousStatus = models.SuccessStatus
		expected.Build.EstimatedDuration = ToInt64(int64(120))
		tUsecase.buildsCache.Add(params, "0", models.SuccessStatus, time.Second*120)
		tile, err = tu.Build(params)
		if assert.NotNil(t, tile) {
			assert.NoError(t, err)
			assert.Equal(t, expected, tile)
		}

		mockRepository.AssertNumberOfCalls(t, "GetLastBuildStatus", 2)
		mockRepository.AssertExpectations(t)
	}
}

//nolint:dupl
func TestBuild_Aborded(t *testing.T) {
	build := buildResponse(branch, "canceled", time.Now(), time.Time{}, 100)

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetLastBuildStatus", AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string")).
		Return(build, nil)

	tu := NewTravisCIUsecase(mockRepository)
	tUsecase, ok := tu.(*travisCIUsecase)
	if assert.True(t, ok) {
		// Expected
		expected := models.NewTile(travisci.TravisCIBuildTileType).WithBuild()
		expected.Label = repo
		expected.Build.Branch = ToString(git.HumanizeBranch(branch))
		expected.Build.ID = ToString("1")

		expected.Status = parseState(build.State)
		expected.Build.PreviousStatus = models.SuccessStatus
		expected.Build.StartedAt = ToTime(build.StartedAt)

		// Without Estimated Duration
		params := &travisModels.BuildParams{Owner: owner, Repository: repo, Branch: branch}
		tUsecase.buildsCache.Add(params, "0", models.SuccessStatus, time.Second*10)
		tile, err := tu.Build(params)
		if assert.NotNil(t, tile) {
			assert.NoError(t, err)
			assert.Equal(t, expected, tile)
		}
	}
}

func TestParseState(t *testing.T) {
	assert.Equal(t, models.QueuedStatus, parseState("created"))
	assert.Equal(t, models.QueuedStatus, parseState("received"))
	assert.Equal(t, models.RunningStatus, parseState("started"))
	assert.Equal(t, models.SuccessStatus, parseState("passed"))
	assert.Equal(t, models.FailedStatus, parseState("failed"))
	assert.Equal(t, models.FailedStatus, parseState("errored"))
	assert.Equal(t, models.CanceledStatus, parseState("canceled"))
	assert.Equal(t, models.UnknownStatus, parseState(""))
}

func buildResponse(branch, state string, startedAt, finishedAt time.Time, duration time.Duration) *travisModels.Build {
	return &travisModels.Build{
		ID:     1,
		Branch: branch,
		Author: models.Author{
			Name:      "me",
			AvatarURL: "http://avatar.com",
		},
		State:      state,
		StartedAt:  startedAt,
		FinishedAt: finishedAt,
		Duration:   duration,
	}
}
