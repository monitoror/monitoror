package repository

import (
	"errors"
	"testing"
	"time"

	"github.com/monitoror/monitoror/monitorables/azuredevops/api/mocks"
	"github.com/monitoror/monitoror/monitorables/azuredevops/api/models"
	"github.com/monitoror/monitoror/monitorables/azuredevops/config"
	mocksBuild "github.com/monitoror/monitoror/pkg/goazuredevops/build/mocks"
	mocksRelease "github.com/monitoror/monitoror/pkg/goazuredevops/release/mocks"

	. "github.com/AlekSi/pointer"
	"github.com/jsdidierlaurent/azure-devops-go-api/azuredevops"
	"github.com/jsdidierlaurent/azure-devops-go-api/azuredevops/build"
	"github.com/jsdidierlaurent/azure-devops-go-api/azuredevops/release"
	"github.com/jsdidierlaurent/azure-devops-go-api/azuredevops/webapi"
	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func initRepository(t *testing.T, buildClient build.Client, releaseClient release.Client) *azureDevOpsRepository {
	conf := &config.AzureDevOps{
		URL:             "http://azure.example.com/",
		Token:           "test",
		Timeout:         1000,
		InitialMaxDelay: 1000,
	}

	mockConnection := new(mocks.Connection)
	if buildClient != nil {
		mockConnection.On("GetBuildConnection").Return(buildClient, nil)
	} else {
		mockConnection.On("GetBuildConnection").
			Return(nil, errors.New("GetBuildConnectionError"))
	}

	if releaseClient != nil {
		mockConnection.On("GetReleaseConnection").Return(releaseClient, nil)
	} else {
		mockConnection.On("GetReleaseConnection").
			Return(nil, errors.New("GetReleaseConnectionError"))
	}

	repository := NewAzureDevOpsRepository(conf)

	assert.Equal(t, "http://azure.example.com", conf.URL)

	apiAzureDevOpsRepository, ok := repository.(*azureDevOpsRepository)
	if assert.True(t, ok) {
		apiAzureDevOpsRepository.connection = mockConnection
		return apiAzureDevOpsRepository
	}
	return nil
}

func TestConnection_GetBuildConnection(t *testing.T) {
	// Fake connection, just fort testing if NewClient is call correctly
	con := &connection{&azuredevops.Connection{}}
	client, err := con.GetBuildConnection()
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestConnection_GetReleaseConnection(t *testing.T) {
	// Fake connection, just fort testing if NewClient is call correctly
	con := &connection{&azuredevops.Connection{}}
	client, err := con.GetReleaseConnection()
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestRepository_GetBuild_Failure_ErrorOnGetClient(t *testing.T) {
	repository := initRepository(t, nil, nil)
	_, err := repository.GetBuild("test", 1, ToString("master"))

	assert.Error(t, err)
	assert.Equal(t, "GetBuildConnectionError", err.Error())
}

func TestRepository_GetBuild_Failure_ErrorOnGetBuilds(t *testing.T) {
	mockBuild := new(mocksBuild.Client)
	mockBuild.On("GetBuilds", Anything, AnythingOfType("build.GetBuildsArgs")).
		Return(nil, errors.New("GetBuildsError"))

	repository := initRepository(t, mockBuild, nil)
	_, err := repository.GetBuild("test", 1, ToString("master"))

	if assert.Error(t, err) {
		assert.Equal(t, "GetBuildsError", err.Error())
		mockBuild.AssertNumberOfCalls(t, "GetBuilds", 1)
		mockBuild.AssertExpectations(t)
	}
}

func TestRepository_GetBuild_Failure_NoBuildsFound(t *testing.T) {
	azureDevOpsBuild := &build.GetBuildsResponseValue{
		Value: []build.Build{},
	}

	mockBuild := new(mocksBuild.Client)
	mockBuild.On("GetBuilds", Anything, AnythingOfType("build.GetBuildsArgs")).
		Return(azureDevOpsBuild, nil)

	repository := initRepository(t, mockBuild, nil)
	bu, err := repository.GetBuild("test", 1, ToString("master"))

	if assert.NoError(t, err) {
		assert.Nil(t, bu)
		mockBuild.AssertNumberOfCalls(t, "GetBuilds", 1)
		mockBuild.AssertExpectations(t)
	}
}

func TestRepository_GetBuild_Success(t *testing.T) {
	author := make(map[string]string)
	author["pr.sender.name"] = "Microsoft.VisualStudio.Services.TFS"
	author["pr.sender.avatarURL"] = "http://avatar.url.com"

	now := time.Now()

	azureDevOpsBuild := &build.GetBuildsResponseValue{
		Value: []build.Build{
			{
				BuildNumber: ToString("1"),
				Definition: &build.DefinitionReference{
					Name: ToString("definitionName"),
				},
				SourceBranch: ToString("refs/heads/master"),
				Status:       &build.BuildStatusValues.Completed,
				Result:       &build.BuildResultValues.Succeeded,
				TriggerInfo:  &author,
				QueueTime:    &azuredevops.Time{Time: now},
				StartTime:    &azuredevops.Time{Time: now},
				FinishTime:   &azuredevops.Time{Time: now},
			},
		},
	}

	mockBuild := new(mocksBuild.Client)
	mockBuild.On("GetBuilds", Anything, AnythingOfType("build.GetBuildsArgs")).
		Return(azureDevOpsBuild, nil)

	// Expected
	expectedBuild := &models.Build{
		BuildNumber:    "1",
		DefinitionName: "definitionName",
		Branch:         "refs/heads/master",
		Author:         nil,
		Status:         "completed",
		Result:         "succeeded",
		FinishedAt:     &now,
		StartedAt:      &now,
		QueuedAt:       &now,
	}

	repository := initRepository(t, mockBuild, nil)
	if repository != nil {
		b, err := repository.GetBuild("test", 1, ToString("master"))
		assert.NoError(t, err)
		assert.Equal(t, expectedBuild, b)
		mockBuild.AssertNumberOfCalls(t, "GetBuilds", 1)
		mockBuild.AssertExpectations(t)
	}
}

func TestRepository_GetBuild_Success_WithoutAuthor(t *testing.T) {
	links := make(map[string]webapi.ReferenceLink)
	links["avatar"] = webapi.ReferenceLink{Href: ToString("http://avatar.url.com")}

	now := time.Now()

	azureDevOpsBuild := &build.GetBuildsResponseValue{
		Value: []build.Build{
			{
				BuildNumber: ToString("1"),
				Definition: &build.DefinitionReference{
					Name: ToString("definitionName"),
				},
				SourceBranch: ToString("master"),
				Status:       &build.BuildStatusValues.Completed,
				Result:       &build.BuildResultValues.Succeeded,
				TriggerInfo:  nil,
				RequestedFor: &webapi.IdentityRef{
					DisplayName: ToString("Microsoft.VisualStudio.Services.TFS"),
					Links:       links,
				},
				QueueTime:  &azuredevops.Time{Time: now},
				StartTime:  &azuredevops.Time{Time: now},
				FinishTime: &azuredevops.Time{Time: now},
			},
		},
	}

	mockBuild := new(mocksBuild.Client)
	mockBuild.On("GetBuilds", Anything, AnythingOfType("build.GetBuildsArgs")).
		Return(azureDevOpsBuild, nil)

	// Expected
	expectedBuild := &models.Build{
		BuildNumber:    "1",
		DefinitionName: "definitionName",
		Branch:         "master",
		Author:         nil,
		Status:         "completed",
		Result:         "succeeded",
		FinishedAt:     &now,
		StartedAt:      &now,
		QueuedAt:       &now,
	}

	repository := initRepository(t, mockBuild, nil)
	if repository != nil {
		b, err := repository.GetBuild("test", 1, ToString("master"))
		assert.NoError(t, err)
		assert.Equal(t, expectedBuild, b)
		mockBuild.AssertNumberOfCalls(t, "GetBuilds", 1)
		mockBuild.AssertExpectations(t)
	}
}

func TestRepository_GetRelease_Failure_ErrorOnGetClient(t *testing.T) {
	repository := initRepository(t, nil, nil)
	_, err := repository.GetRelease("test", 1)

	assert.Error(t, err)
	assert.Equal(t, "GetReleaseConnectionError", err.Error())
}

func TestRepository_GetRelease_Failure_ErrorOnGetBuilds(t *testing.T) {
	mockRelease := new(mocksRelease.Client)
	mockRelease.On("GetDeployments", Anything, AnythingOfType("release.GetDeploymentsArgs")).
		Return(nil, errors.New("GetDeploymentsError"))

	repository := initRepository(t, nil, mockRelease)
	_, err := repository.GetRelease("test", 1)

	assert.Error(t, err)
	assert.Equal(t, "GetDeploymentsError", err.Error())

	if assert.Error(t, err) {
		assert.Equal(t, "GetDeploymentsError", err.Error())
		mockRelease.AssertNumberOfCalls(t, "GetDeployments", 1)
		mockRelease.AssertExpectations(t)
	}
}

func TestRepository_GetRelease_Failure_NoBuildsFound(t *testing.T) {
	azureDevOpsDeployments := &release.GetDeploymentsResponseValue{
		Value: []release.Deployment{},
	}

	mockRelease := new(mocksRelease.Client)
	mockRelease.On("GetDeployments", Anything, AnythingOfType("release.GetDeploymentsArgs")).
		Return(azureDevOpsDeployments, nil)

	repository := initRepository(t, nil, mockRelease)
	bu, err := repository.GetRelease("test", 1)

	if assert.NoError(t, err) {
		assert.Nil(t, bu)
		mockRelease.AssertNumberOfCalls(t, "GetDeployments", 1)
		mockRelease.AssertExpectations(t)
	}
}

func TestRepository_GetRelease_Success(t *testing.T) {
	links := make(map[string]webapi.ReferenceLink)
	links["avatar"] = webapi.ReferenceLink{Href: ToString("http://avatar.url.com")}

	now := time.Now()

	azureDevOpsDeployments := &release.GetDeploymentsResponseValue{
		Value: []release.Deployment{
			{
				Release: &release.ReleaseReference{
					Name: ToString("1"),
				},
				ReleaseDefinition: &release.ReleaseDefinitionShallowReference{
					Name: ToString("definitionName"),
				},
				DeploymentStatus: &release.DeploymentStatusValues.Succeeded,
				RequestedFor: &webapi.IdentityRef{
					DisplayName: ToString("Microsoft.VisualStudio.Services.TFS"),
					Links:       links,
				},
				QueuedOn:    &azuredevops.Time{Time: now},
				StartedOn:   &azuredevops.Time{Time: now},
				CompletedOn: &azuredevops.Time{Time: now},
			},
		},
	}

	mockRelease := new(mocksRelease.Client)
	mockRelease.On("GetDeployments", Anything, AnythingOfType("release.GetDeploymentsArgs")).
		Return(azureDevOpsDeployments, nil)

	// Expected
	expectedRelease := &models.Release{
		ReleaseNumber:  "1",
		DefinitionName: "definitionName",
		Author:         nil,
		Status:         "succeeded",
		FinishedAt:     &now,
		StartedAt:      &now,
		QueuedAt:       &now,
	}

	repository := initRepository(t, nil, mockRelease)
	if repository != nil {
		r, err := repository.GetRelease("test", 1)
		assert.NoError(t, err)
		assert.Equal(t, expectedRelease, r)
		mockRelease.AssertNumberOfCalls(t, "GetDeployments", 1)
		mockRelease.AssertExpectations(t)
	}
}
