package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/jenkins"
	"github.com/monitoror/monitoror/monitorable/jenkins/mocks"
	jenkinsModels "github.com/monitoror/monitoror/monitorable/jenkins/models"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/git"

	. "github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

var job, branch = "test", "master"

func TestBuild_Error(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetJob", AnythingOfType("string"), AnythingOfType("string")).
		Return(nil, errors.New("boom"))

	tu := NewJenkinsUsecase(mockRepository)

	tile, err := tu.Build(&jenkinsModels.BuildParams{Job: job, Branch: branch})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &models.MonitororError{}, err)
		assert.Equal(t, "unable to find job", err.Error())
		mockRepository.AssertNumberOfCalls(t, "GetJob", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestBuild_DisabledBuild(t *testing.T) {
	repositoryJob := &jenkinsModels.Job{
		Buildable: false,
	}

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetJob", AnythingOfType("string"), AnythingOfType("string")).
		Return(repositoryJob, nil)

	tu := NewJenkinsUsecase(mockRepository)

	tile, err := tu.Build(&jenkinsModels.BuildParams{Job: job})
	if assert.NoError(t, err) {
		assert.Equal(t, job, tile.Label)
		assert.Equal(t, models.DisabledStatus, tile.Status)
		mockRepository.AssertNumberOfCalls(t, "GetJob", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestBuild_Error_NoBuild(t *testing.T) {
	repositoryJob := &jenkinsModels.Job{
		Buildable: true,
	}

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetJob", AnythingOfType("string"), AnythingOfType("string")).
		Return(repositoryJob, nil)
	mockRepository.On("GetLastBuildStatus", Anything).
		Return(nil, errors.New("boom"))

	tu := NewJenkinsUsecase(mockRepository)

	tile, err := tu.Build(&jenkinsModels.BuildParams{Job: job, Branch: branch})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &models.MonitororError{}, err)
		assert.Equal(t, "no build found", err.Error())
		assert.Equal(t, models.UnknownStatus, err.(*models.MonitororError).ErrorStatus)
		mockRepository.AssertNumberOfCalls(t, "GetJob", 1)
		mockRepository.AssertNumberOfCalls(t, "GetLastBuildStatus", 1)
		mockRepository.AssertExpectations(t)
	}
}

func CheckBuild(t *testing.T, result string) {
	repositoryJob := &jenkinsModels.Job{
		Buildable: true,
	}
	repositoryBuild := buildResponse(result, time.Date(2000, 01, 01, 10, 00, 00, 00, time.UTC), time.Minute)

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetJob", AnythingOfType("string"), AnythingOfType("string")).
		Return(repositoryJob, nil)
	mockRepository.On("GetLastBuildStatus", Anything).
		Return(repositoryBuild, nil)

	tu := NewJenkinsUsecase(mockRepository)
	tUsecase, ok := tu.(*jenkinsUsecase)
	if assert.True(t, ok, "enable to case tu into travisCIUsecase") {
		expected := models.NewTile(jenkins.JenkinsBuildTileType).WithBuild()
		expected.Label = job
		expected.Build.ID = ToString("1")
		expected.Build.Branch = ToString(git.HumanizeBranch(branch))

		expected.Status = parseResult(repositoryBuild.Result)
		expected.Build.PreviousStatus = models.SuccessStatus
		expected.Build.StartedAt = ToTime(repositoryBuild.StartedAt)
		expected.Build.FinishedAt = ToTime(repositoryBuild.StartedAt.Add(repositoryBuild.Duration))

		if result == "FAILURE" {
			expected.Build.Author = &models.Author{
				Name:      repositoryBuild.Author.Name,
				AvatarURL: repositoryBuild.Author.AvatarURL,
			}
		}

		// Add cache for previousStatus
		params := &jenkinsModels.BuildParams{Job: job, Branch: branch}
		tUsecase.buildsCache.Add(params, "0", models.SuccessStatus, time.Second*120)
		tile, err := tu.Build(params)
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
	repositoryJob := &jenkinsModels.Job{
		Buildable: true,
		InQueue:   true,
		QueuedAt:  ToTime(time.Date(2000, 01, 01, 10, 00, 00, 00, time.UTC)),
	}

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetJob", AnythingOfType("string"), AnythingOfType("string")).
		Return(repositoryJob, nil)

	tu := NewJenkinsUsecase(mockRepository)
	tUsecase, ok := tu.(*jenkinsUsecase)
	if assert.True(t, ok, "enable to case tu into travisCIUsecase") {
		expected := models.NewTile(jenkins.JenkinsBuildTileType).WithBuild()
		expected.Label = job
		expected.Build.Branch = ToString(git.HumanizeBranch(branch))

		expected.Status = models.QueuedStatus
		expected.Build.PreviousStatus = models.SuccessStatus
		expected.Build.StartedAt = repositoryJob.QueuedAt

		// Add cache for previousStatus
		params := &jenkinsModels.BuildParams{Job: job, Branch: branch}
		tUsecase.buildsCache.Add(params, "0", models.SuccessStatus, time.Second*120)
		tile, err := tu.Build(params)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, tile)
			mockRepository.AssertNumberOfCalls(t, "GetJob", 1)
			mockRepository.AssertExpectations(t)
		}
	}
}

func TestBuild_Running(t *testing.T) {
	repositoryJob := &jenkinsModels.Job{
		Buildable: true,
	}
	repositoryBuild := buildResponse("null", time.Now(), 0)
	repositoryBuild.Building = true

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetJob", AnythingOfType("string"), AnythingOfType("string")).
		Return(repositoryJob, nil)
	mockRepository.On("GetLastBuildStatus", Anything).
		Return(repositoryBuild, nil)

	ju := NewJenkinsUsecase(mockRepository)
	jUsecase, ok := ju.(*jenkinsUsecase)
	if assert.True(t, ok, "enable to case ju into jenkinsUsecase") {
		// Without cached build
		expected := models.NewTile(jenkins.JenkinsBuildTileType).WithBuild()
		expected.Label = job
		expected.Build.ID = ToString("1")
		expected.Build.Branch = ToString(git.HumanizeBranch(branch))

		expected.Status = models.RunningStatus
		expected.Build.PreviousStatus = models.UnknownStatus
		expected.Build.StartedAt = ToTime(repositoryBuild.StartedAt)
		expected.Build.Duration = ToInt64(int64(0))
		expected.Build.EstimatedDuration = ToInt64(int64(0))

		params := &jenkinsModels.BuildParams{Job: job, Branch: branch}
		tile, err := ju.Build(params)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, tile)
		}

		// With cached build
		jUsecase.buildsCache.Add(params, "0", models.SuccessStatus, time.Second*120)

		expected.Build.PreviousStatus = models.SuccessStatus
		expected.Build.EstimatedDuration = ToInt64(int64(120))

		tile, err = ju.Build(params)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, tile)
		}

		mockRepository.AssertNumberOfCalls(t, "GetJob", 2)
		mockRepository.AssertNumberOfCalls(t, "GetLastBuildStatus", 2)
		mockRepository.AssertExpectations(t)
	}
}

