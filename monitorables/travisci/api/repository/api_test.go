package repository

import (
	"errors"
	"testing"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/travisci/api/models"
	"github.com/monitoror/monitoror/monitorables/travisci/config"
	pkgTravis "github.com/monitoror/monitoror/pkg/gotravis"
	"github.com/monitoror/monitoror/pkg/gotravis/mocks"

	. "github.com/AlekSi/pointer"
	"github.com/shuheiktgw/go-travis"
	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func initRepository(t *testing.T, buildsAPI pkgTravis.TravisCI) *travisCIRepository {
	conf := &config.TravisCI{
		URL:             config.Default.URL,
		Token:           config.Default.Token,
		GithubToken:     config.Default.GithubToken,
		Timeout:         config.Default.Timeout,
		InitialMaxDelay: config.Default.InitialMaxDelay,
	}

	repository := NewTravisCIRepository(conf)

	apiTravisCIRepository, ok := repository.(*travisCIRepository)
	if assert.True(t, ok) {
		apiTravisCIRepository.travisBuildsAPI = buildsAPI
		return apiTravisCIRepository
	}
	return nil
}

func TestNewApiTravisCIRepository_Panic(t *testing.T) {
	conf := &config.TravisCI{
		URL:             "",
		Token:           config.Default.Token,
		GithubToken:     "token",
		Timeout:         config.Default.Timeout,
		InitialMaxDelay: config.Default.InitialMaxDelay,
	}
	// Panic because ApiURL is not define
	assert.Panics(t, func() {
		_ = NewTravisCIRepository(conf)
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
		Author: coreModels.Author{
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
