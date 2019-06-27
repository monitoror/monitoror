package usecase

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	. "github.com/AlekSi/pointer"

	"github.com/monitoror/monitoror/config"

	"github.com/monitoror/monitoror/monitorable/travisci"

	mErrors "github.com/monitoror/monitoror/models/errors"
	. "github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/monitorable/travisci/mocks"
	"github.com/monitoror/monitoror/monitorable/travisci/models"
	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

var group, repo, branch = "test", "test", "master"

func TestBuild_Error_NoHost(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetBuildStatus", Anything, AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string")).
		Return(nil, errors.New("no such host"))

	conf := config.InitConfig()
	tu := NewTravisCIUsecase(conf, mockRepository)

	tile, err := tu.Build(&models.BuildParams{Group: group, Repository: repo, Branch: branch})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &mErrors.TimeoutError{}, err)
		mockRepository.AssertNumberOfCalls(t, "GetBuildStatus", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestBuild_Error_NoNetwork(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetBuildStatus", Anything, AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string")).
		Return(nil, errors.New("dial tcp: lookup"))

	conf := config.InitConfig()
	tu := NewTravisCIUsecase(conf, mockRepository)

	tile, err := tu.Build(&models.BuildParams{Group: group, Repository: repo, Branch: branch})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &mErrors.TimeoutError{}, err)
		mockRepository.AssertNumberOfCalls(t, "GetBuildStatus", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestBuild_Timeout(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetBuildStatus", Anything, AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string")).
		Return(nil, context.DeadlineExceeded)

	conf := config.InitConfig()
	tu := NewTravisCIUsecase(conf, mockRepository)

	tile, err := tu.Build(&models.BuildParams{Group: group, Repository: repo, Branch: branch})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &mErrors.TimeoutError{}, err)
		mockRepository.AssertNumberOfCalls(t, "GetBuildStatus", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestBuild_Error_System(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetBuildStatus", Anything, AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string")).
		Return(nil, errors.New("boom"))

	conf := config.InitConfig()
	tu := NewTravisCIUsecase(conf, mockRepository)

	tile, err := tu.Build(&models.BuildParams{Group: group, Repository: repo, Branch: branch})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &mErrors.SystemError{}, err)
		mockRepository.AssertNumberOfCalls(t, "GetBuildStatus", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestBuild_Error_NoBuild(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetBuildStatus", Anything, AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string")).
		Return(nil, nil)

	conf := config.InitConfig()
	tu := NewTravisCIUsecase(conf, mockRepository)

	tile, err := tu.Build(&models.BuildParams{Group: group, Repository: repo, Branch: branch})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &mErrors.NoBuildError{}, err)
		mockRepository.AssertNumberOfCalls(t, "GetBuildStatus", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestBuild_Success(t *testing.T) {
	build := buildResponse(branch, "passed", "", time.Now(), time.Now(), 100)

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetBuildStatus", Anything, AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string")).
		Return(build, nil)

	conf := config.InitConfig()
	tu := NewTravisCIUsecase(conf, mockRepository)
	tUsecase, ok := tu.(*travisCIUsecase)
	if assert.True(t, ok, "enable to case tu into travisCIUsecase") {
		// Expected
		expected := NewBuildTile(travisci.TravisCIBuildTileType)
		expected.Label = fmt.Sprintf("%s : #%s", repo, branch)
		expected.Status = parseState(build.State)
		expected.StartedAt = ToInt64(build.StartedAt.Unix())
		expected.FinishedAt = ToInt64(build.FinishedAt.Unix())
		expected.Author = &Author{
			Name:      build.Author.Name,
			AvatarUrl: build.Author.AvatarUrl,
		}

		// Tests
		tile, err := tu.Build(&models.BuildParams{Group: group, Repository: repo, Branch: branch})
		if assert.NotNil(t, tile) {
			assert.NoError(t, err)
			assert.Equal(t, expected, tile)

			// Check if duration is added into cache
			previousDuration, ok := tUsecase.estimatedDurations[tile.Label]
			assert.True(t, ok)
			assert.Equal(t, build.Duration, previousDuration)

			mockRepository.AssertNumberOfCalls(t, "GetBuildStatus", 1)
			mockRepository.AssertExpectations(t)
		}
	}
}

