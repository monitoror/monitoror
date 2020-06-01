//+build faker

package usecase

import (
	"fmt"
	"time"

	"github.com/AlekSi/pointer"
	cmap "github.com/orcaman/concurrent-map"

	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
	"github.com/monitoror/monitoror/internal/pkg/monitorable/faker"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/gitlab/api"
	"github.com/monitoror/monitoror/monitorables/gitlab/api/models"
	"github.com/monitoror/monitoror/pkg/git"
	"github.com/monitoror/monitoror/pkg/nonempty"
)

type (
	gitlabUsecase struct {
		timeRefByProject cmap.ConcurrentMap
	}
)

var availableBuildStatus = faker.Statuses{
	{coreModels.SuccessStatus, time.Second * 30},
	{coreModels.FailedStatus, time.Second * 30},
	{coreModels.CanceledStatus, time.Second * 20},
	{coreModels.RunningStatus, time.Second * 60},
	{coreModels.QueuedStatus, time.Second * 30},
	{coreModels.ActionRequiredStatus, time.Second * 20},
}

func NewGitlabUsecase() api.Usecase {
	return &gitlabUsecase{cmap.New()}
}

func (gu *gitlabUsecase) Issues(params *models.IssuesParams) (*coreModels.Tile, error) {
	tile := coreModels.NewTile(api.GitlabIssuesTileType).WithValue(coreModels.NumberUnit)
	tile.Label = "GitLab issues"

	tile.Status = coreModels.SuccessStatus

	if len(params.ValueValues) != 0 {
		tile.Value.Values = params.ValueValues
	} else {
		tile.Value.Values = []string{"42"}
	}

	return tile, nil
}

func (gu *gitlabUsecase) Pipeline(params *models.PipelineParams) (tile *coreModels.Tile, err error) {
	tile = coreModels.NewTile(api.GitlabPipelineTileType).WithBuild()
	tile.Label = fmt.Sprintf("Project %d name", *params.ProjectID)

	projectID := fmt.Sprintf("%d-%s", params.ProjectID, params.Ref)
	tile.Status = nonempty.Struct(params.Status, gu.computeStatus(projectID)).(coreModels.TileStatus)

	tile.Build.Branch = pointer.ToString(git.HumanizeBranch(params.Ref))
	tile.Build.PreviousStatus = nonempty.Struct(params.PreviousStatus, coreModels.SuccessStatus).(coreModels.TileStatus)

	// Author
	if tile.Status == coreModels.FailedStatus {
		tile.Build.Author = &coreModels.Author{}
		tile.Build.Author.Name = nonempty.String(params.AuthorName, "John Doe")
		tile.Build.Author.AvatarURL = nonempty.String(params.AuthorAvatarURL, "https://monitoror.com/assets/images/avatar.png")
	}

	// Duration / EstimatedDuration
	if tile.Status == coreModels.RunningStatus {
		estimatedDuration := nonempty.Duration(time.Duration(params.EstimatedDuration), time.Second*300)
		tile.Build.Duration = pointer.ToInt64(nonempty.Int64(params.Duration, int64(gu.computeDuration(projectID, estimatedDuration).Seconds())))

		if tile.Build.PreviousStatus != coreModels.UnknownStatus {
			tile.Build.EstimatedDuration = pointer.ToInt64(int64(estimatedDuration.Seconds()))
		} else {
			tile.Build.EstimatedDuration = pointer.ToInt64(0)
		}
	}

	// StartedAt / FinishedAt
	if tile.Build.Duration == nil {
		tile.Build.StartedAt = pointer.ToTime(nonempty.Time(params.StartedAt, time.Now().Add(-time.Minute*10)))
	} else {
		tile.Build.StartedAt = pointer.ToTime(nonempty.Time(params.StartedAt, time.Now().Add(-time.Second*time.Duration(*tile.Build.Duration))))
	}

	if tile.Status != coreModels.QueuedStatus && tile.Status != coreModels.RunningStatus {
		tile.Build.FinishedAt = pointer.ToTime(nonempty.Time(params.FinishedAt, tile.Build.StartedAt.Add(time.Minute*5)))
	}

	return tile, nil
}

