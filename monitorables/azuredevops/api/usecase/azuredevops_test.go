package usecase

import (
	"errors"
	"testing"
	"time"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/azuredevops/api"
	"github.com/monitoror/monitoror/monitorables/azuredevops/api/mocks"
	"github.com/monitoror/monitoror/monitorables/azuredevops/api/models"

	. "github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAzureDevOpsUsecase_Build_ErrorOnGetBuild(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetBuild", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("GetBuildError"))

	usecase := NewAzureDevOpsUsecase(mockRepository)
	tile, err := usecase.Build(&models.BuildParams{Project: "test", Definition: ToInt(1), Branch: ToString("master")})

	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.Equal(t, "unable to find build", err.Error())
		mockRepository.AssertNumberOfCalls(t, "GetBuild", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestAzureDevOpsUsecase_Build_ErrorNoBuildFound(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetBuild", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

	usecase := NewAzureDevOpsUsecase(mockRepository)
	tile, err := usecase.Build(&models.BuildParams{Project: "test", Definition: ToInt(1), Branch: ToString("master")})

	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.Equal(t, "no build found", err.Error())
		mockRepository.AssertNumberOfCalls(t, "GetBuild", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestAzureDevOpsUsecase_Build_Success(t *testing.T) {
	now := time.Now()

	build := &models.Build{
		BuildNumber:    "1",
		DefinitionName: "definitionName",
		Branch:         "master",
		Author: &coreModels.Author{
			Name:      "test",
			AvatarURL: "monitoror.example.com",
		},
		Status:     "completed",
		Result:     "succeeded",
		FinishedAt: &now,
		StartedAt:  &now,
	}

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetBuild", mock.Anything, mock.Anything, mock.Anything).Return(build, nil)

	expected := coreModels.NewTile(api.AzureDevOpsBuildTileType).WithBuild()
	expected.Label = "test (definitionName)"
	expected.Build.ID = ToString("1")
	expected.Build.Branch = ToString("master")

	expected.Status = coreModels.SuccessStatus
	expected.Build.PreviousStatus = coreModels.UnknownStatus
	expected.Build.StartedAt = &now
	expected.Build.FinishedAt = &now

	params := &models.BuildParams{Project: "test", Definition: ToInt(1), Branch: ToString("master")}

	usecase := NewAzureDevOpsUsecase(mockRepository)
	tile, err := usecase.Build(params)
	if assert.NoError(t, err) {
		assert.NotNil(t, tile)
		assert.Equal(t, expected, tile)
		mockRepository.AssertNumberOfCalls(t, "GetBuild", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestAzureDevOpsUsecase_Build_Failed(t *testing.T) {
	now := time.Now()

	build := &models.Build{
		BuildNumber:    "1",
		DefinitionName: "definitionName",
		Branch:         "master",
		Author: &coreModels.Author{
			Name:      "test",
			AvatarURL: "monitoror.example.com",
		},
		Status:     "completed",
		Result:     "failed",
		FinishedAt: &now,
		StartedAt:  &now,
	}

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetBuild", mock.Anything, mock.Anything, mock.Anything).Return(build, nil)

	expected := coreModels.NewTile(api.AzureDevOpsBuildTileType).WithBuild()
	expected.Label = "test (definitionName)"
	expected.Build.ID = ToString("1")
	expected.Build.Branch = ToString("master")

	expected.Build.Author = &coreModels.Author{
		Name:      "test",
		AvatarURL: "monitoror.example.com",
	}
	expected.Status = coreModels.FailedStatus
	expected.Build.PreviousStatus = coreModels.UnknownStatus
	expected.Build.StartedAt = &now
	expected.Build.FinishedAt = &now

	params := &models.BuildParams{Project: "test", Definition: ToInt(1), Branch: ToString("master")}

	usecase := NewAzureDevOpsUsecase(mockRepository)
	tile, err := usecase.Build(params)
	if assert.NoError(t, err) {
		assert.NotNil(t, tile)
		assert.Equal(t, expected, tile)
		mockRepository.AssertNumberOfCalls(t, "GetBuild", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestAzureDevOpsUsecase_Build_Running(t *testing.T) {
	now := time.Now()

	build := &models.Build{
		BuildNumber:    "1",
		DefinitionName: "definitionName",
		Branch:         "master",
		Author:         nil,
		Status:         "inProgress",
		Result:         "",
		StartedAt:      &now,
	}

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetBuild", mock.Anything, mock.Anything, mock.Anything).Return(build, nil)

	au := NewAzureDevOpsUsecase(mockRepository)
	aUsecase, ok := au.(*azureDevOpsUsecase)
	if assert.True(t, ok, "enable to case au into azureDevOpsUsecase") {
		expected := coreModels.NewTile(api.AzureDevOpsBuildTileType).WithBuild()
		expected.Label = "test (definitionName)"
		expected.Build.ID = ToString("1")
		expected.Build.Branch = ToString("master")

		expected.Status = coreModels.RunningStatus
		expected.Build.PreviousStatus = coreModels.UnknownStatus
		expected.Build.StartedAt = &now
		expected.Build.Duration = ToInt64(0)
		expected.Build.EstimatedDuration = ToInt64(0)

		params := &models.BuildParams{Project: "test", Definition: ToInt(1), Branch: ToString("master")}
		tile, err := au.Build(params)
		if assert.NoError(t, err) {
			assert.NotNil(t, tile)
			assert.Equal(t, expected, tile)
		}

		// With cached build
		aUsecase.buildsCache.Add(params, "0", coreModels.SuccessStatus, time.Second*120)
		expected.Build.PreviousStatus = coreModels.SuccessStatus
		expected.Build.EstimatedDuration = ToInt64(int64(120))

		tile, err = au.Build(params)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, tile)
		}

		mockRepository.AssertNumberOfCalls(t, "GetBuild", 2)
		mockRepository.AssertExpectations(t)
	}
}

func TestAzureDevOpsUsecase_Build_Queued(t *testing.T) {
	now := time.Now()

	build := &models.Build{
		BuildNumber:    "1",
		DefinitionName: "definitionName",
		Branch:         "master",
		Author:         nil,
		Status:         "notStarted",
		Result:         "",
		QueuedAt:       &now,
	}

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetBuild", mock.Anything, mock.Anything, mock.Anything).Return(build, nil)

	au := NewAzureDevOpsUsecase(mockRepository)
	expected := coreModels.NewTile(api.AzureDevOpsBuildTileType).WithBuild()
	expected.Label = "test (definitionName)"
	expected.Build.ID = ToString("1")
	expected.Build.Branch = ToString("master")

	expected.Status = coreModels.QueuedStatus
	expected.Build.PreviousStatus = coreModels.UnknownStatus
	expected.Build.StartedAt = &now

	params := &models.BuildParams{Project: "test", Definition: ToInt(1), Branch: ToString("master")}
	tile, err := au.Build(params)
	if assert.NoError(t, err) {
		assert.NotNil(t, tile)
		assert.Equal(t, expected, tile)
		mockRepository.AssertNumberOfCalls(t, "GetBuild", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestAzureDevOpsUsecase_Release_ErrorOnGetRelease(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetRelease", mock.Anything, mock.Anything).Return(nil, errors.New("GetReleaseError"))

	usecase := NewAzureDevOpsUsecase(mockRepository)
	tile, err := usecase.Release(&models.ReleaseParams{Project: "test", Definition: ToInt(1)})

	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.Equal(t, "unable to find release", err.Error())
		mockRepository.AssertNumberOfCalls(t, "GetRelease", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestAzureDevOpsUsecase_Release_ErrorNoBuildFound(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetRelease", mock.Anything, mock.Anything).Return(nil, nil)

	usecase := NewAzureDevOpsUsecase(mockRepository)
	tile, err := usecase.Release(&models.ReleaseParams{Project: "test", Definition: ToInt(1)})

	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.Equal(t, "no release found", err.Error())
		mockRepository.AssertNumberOfCalls(t, "GetRelease", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestAzureDevOpsUsecase_Release_Success(t *testing.T) {
	now := time.Now()

	release := &models.Release{
		ReleaseNumber:  "1",
		DefinitionName: "definitionName",
		Author: &coreModels.Author{
			Name:      "test",
			AvatarURL: "monitoror.example.com",
		},
		Status:     "succeeded",
		FinishedAt: &now,
		StartedAt:  &now,
	}

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetRelease", mock.Anything, mock.Anything).Return(release, nil)

	expected := coreModels.NewTile(api.AzureDevOpsReleaseTileType).WithBuild()
	expected.Label = "test (definitionName)"
	expected.Build.ID = ToString("1")

	expected.Status = coreModels.SuccessStatus
	expected.Build.PreviousStatus = coreModels.UnknownStatus
	expected.Build.StartedAt = &now
	expected.Build.FinishedAt = &now

	params := &models.ReleaseParams{Project: "test", Definition: ToInt(1)}

	usecase := NewAzureDevOpsUsecase(mockRepository)
	tile, err := usecase.Release(params)
	if assert.NoError(t, err) {
		assert.NotNil(t, tile)
		assert.Equal(t, expected, tile)
		mockRepository.AssertNumberOfCalls(t, "GetRelease", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestAzureDevOpsUsecase_Release_Failed(t *testing.T) {
	now := time.Now()

	release := &models.Release{
		ReleaseNumber:  "1",
		DefinitionName: "definitionName",
		Author: &coreModels.Author{
			Name:      "test",
			AvatarURL: "monitoror.example.com",
		},
		Status:     "failed",
		FinishedAt: &now,
		StartedAt:  &now,
	}

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetRelease", mock.Anything, mock.Anything).Return(release, nil)

	expected := coreModels.NewTile(api.AzureDevOpsReleaseTileType).WithBuild()
	expected.Label = "test (definitionName)"
	expected.Build.ID = ToString("1")

	expected.Build.Author = &coreModels.Author{
		Name:      "test",
		AvatarURL: "monitoror.example.com",
	}
	expected.Status = coreModels.FailedStatus
	expected.Build.PreviousStatus = coreModels.UnknownStatus
	expected.Build.StartedAt = &now
	expected.Build.FinishedAt = &now

	params := &models.ReleaseParams{Project: "test", Definition: ToInt(1)}

	usecase := NewAzureDevOpsUsecase(mockRepository)
	tile, err := usecase.Release(params)
	if assert.NoError(t, err) {
		assert.NotNil(t, tile)
		assert.Equal(t, expected, tile)
		mockRepository.AssertNumberOfCalls(t, "GetRelease", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestAzureDevOpsUsecase_Release_Running(t *testing.T) {
	now := time.Now()

	release := &models.Release{
		ReleaseNumber:  "1",
		DefinitionName: "definitionName",
		Author:         nil,
		Status:         "inProgress",
		StartedAt:      &now,
	}

	mockRepository := new(mocks.Repository)
	mockRepository.On("GetRelease", mock.Anything, mock.Anything).Return(release, nil)

	au := NewAzureDevOpsUsecase(mockRepository)
	aUsecase, ok := au.(*azureDevOpsUsecase)
	if assert.True(t, ok, "enable to case au into azureDevOpsUsecase") {
		expected := coreModels.NewTile(api.AzureDevOpsReleaseTileType).WithBuild()
		expected.Label = "test (definitionName)"
		expected.Build.ID = ToString("1")

		expected.Status = coreModels.RunningStatus
		expected.Build.PreviousStatus = coreModels.UnknownStatus
		expected.Build.StartedAt = &now
		expected.Build.Duration = ToInt64(0)
		expected.Build.EstimatedDuration = ToInt64(0)

		params := &models.ReleaseParams{Project: "test", Definition: ToInt(1)}
		tile, err := au.Release(params)
		if assert.NoError(t, err) {
			assert.NotNil(t, tile)
			assert.Equal(t, expected, tile)
		}

		// With cached build
		aUsecase.buildsCache.Add(params, "0", coreModels.SuccessStatus, time.Second*120)
		expected.Build.PreviousStatus = coreModels.SuccessStatus
		expected.Build.EstimatedDuration = ToInt64(int64(120))

		tile, err = au.Release(params)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, tile)
		}

		mockRepository.AssertNumberOfCalls(t, "GetRelease", 2)
		mockRepository.AssertExpectations(t)
	}
}

func Test_parseBuildResult(t *testing.T) {
	assert.Equal(t, coreModels.RunningStatus, parseBuildResult("inProgress", ""))
	assert.Equal(t, coreModels.RunningStatus, parseBuildResult("cancelling", ""))
	assert.Equal(t, coreModels.QueuedStatus, parseBuildResult("notStarted", ""))
	assert.Equal(t, coreModels.SuccessStatus, parseBuildResult("completed", "succeeded"))
	assert.Equal(t, coreModels.WarningStatus, parseBuildResult("completed", "partiallySucceeded"))
	assert.Equal(t, coreModels.FailedStatus, parseBuildResult("completed", "failed"))
	assert.Equal(t, coreModels.CanceledStatus, parseBuildResult("completed", "canceled"))
	assert.Equal(t, coreModels.UnknownStatus, parseBuildResult("completed", ""))
	assert.Equal(t, coreModels.UnknownStatus, parseBuildResult("", ""))
}

func Test_parseReleaseStatus(t *testing.T) {
	assert.Equal(t, coreModels.FailedStatus, parseReleaseStatus("failed"))
	assert.Equal(t, coreModels.SuccessStatus, parseReleaseStatus("succeeded"))
	assert.Equal(t, coreModels.WarningStatus, parseReleaseStatus("partiallySucceeded"))
	assert.Equal(t, coreModels.RunningStatus, parseReleaseStatus("inProgress"))
	assert.Equal(t, coreModels.UnknownStatus, parseReleaseStatus("all"))
	assert.Equal(t, coreModels.UnknownStatus, parseReleaseStatus("notDeployed"))
	assert.Equal(t, coreModels.UnknownStatus, parseReleaseStatus(""))
}
