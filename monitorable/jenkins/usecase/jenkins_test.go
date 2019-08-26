package usecase

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/monitoror/monitoror/config"

	. "github.com/AlekSi/pointer"
	"github.com/monitoror/monitoror/monitorable/jenkins"

	. "github.com/monitoror/monitoror/models/tiles"

	mErrors "github.com/monitoror/monitoror/models/errors"
	"github.com/monitoror/monitoror/monitorable/jenkins/mocks"
	"github.com/monitoror/monitoror/monitorable/jenkins/models"

	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

var job, branch = "test", "master"

func TestBuild_Error_NoHost(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetJob", AnythingOfType("string"), AnythingOfType("string")).
		Return(nil, errors.New("no such host"))

	tu := NewJenkinsUsecase(mockRepository, config.Cache{})

	tile, err := tu.Build(&models.BuildParams{Job: job, Branch: branch})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &mErrors.TimeoutError{}, err)
		mockRepository.AssertNumberOfCalls(t, "GetJob", 1)
		mockRepository.AssertExpectations(t)
	}

	repositoryJob := &models.Job{
		Buildable: true,
	}

	mockRepository = new(mocks.Repository)
	mockRepository.On("GetJob", AnythingOfType("string"), AnythingOfType("string")).
		Return(repositoryJob, nil)
	mockRepository.On("GetLastBuildStatus", Anything).
		Return(nil, errors.New("no such host"))

	tu = NewJenkinsUsecase(mockRepository, config.Cache{})

	tile, err = tu.Build(&models.BuildParams{Job: job, Branch: branch})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &mErrors.TimeoutError{}, err)
		mockRepository.AssertNumberOfCalls(t, "GetJob", 1)
		mockRepository.AssertNumberOfCalls(t, "GetLastBuildStatus", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestBuild_Error_NoNetwork(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetJob", AnythingOfType("string"), AnythingOfType("string")).
		Return(nil, errors.New("dial tcp: lookup"))

	tu := NewJenkinsUsecase(mockRepository, config.Cache{})

	tile, err := tu.Build(&models.BuildParams{Job: job, Branch: branch})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &mErrors.TimeoutError{}, err)
		mockRepository.AssertNumberOfCalls(t, "GetJob", 1)
		mockRepository.AssertExpectations(t)
	}

	repositoryJob := &models.Job{
		Buildable: true,
	}

	mockRepository = new(mocks.Repository)
	mockRepository.On("GetJob", AnythingOfType("string"), AnythingOfType("string")).
		Return(repositoryJob, nil)
	mockRepository.On("GetLastBuildStatus", Anything).
		Return(nil, errors.New("dial tcp: lookup"))

	tu = NewJenkinsUsecase(mockRepository, config.Cache{})

	tile, err = tu.Build(&models.BuildParams{Job: job, Branch: branch})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &mErrors.TimeoutError{}, err)
		mockRepository.AssertNumberOfCalls(t, "GetJob", 1)
		mockRepository.AssertNumberOfCalls(t, "GetLastBuildStatus", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestBuild_Timeout(t *testing.T) {
	var errRequestCanceledConn = errors.New("net/http: request canceled while waiting for connection")

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetJob", AnythingOfType("string"), AnythingOfType("string")).
		Return(nil, errRequestCanceledConn)

	tu := NewJenkinsUsecase(mockRepository, config.Cache{})

	tile, err := tu.Build(&models.BuildParams{Job: job, Branch: branch})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &mErrors.TimeoutError{}, err)
		mockRepository.AssertNumberOfCalls(t, "GetJob", 1)
		mockRepository.AssertExpectations(t)
	}

	repositoryJob := &models.Job{
		Buildable: true,
	}

	mockRepository = new(mocks.Repository)
	mockRepository.On("GetJob", AnythingOfType("string"), AnythingOfType("string")).
		Return(repositoryJob, nil)
	mockRepository.On("GetLastBuildStatus", Anything).
		Return(nil, errRequestCanceledConn)

	tu = NewJenkinsUsecase(mockRepository, config.Cache{})

	tile, err = tu.Build(&models.BuildParams{Job: job, Branch: branch})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &mErrors.TimeoutError{}, err)
		mockRepository.AssertNumberOfCalls(t, "GetJob", 1)
		mockRepository.AssertNumberOfCalls(t, "GetLastBuildStatus", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestBuild_Error_System(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetJob", AnythingOfType("string"), AnythingOfType("string")).
		Return(nil, errors.New("boom"))

	tu := NewJenkinsUsecase(mockRepository, config.Cache{})

	tile, err := tu.Build(&models.BuildParams{Job: job, Branch: branch})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &mErrors.SystemError{}, err)
		mockRepository.AssertNumberOfCalls(t, "GetJob", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestBuild_DisabledBuild(t *testing.T) {
	repositoryJob := &models.Job{
		Buildable: false,
	}

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetJob", AnythingOfType("string"), AnythingOfType("string")).
		Return(repositoryJob, nil)

	tu := NewJenkinsUsecase(mockRepository, config.Cache{})

	tile, err := tu.Build(&models.BuildParams{Job: job})
	if assert.NoError(t, err) {
		assert.Equal(t, job, tile.Label)
		assert.Equal(t, DisabledStatus, tile.Status)
		mockRepository.AssertNumberOfCalls(t, "GetJob", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestBuild_Error_NoBuild(t *testing.T) {
	repositoryJob := &models.Job{
		Buildable: true,
	}

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetJob", AnythingOfType("string"), AnythingOfType("string")).
		Return(repositoryJob, nil)
	mockRepository.On("GetLastBuildStatus", Anything).
		Return(nil, errors.New("boom"))

	tu := NewJenkinsUsecase(mockRepository, config.Cache{})

	tile, err := tu.Build(&models.BuildParams{Job: job, Branch: branch})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &mErrors.NoBuildError{}, err)
		mockRepository.AssertNumberOfCalls(t, "GetJob", 1)
		mockRepository.AssertNumberOfCalls(t, "GetLastBuildStatus", 1)
		mockRepository.AssertExpectations(t)
	}
}

func CheckBuild(t *testing.T, result string) {
	repositoryJob := &models.Job{
		Buildable: true,
	}
	repositoryBuild := buildResponse(result, time.Date(2000, 01, 01, 10, 00, 00, 00, time.UTC), time.Minute)

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetJob", AnythingOfType("string"), AnythingOfType("string")).
		Return(repositoryJob, nil)
	mockRepository.On("GetLastBuildStatus", Anything).
		Return(repositoryBuild, nil)

	tu := NewJenkinsUsecase(mockRepository, config.Cache{})
	tUsecase, ok := tu.(*jenkinsUsecase)
	if assert.True(t, ok, "enable to case tu into travisCIUsecase") {
		// Add cache for previousStatus
		tUsecase.buildsCache.Add(fmt.Sprintf("%s : #%s", job, branch), "0", SuccessStatus, time.Second*120)

		expected := NewBuildTile(jenkins.JenkinsBuildTileType)
		expected.Label = fmt.Sprintf("%s : #%s", job, branch)
		expected.Status = parseResult(repositoryBuild.Result)
		expected.PreviousStatus = SuccessStatus
		expected.StartedAt = ToTime(repositoryBuild.StartedAt)
		expected.FinishedAt = ToTime(repositoryBuild.StartedAt.Add(repositoryBuild.Duration))
		expected.Author = &Author{
			Name:      repositoryBuild.Author.Name,
			AvatarUrl: repositoryBuild.Author.AvatarUrl,
		}

		tile, err := tu.Build(&models.BuildParams{Job: job, Branch: branch})
		if assert.NoError(t, err) {
			assert.Equal(t, expected, tile)
			mockRepository.AssertNumberOfCalls(t, "GetJob", 1)
			mockRepository.AssertNumberOfCalls(t, "GetLastBuildStatus", 1)
			mockRepository.AssertExpectations(t)
		}
	}
}

func TestBuild_Success(t *testing.T) {
	CheckBuild(t, "SUCCESS")
}

func TestBuild_Unstable(t *testing.T) {
	CheckBuild(t, "UNSTABLE")
}

func TestBuild_Failure(t *testing.T) {
	CheckBuild(t, "FAILURE")
}

func TestBuild_Aborted(t *testing.T) {
	CheckBuild(t, "ABORTED")
}

func TestBuild_Queued(t *testing.T) {
	repositoryJob := &models.Job{
		Buildable: true,
		InQueue:   true,
		QueuedAt:  ToTime(time.Date(2000, 01, 01, 10, 00, 00, 00, time.UTC)),
	}

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetJob", AnythingOfType("string"), AnythingOfType("string")).
		Return(repositoryJob, nil)

	tu := NewJenkinsUsecase(mockRepository, config.Cache{})
	tUsecase, ok := tu.(*jenkinsUsecase)
	if assert.True(t, ok, "enable to case tu into travisCIUsecase") {
		// Add cache for previousStatus
		tUsecase.buildsCache.Add(fmt.Sprintf("%s : #%s", job, branch), "0", SuccessStatus, time.Second*120)

		expected := NewBuildTile(jenkins.JenkinsBuildTileType)
		expected.Label = fmt.Sprintf("%s : #%s", job, branch)
		expected.Status = QueuedStatus
		expected.PreviousStatus = SuccessStatus
		expected.StartedAt = repositoryJob.QueuedAt

		tile, err := tu.Build(&models.BuildParams{Job: job, Branch: branch})
		if assert.NoError(t, err) {
			assert.Equal(t, expected, tile)
			mockRepository.AssertNumberOfCalls(t, "GetJob", 1)
			mockRepository.AssertExpectations(t)
		}
	}
}

func TestBuild_Running(t *testing.T) {
	repositoryJob := &models.Job{
		Buildable: true,
	}
	repositoryBuild := buildResponse("null", time.Now(), 0)
	repositoryBuild.Building = true

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetJob", AnythingOfType("string"), AnythingOfType("string")).
		Return(repositoryJob, nil)
	mockRepository.On("GetLastBuildStatus", Anything).
		Return(repositoryBuild, nil)

	tu := NewJenkinsUsecase(mockRepository, config.Cache{})
	tUsecase, ok := tu.(*jenkinsUsecase)
	if assert.True(t, ok, "enable to case tu into travisCIUsecase") {
		// Without cached build
		expected := NewBuildTile(jenkins.JenkinsBuildTileType)
		expected.Label = fmt.Sprintf("%s : #%s", job, branch)
		expected.Status = RunningStatus
		expected.PreviousStatus = UnknownStatus
		expected.StartedAt = ToTime(repositoryBuild.StartedAt)
		expected.Duration = ToInt64(int64(0))
		expected.EstimatedDuration = ToInt64(int64(0))
		expected.Author = &Author{
			Name:      repositoryBuild.Author.Name,
			AvatarUrl: repositoryBuild.Author.AvatarUrl,
		}

		tile, err := tu.Build(&models.BuildParams{Job: job, Branch: branch})
		if assert.NoError(t, err) {
			assert.Equal(t, expected, tile)
		}

		// With cached build
		tUsecase.buildsCache.Add(fmt.Sprintf("%s : #%s", job, branch), "0", SuccessStatus, time.Second*120)

		expected.PreviousStatus = SuccessStatus
		expected.EstimatedDuration = ToInt64(int64(120))

		tile, err = tu.Build(&models.BuildParams{Job: job, Branch: branch})
		if assert.NoError(t, err) {
			assert.Equal(t, expected, tile)
		}

		mockRepository.AssertNumberOfCalls(t, "GetJob", 2)
		mockRepository.AssertNumberOfCalls(t, "GetLastBuildStatus", 2)
		mockRepository.AssertExpectations(t)
	}
}

func TestListDynamicTile_Success(t *testing.T) {
	repositoryJob := &models.Job{
		ID:        job,
		Buildable: false,
		InQueue:   false,
		Branches:  []string{branch, "develop", "feat%2Ftest-deploy"},
	}

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetJob", AnythingOfType("string"), AnythingOfType("string")).
		Return(repositoryJob, nil)

	tu := NewJenkinsUsecase(mockRepository, config.Cache{})

	tiles, err := tu.ListDynamicTile(&models.MultiBranchParams{Job: job})
	if assert.NoError(t, err) {
		assert.Len(t, tiles, 3)
		assert.Equal(t, jenkins.JenkinsBuildTileType, tiles[0].TileType)
		assert.Equal(t, job, tiles[0].Params["job"])
		assert.Equal(t, "master", tiles[0].Params["branch"])
		assert.Equal(t, jenkins.JenkinsBuildTileType, tiles[1].TileType)
		assert.Equal(t, job, tiles[1].Params["job"])
		assert.Equal(t, "develop", tiles[1].Params["branch"])
		assert.Equal(t, jenkins.JenkinsBuildTileType, tiles[2].TileType)
		assert.Equal(t, job, tiles[2].Params["job"])
		assert.Equal(t, "feat%2Ftest-deploy", tiles[2].Params["branch"])
	}

	tiles, err = tu.ListDynamicTile(&models.MultiBranchParams{Job: job, Filter: "feat/*"})
	if assert.NoError(t, err) {
		assert.Len(t, tiles, 1)
		assert.Equal(t, jenkins.JenkinsBuildTileType, tiles[0].TileType)
		assert.Equal(t, job, tiles[0].Params["job"])
		assert.Equal(t, "feat%2Ftest-deploy", tiles[0].Params["branch"])
	}

	mockRepository.AssertNumberOfCalls(t, "GetJob", 2)
	mockRepository.AssertExpectations(t)
}

func TestListDynamicTile_Timeout(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetJob", AnythingOfType("string"), AnythingOfType("string")).
		Return(nil, errors.New("request canceled"))

	tu := NewJenkinsUsecase(mockRepository, config.Cache{Expire: 120000, CleanupInterval: 1000})

	_, err := tu.ListDynamicTile(&models.MultiBranchParams{Job: job})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "timeout/host unreachable")

	mockRepository.AssertNumberOfCalls(t, "GetJob", 1)
	mockRepository.AssertExpectations(t)
}

func TestListDynamicTile_Timeout_WithCache(t *testing.T) {
	repositoryJob := &models.Job{
		ID:        job,
		Buildable: false,
		InQueue:   false,
		Branches:  []string{branch},
	}

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetJob", AnythingOfType("string"), AnythingOfType("string")).
		Return(nil, errors.New("request canceled"))

	tu := NewJenkinsUsecase(mockRepository, config.Cache{Expire: 120000, CleanupInterval: 1000})
	tUsecase, ok := tu.(*jenkinsUsecase)
	if assert.True(t, ok) {
		tUsecase.jobsCache.Set(job, repositoryJob, 0)

		tiles, err := tu.ListDynamicTile(&models.MultiBranchParams{Job: job})
		assert.NoError(t, err)

		assert.Len(t, tiles, 1)
		assert.Equal(t, jenkins.JenkinsBuildTileType, tiles[0].TileType)
		assert.Equal(t, job, tiles[0].Params["job"])
		assert.Equal(t, "master", tiles[0].Params["branch"])

		mockRepository.AssertNumberOfCalls(t, "GetJob", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestListDynamicTile_Error(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetJob", AnythingOfType("string"), AnythingOfType("string")).
		Return(nil, errors.New("boom"))

	tu := NewJenkinsUsecase(mockRepository, config.Cache{})

	_, err := tu.ListDynamicTile(&models.MultiBranchParams{Job: "test"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unable to found job")

	mockRepository.AssertNumberOfCalls(t, "GetJob", 1)
	mockRepository.AssertExpectations(t)
}

func TestListDynamicTile_ErrorWithRegex(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetJob", AnythingOfType("string"), AnythingOfType("string")).
		Return(nil, nil)

	tu := NewJenkinsUsecase(mockRepository, config.Cache{})

	_, err := tu.ListDynamicTile(&models.MultiBranchParams{Job: "test", Filter: "("})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error parsing regexp")

	mockRepository.AssertNumberOfCalls(t, "GetJob", 1)
	mockRepository.AssertExpectations(t)
}

func TestParseResult(t *testing.T) {
	assert.Equal(t, SuccessStatus, parseResult("SUCCESS"))
	assert.Equal(t, WarningStatus, parseResult("UNSTABLE"))
	assert.Equal(t, FailedStatus, parseResult("FAILURE"))
	assert.Equal(t, AbortedStatus, parseResult("ABORTED"))
	assert.Equal(t, UnknownStatus, parseResult(""))
}

func buildResponse(result string, startedAt time.Time, duration time.Duration) *models.Build {
	repositoryBuild := &models.Build{
		Number:    "1",
		FullName:  "Test-Build",
		Result:    result,
		StartedAt: startedAt,
		Duration:  duration,
		Author: &models.Author{
			Name:      "me",
			AvatarUrl: "http://avatar.com",
		},
	}
	return repositoryBuild
}