func TestBuild_Failed(t *testing.T) {
	build := buildResponse(branch, "failed", "", time.Now(), time.Now(), 100)

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetBuildStatus", Anything, AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string")).
		Return(build, nil)

	conf := config.InitConfig()
	tu := NewTravisCIUsecase(conf, mockRepository)
	tUsecase, ok := tu.(*travisCIUsecase)
	if assert.True(t, ok, "enable to case tu into travisCIUsecase") {
		// Expected
		expected := NewBuildTile(travisci.TravisCIBuildTileType)
		expected.Label = fmt.Sprintf("%s : #%s", repo, branch)
		expected.Status = parseState(build.State)
		expected.StartedAt = ToInt64(build.StartedAt.Unix())
		expected.FinishedAt = ToInt64(build.FinishedAt.Unix())
		expected.Author = &Author{
			Name:      build.Author.Name,
			AvatarUrl: build.Author.AvatarUrl,
		}

		tile, err := tu.Build(&models.BuildParams{Group: group, Repository: repo, Branch: branch})
		if assert.NotNil(t, tile) {
			assert.NoError(t, err)
			assert.Equal(t, expected, tile)

			// Check if duration is not added into cache
			_, ok = tUsecase.estimatedDurations[tile.Label]
			assert.False(t, ok)

			mockRepository.AssertNumberOfCalls(t, "GetBuildStatus", 1)
			mockRepository.AssertExpectations(t)
		}
	}
}

func TestBuild_Queued(t *testing.T) {
	build := buildResponse(branch, "received", "passed", time.Now(), time.Time{}, 100)

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetBuildStatus", Anything, AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string")).
		Return(build, nil)

	conf := config.InitConfig()
	tu := NewTravisCIUsecase(conf, mockRepository)

	// Expected
	expected := NewBuildTile(travisci.TravisCIBuildTileType)
	expected.Label = fmt.Sprintf("%s : #%s", repo, branch)
	expected.Status = parseState(build.State)
	expected.PreviousStatus = parseState(build.PreviousState)
	expected.StartedAt = ToInt64(build.StartedAt.Unix())
	expected.Author = &Author{
		Name:      build.Author.Name,
		AvatarUrl: build.Author.AvatarUrl,
	}

	// Without Estimated Duration
	tile, err := tu.Build(&models.BuildParams{Group: group, Repository: repo, Branch: branch})
	if assert.NotNil(t, tile) {
		assert.NoError(t, err)
		assert.Equal(t, expected, tile)
	}
}

func TestBuild_Running(t *testing.T) {
	build := buildResponse(branch, "started", "passed", time.Now(), time.Time{}, 100)

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetBuildStatus", Anything, AnythingOfType("string"), AnythingOfType("string"), AnythingOfType("string")).
		Return(build, nil)

	conf := config.InitConfig()
	tu := NewTravisCIUsecase(conf, mockRepository)
	tUsecase, ok := tu.(*travisCIUsecase)
	if assert.True(t, ok, "enable to case tu into travisCIUsecase") {

		// Expected
		expected := NewBuildTile(travisci.TravisCIBuildTileType)
		expected.Label = fmt.Sprintf("%s : #%s", repo, branch)
		expected.Status = parseState(build.State)
		expected.PreviousStatus = parseState(build.PreviousState)
		expected.Duration = ToInt64(int64(build.Duration / time.Second))
		expected.StartedAt = ToInt64(build.StartedAt.Unix())
		expected.Author = &Author{
			Name:      build.Author.Name,
			AvatarUrl: build.Author.AvatarUrl,
		}

		// Without Estimated Duration
		tile, err := tu.Build(&models.BuildParams{Group: group, Repository: repo, Branch: branch})
		if assert.NotNil(t, tile) {
			assert.NoError(t, err)
			assert.Equal(t, expected, tile)
		}

		// With Estimated Duration
		expected.EstimatedDuration = ToInt64(int64(120))
		tUsecase.estimatedDurations[expected.Label] = time.Second * 120
		tile, err = tu.Build(&models.BuildParams{Group: group, Repository: repo, Branch: branch})
		if assert.NotNil(t, tile) {
			assert.NoError(t, err)
			assert.Equal(t, expected, tile)
		}

		mockRepository.AssertNumberOfCalls(t, "GetBuildStatus", 2)
		mockRepository.AssertExpectations(t)
	}
}

func TestParseState(t *testing.T) {
	assert.Equal(t, QueuedStatus, parseState("created"))
	assert.Equal(t, QueuedStatus, parseState("received"))
	assert.Equal(t, RunningStatus, parseState("started"))
	assert.Equal(t, SuccessStatus, parseState("passed"))
	assert.Equal(t, FailedStatus, parseState("failed"))
	assert.Equal(t, FailedStatus, parseState("errored"))
	assert.Equal(t, UnknownStatus, parseState(""))
}

func buildResponse(branch, state, previousState string, startedAt, finishedAt time.Time, duration time.Duration) *models.Build {
	return &models.Build{
		Branch: branch,
		Author: models.Author{
			Name:      "me",
			AvatarUrl: "http://avatar.com",
		},
		State:         state,
		PreviousState: previousState,
		StartedAt:     startedAt,
		FinishedAt:    finishedAt,
		Duration:      duration,
	}
}