func (gu *gitlabUsecase) MergeRequest(params *models.MergeRequestParams) (tile *coreModels.Tile, err error) {
	tile = coreModels.NewTile(api.GitlabMergeRequestTileType).WithBuild()
	tile.Label = fmt.Sprintf("Project %d name", *params.ProjectID)

	projectID := fmt.Sprintf("%d-%d", params.ProjectID, params.ID)
	tile.Status = nonempty.Struct(params.Status, gu.computeStatus(projectID)).(coreModels.TileStatus)

	tile.Build.Branch = pointer.ToString(nonempty.String(git.HumanizeBranch(params.Branch), "feature-branch"))
	tile.Build.PreviousStatus = nonempty.Struct(params.PreviousStatus, coreModels.SuccessStatus).(coreModels.TileStatus)
	tile.Build.MergeRequest = &coreModels.TileMergeRequest{
		ID:    *params.ID,
		Title: nonempty.String(params.MergeRequestTitle, "Feature branch title"),
	}

	// Author
	if tile.Status == coreModels.FailedStatus {
		tile.Build.Author = &coreModels.Author{}
		tile.Build.Author.Name = nonempty.String(params.AuthorName, "John Doe")
		tile.Build.Author.AvatarURL = nonempty.String(params.AuthorAvatarURL, "https://monitoror.com/assets/images/avatar.png")
	}

	// Duration / EstimatedDuration
	if tile.Status == coreModels.RunningStatus {
		estimatedDuration := nonempty.Duration(time.Duration(params.EstimatedDuration), time.Second*300)
		tile.Build.Duration = pointer.ToInt64(nonempty.Int64(params.Duration, int64(gu.computeDuration(projectID, estimatedDuration).Seconds())))

		if tile.Build.PreviousStatus != coreModels.UnknownStatus {
			tile.Build.EstimatedDuration = pointer.ToInt64(int64(estimatedDuration.Seconds()))
		} else {
			tile.Build.EstimatedDuration = pointer.ToInt64(0)
		}
	}

	// StartedAt / FinishedAt
	if tile.Build.Duration == nil {
		tile.Build.StartedAt = pointer.ToTime(nonempty.Time(params.StartedAt, time.Now().Add(-time.Minute*10)))
	} else {
		tile.Build.StartedAt = pointer.ToTime(nonempty.Time(params.StartedAt, time.Now().Add(-time.Second*time.Duration(*tile.Build.Duration))))
	}

	if tile.Status != coreModels.QueuedStatus && tile.Status != coreModels.RunningStatus {
		tile.Build.FinishedAt = pointer.ToTime(nonempty.Time(params.FinishedAt, tile.Build.StartedAt.Add(time.Minute*5)))
	}

	return tile, nil
}

func (gu *gitlabUsecase) MergeRequestsGenerator(params interface{}) ([]uiConfigModels.GeneratedTile, error) {
	panic("unimplemented")
}

func (gu *gitlabUsecase) computeStatus(projectUID string) coreModels.TileStatus {
	value, ok := gu.timeRefByProject.Get(projectUID)
	if !ok || value == nil {
		value = faker.GetRefTime()
		gu.timeRefByProject.Set(projectUID, value)
	}

	return faker.ComputeStatus(value.(time.Time), availableBuildStatus)
}

func (gu *gitlabUsecase) computeDuration(projectUID string, duration time.Duration) time.Duration {
	value, ok := gu.timeRefByProject.Get(projectUID)
	if !ok || value == nil {
		value = faker.GetRefTime()
		gu.timeRefByProject.Set(projectUID, value)
	}

	return faker.ComputeDuration(value.(time.Time), duration)
}
