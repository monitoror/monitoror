package repository

import (
	"errors"
	"testing"

	. "github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/monitorable/travisci/models"
	pkgTravis "github.com/monitoror/monitoror/pkg/gotravis"
	"github.com/monitoror/monitoror/pkg/gotravis/mocks"

	. "github.com/AlekSi/pointer"
	"github.com/shuheiktgw/go-travis"
	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func initRepository(t *testing.T, buildsAPI pkgTravis.TravisCI) *travisCIRepository {
	conf := InitConfig()
	repository := NewTravisCIRepository(conf.Monitorable.TravisCI[DefaultVariant], conf.Monitorable.Github[DefaultVariant])

	apiTravisCIRepository, ok := repository.(*travisCIRepository)
	if assert.True(t, ok) {
		apiTravisCIRepository.travisBuildsAPI = buildsAPI
		return apiTravisCIRepository
	}
	return nil
}

func TestNewApiTravisCIRepository_Panic(t *testing.T) {
	conf := InitConfig()
	conf.Monitorable.Github[DefaultVariant].Token = "token"
	conf.Monitorable.TravisCI[DefaultVariant].URL = ""

	// Panic because ApiURL is not define
	assert.Panics(t, func() {
		_ = NewTravisCIRepository(conf.Monitorable.TravisCI[DefaultVariant], conf.Monitorable.Github[DefaultVariant])
	})
}

func TestRepository_GetLastBuildStatus_Error(t *testing.T) {
	// Params
	travisErr := errors.New("TravisCI Error")

	mockTravis := new(mocks.TravisCI)
	mockTravis.On("ListByRepoSlug", Anything, AnythingOfType("string"), Anything).
		Return([]*travis.Build{}, nil, travisErr)

	repository := initRepository(t, mockTravis)
	if repository != nil {
		_, err := repository.GetLastBuildStatus("test", "test", "test")
		assert.Error(t, err)
		assert.Equal(t, travisErr, err)
		mockTravis.AssertNumberOfCalls(t, "ListByRepoSlug", 1)
		mockTravis.AssertExpectations(t)
	}
}

func TestRepository_GetLastBuildStatus_NoBuild(t *testing.T) {
	mockTravis := new(mocks.TravisCI)
	mockTravis.On("ListByRepoSlug", Anything, AnythingOfType("string"), Anything).
		Return([]*travis.Build{}, nil, nil)

	repository := initRepository(t, mockTravis)
	if repository != nil {
		build, err := repository.GetLastBuildStatus("test", "test", "test")
		assert.NoError(t, err)
		assert.Nil(t, build)
		mockTravis.AssertNumberOfCalls(t, "ListByRepoSlug", 1)
		mockTravis.AssertExpectations(t)
	}
}

func TestRepository_GetLastBuildStatus_Success(t *testing.T) {
	// Params
	travisBuild := &travis.Build{
		Id: ToUint(1),
		Branch: &travis.Branch{
			Name: ToString("test"),
		},
		Commit: &travis.Commit{
			Author: &travis.Author{
				Name:      "test",
				AvatarURL: "monitoror.example.com",
			},
		},
		State:         ToString("passed"),
		PreviousState: ToString("passed"),
		StartedAt:     ToString("2019-04-12T20:39:59Z"),
		FinishedAt:    ToString("2019-04-12T20:39:59Z"),
		Duration:      ToUint(154),
	}

	mockTravis := new(mocks.TravisCI)
	mockTravis.On("ListByRepoSlug", Anything, AnythingOfType("string"), Anything).
		Return([]*travis.Build{travisBuild}, nil, nil)

	// Expected
	expectedBuild := &models.Build{
		ID:     1,
		Branch: *travisBuild.Branch.Name,
		Author: models.Author{
			Name:      travisBuild.Commit.Author.Name,
			AvatarURL: travisBuild.Commit.Author.AvatarURL,
		},
		State:      *travisBuild.State,
		StartedAt:  parseDate(*travisBuild.StartedAt),
		FinishedAt: parseDate(*travisBuild.FinishedAt),
		Duration:   parseDuration(*travisBuild.Duration),
	}

	repository := initRepository(t, mockTravis)
	if repository != nil {
		build, err := repository.GetLastBuildStatus("test", "test", "test")
		assert.NoError(t, err)
		assert.Equal(t, expectedBuild, build)
		mockTravis.AssertNumberOfCalls(t, "ListByRepoSlug", 1)
		mockTravis.AssertExpectations(t)
	}
}
