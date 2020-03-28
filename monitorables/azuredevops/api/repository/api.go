package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/azuredevops/api"
	"github.com/monitoror/monitoror/monitorables/azuredevops/api/models"
	"github.com/monitoror/monitoror/monitorables/azuredevops/config"

	"github.com/AlekSi/pointer"
	azureDevOpsApi "github.com/jsdidierlaurent/azure-devops-go-api/azuredevops"
	"github.com/jsdidierlaurent/azure-devops-go-api/azuredevops/build"
	"github.com/jsdidierlaurent/azure-devops-go-api/azuredevops/release"
)

type (
	azureDevOpsRepository struct {
		connection api.Connection
		config     *config.AzureDevOps
	}

	connection struct {
		connection *azureDevOpsApi.Connection
	}
)

func (c *connection) GetBuildConnection() (build.Client, error) {
	return build.NewClient(context.TODO(), c.connection)
}

func (c *connection) GetReleaseConnection() (release.Client, error) {
	return release.NewClient(context.TODO(), c.connection)
}

func NewAzureDevOpsRepository(config *config.AzureDevOps) api.Repository {
	// Remove last /
	config.URL = strings.TrimRight(config.URL, "/")

	conn := azureDevOpsApi.NewPatConnection(config.URL, config.Token)

	// Setup timeout
	timeout := time.Duration(config.Timeout) * time.Millisecond
	conn.Timeout = &timeout

	return &azureDevOpsRepository{
		connection: &connection{conn},
		config:     config,
	}
}

func (r *azureDevOpsRepository) GetBuild(project string, definition int, branch *string) (*models.Build, error) {
	// Inject "refs/heads/" in branch name
	if branch != nil && !strings.HasPrefix(*branch, "refs/") {
		branch = pointer.ToString(fmt.Sprintf("refs/heads/%s", *branch))
	}

	ids := []int{definition}
	args := build.GetBuildsArgs{
		Project:                pointer.ToString(project),
		Definitions:            &ids,
		BranchName:             branch, // Can be nil
		MaxBuildsPerDefinition: pointer.ToInt(1),
	}

	client, err := r.connection.GetBuildConnection()
	if err != nil {
		return nil, err
	}

	aBuilds, err := client.GetBuilds(context.TODO(), args)
	if err != nil {
		return nil, err
	}

	// No build found
	if len(aBuilds.Value) == 0 {
		return nil, nil
	}
	aBuild := aBuilds.Value[0]

	result := &models.Build{
		BuildNumber:    *aBuild.BuildNumber,
		DefinitionName: *aBuild.Definition.Name,
	}

	// Branch
	if aBuild.SourceBranch != nil {
		result.Branch = *aBuild.SourceBranch
	}

	// Status
	if aBuild.Status != nil {
		result.Status = string(*aBuild.Status)
	}

	// DynamicTileResult
	if aBuild.Result != nil {
		result.Result = string(*aBuild.Result)
	}

	// Author
	result.Author = &coreModels.Author{}
	if aBuild.TriggerInfo != nil {
		if value, ok := (*aBuild.TriggerInfo)["pr.sender.name"]; ok {
			result.Author.Name = value
		}
		if value, ok := (*aBuild.TriggerInfo)["pr.sender.avatarURL"]; ok {
			result.Author.AvatarURL = value
		}
	}

	if aBuild.RequestedFor != nil {
		if result.Author.Name == "" && aBuild.RequestedFor.DisplayName != nil {
			result.Author.Name = *aBuild.RequestedFor.DisplayName
		}
		if result.Author.AvatarURL == "" {
			if link, ok := aBuild.RequestedFor.Links["avatar"]; ok {
				result.Author.AvatarURL = *link.Href
			}
		}
	}

	// HACK: Remove author if user is the default Azure user or empty
	if result.Author.Name == "" || result.Author.Name == "Microsoft.VisualStudio.Services.TFS" {
		result.Author = nil
	}

	if aBuild.QueueTime != nil {
		result.QueuedAt = &aBuild.QueueTime.Time
	}
	if aBuild.StartTime != nil {
		result.StartedAt = &aBuild.StartTime.Time
	}
	if aBuild.FinishTime != nil {
		result.FinishedAt = &aBuild.FinishTime.Time
	}

	return result, nil
}

func (r *azureDevOpsRepository) GetRelease(project string, definition int) (*models.Release, error) {
	args := release.GetDeploymentsArgs{
		Project:            pointer.ToString(project),
		DefinitionId:       pointer.ToInt(definition),
		LatestAttemptsOnly: pointer.ToBool(true),
		Top:                pointer.ToInt(1),
	}

	client, err := r.connection.GetReleaseConnection()
	if err != nil {
		return nil, err
	}

	aReleases, err := client.GetDeployments(context.TODO(), args)
	if err != nil {
		return nil, err
	}

	// No build found
	if len(aReleases.Value) == 0 {
		return nil, nil
	}
	aRelease := aReleases.Value[0]

	result := &models.Release{
		ReleaseNumber:  *aRelease.Release.Name,
		DefinitionName: *aRelease.ReleaseDefinition.Name,
		Status:         string(*aRelease.DeploymentStatus),
	}

	// Author
	if aRelease.RequestedFor != nil {
		result.Author = &coreModels.Author{}

		if aRelease.RequestedFor.DisplayName != nil {
			result.Author.Name = *aRelease.RequestedFor.DisplayName
		}
		if link, ok := aRelease.RequestedFor.Links["avatar"]; ok {
			result.Author.AvatarURL = *link.Href
		}

		// HACK: Remove author if user is the default Azure user or empty
		if result.Author.Name == "" || result.Author.Name == "Microsoft.VisualStudio.Services.TFS" {
			result.Author = nil
		}
	}

	if aRelease.QueuedOn != nil {
		result.QueuedAt = &aRelease.QueuedOn.Time
	}
	if aRelease.StartedOn != nil {
		result.StartedAt = &aRelease.StartedOn.Time
	}
	if aRelease.CompletedOn != nil {
		result.FinishedAt = &aRelease.CompletedOn.Time
	}

	return result, nil
}
