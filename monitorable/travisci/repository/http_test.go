package repository

import (
	"context"
	"errors"
	"testing"

	. "github.com/AlekSi/pointer"

	. "github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/monitorable/travisci/models"
	"github.com/stretchr/testify/assert"

	pkgTravis "github.com/monitoror/monitoror/pkg/gotravis"
	"github.com/shuheiktgw/go-travis"
	. "github.com/stretchr/testify/mock"

	"github.com/monitoror/monitoror/pkg/gotravis/mocks"
)

func initRepository(t *testing.T, buildsApi pkgTravis.Builds) *travisCIRepository {
	conf := InitConfig()
	repository := NewTravisCIRepository(conf)

	apiTravisCIRepository, ok := repository.(*travisCIRepository)
	if assert.True(t, ok) {
		apiTravisCIRepository.travisBuildsApi = buildsApi
		return apiTravisCIRepository
	}
	return nil
}

func TestNewApiTravisCIRepository_Panic(t *testing.T) {
	conf := InitConfig()
	conf.Monitorable.Github.Token = "token"
	conf.Monitorable.TravisCI.Url = ""

	// Panic because ApiUrl is not define
	assert.Panics(t, func() { _ = NewTravisCIRepository(conf) })
}

func TestRepository_Build_Error(t *testing.T) {
	// Params
	travisErr := errors.New("TravisCI Error")

	mockTravis := new(mocks.Builds)
	mockTravis.On("ListByRepoSlug", Anything, AnythingOfType("string"), Anything).
		Return([]*travis.Build{}, nil, travisErr)

	repository := initRepository(t, mockTravis)
	if repository != nil {
		_, err := repository.GetBuildStatus(context.Background(), "test", "test", "test")
		assert.Error(t, err)
		assert.Equal(t, travisErr, err)
		mockTravis.AssertNumberOfCalls(t, "ListByRepoSlug", 1)
		mockTravis.AssertExpectations(t)
	}
}

func TestRepository_Build_NoBuild(t *testing.T) {
	mockTravis := new(mocks.Builds)
	mockTravis.On("ListByRepoSlug", Anything, AnythingOfType("string"), Anything).
		Return([]*travis.Build{}, nil, nil)

	repository := initRepository(t, mockTravis)
	if repository != nil {
		build, err := repository.GetBuildStatus(context.Background(), "test", "test", "test")
		assert.NoError(t, err)
		assert.Nil(t, build)
		mockTravis.AssertNumberOfCalls(t, "ListByRepoSlug", 1)
		mockTravis.AssertExpectations(t)
	}
}

func TestRepository_Build_Success(t *testing.T) {
	// Params
	travisBuild := &travis.Build{
		Branch: &travis.Branch{
			Name: ToString("test"),
		},
		Commit: &travis.Commit{
			Author: &travis.Author{
				Name:      "test",
				AvatarURL: "test.com",
			},
		},
		State:         ToString("passed"),
		PreviousState: ToString("passed"),
		StartedAt:     ToString("2019-04-12T20:39:59Z"),
		FinishedAt:    ToString("2019-04-12T20:39:59Z"),
		Duration:      ToUint(154),
	}

	mockTravis := new(mocks.Builds)
	mockTravis.On("ListByRepoSlug", Anything, AnythingOfType("string"), Anything).
		Return([]*travis.Build{travisBuild}, nil, nil)

	// Expected
	expectedBuild := &models.Build{
		Branch: *travisBuild.Branch.Name,
		Author: models.Author{
			Name:      travisBuild.Commit.Author.Name,
			AvatarUrl: travisBuild.Commit.Author.AvatarURL,
		},
		State:         *travisBuild.State,
		PreviousState: *travisBuild.PreviousState,
		StartedAt:     parseDate(*travisBuild.StartedAt),
		FinishedAt:    parseDate(*travisBuild.FinishedAt),
		Duration:      parseDuration(*travisBuild.Duration),
	}

	repository := initRepository(t, mockTravis)
	if repository != nil {
		build, err := repository.GetBuildStatus(context.Background(), "test", "test", "test")
		assert.NoError(t, err)
		assert.Equal(t, expectedBuild, build)
		mockTravis.AssertNumberOfCalls(t, "ListByRepoSlug", 1)
		mockTravis.AssertExpectations(t)
	}
}