func TestListDynamicTile_Success(t *testing.T) {
	repositoryJob := &jenkinsModels.Job{
		ID:        job,
		Buildable: false,
		InQueue:   false,
		Branches:  []string{branch, "develop", "feat%2Ftest-deploy"},
	}

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetJob", AnythingOfType("string"), AnythingOfType("string")).
		Return(repositoryJob, nil)

	tu := NewJenkinsUsecase(mockRepository)

	tiles, err := tu.ListDynamicTile(&jenkinsModels.MultiBranchParams{Job: job})
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

	tiles, err = tu.ListDynamicTile(&jenkinsModels.MultiBranchParams{Job: job, Match: "feat/*"})
	if assert.NoError(t, err) {
		assert.Len(t, tiles, 1)
		assert.Equal(t, jenkins.JenkinsBuildTileType, tiles[0].TileType)
		assert.Equal(t, job, tiles[0].Params["job"])
		assert.Equal(t, "feat%2Ftest-deploy", tiles[0].Params["branch"])
	}

	mockRepository.AssertNumberOfCalls(t, "GetJob", 2)
	mockRepository.AssertExpectations(t)
}

func TestListDynamicTile_Error(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetJob", AnythingOfType("string"), AnythingOfType("string")).
		Return(nil, errors.New("boom"))

	tu := NewJenkinsUsecase(mockRepository)

	_, err := tu.ListDynamicTile(&jenkinsModels.MultiBranchParams{Job: "test"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unable to find job")

	mockRepository.AssertNumberOfCalls(t, "GetJob", 1)
	mockRepository.AssertExpectations(t)
}

func TestListDynamicTile_ErrorWithRegex(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetJob", AnythingOfType("string"), AnythingOfType("string")).
		Return(nil, nil)

	tu := NewJenkinsUsecase(mockRepository)

	_, err := tu.ListDynamicTile(&jenkinsModels.MultiBranchParams{Job: "test", Match: "("})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error parsing regexp")

	_, err = tu.ListDynamicTile(&jenkinsModels.MultiBranchParams{Job: "test", Unmatch: "("})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error parsing regexp")

	mockRepository.AssertNumberOfCalls(t, "GetJob", 2)
	mockRepository.AssertExpectations(t)
}

func TestParseResult(t *testing.T) {
	assert.Equal(t, models.SuccessStatus, parseResult("SUCCESS"))
	assert.Equal(t, models.WarningStatus, parseResult("UNSTABLE"))
	assert.Equal(t, models.FailedStatus, parseResult("FAILURE"))
	assert.Equal(t, models.CanceledStatus, parseResult("ABORTED"))
	assert.Equal(t, models.UnknownStatus, parseResult(""))
}

func buildResponse(result string, startedAt time.Time, duration time.Duration) *jenkinsModels.Build {
	repositoryBuild := &jenkinsModels.Build{
		Number:    "1",
		FullName:  "Test-Build",
		Result:    result,
		StartedAt: startedAt,
		Duration:  duration,
		Author: &models.Author{
			Name:      "me",
			AvatarURL: "http://avatar.com",
		},
	}
	return repositoryBuild
}
