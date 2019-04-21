//+build !faker

package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	. "github.com/AlekSi/pointer"

	"github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/models/errors"

	. "github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/monitorable/travisci"
	"github.com/monitoror/monitoror/monitorable/travisci/model"
)

type (
	travisCIUsecase struct {
		config     *config.Config
		repository travisci.Repository

		// Estimated duration cache
		estimatedDurations map[string]time.Duration
	}
)

func NewTravisCIUsecase(conf *config.Config, repository travisci.Repository) travisci.Usecase {
	return &travisCIUsecase{conf, repository, make(map[string]time.Duration)}
}

func (tu *travisCIUsecase) Build(params *model.BuildParams) (tile *BuildTile, err error) {
	tile = NewBuildTile(travisci.TravisCIBuildTileSubType)
	tile.Label = fmt.Sprintf("%s : #%s", params.Repository, params.Branch)

	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, time.Duration(tu.config.Monitorable.TravisCI.Timeout)*time.Millisecond)

	// Request
	build, err := tu.repository.Build(ctx, params.Group, params.Repository, params.Branch)
	if err != nil {
		if err == context.DeadlineExceeded || strings.Contains(err.Error(), "no such host") || strings.Contains(err.Error(), "dial tcp: lookup") {
			err = errors.NewTimeoutError(tile.Tile, "Timeout/Host Unreachable")
		} else {
			err = errors.NewSystemError("unable to get travisci build", nil)
		}
		return nil, err
	}
	if build == nil {
		err = errors.NewNoBuildError(tile)
		return nil, err
	}

	// Parsing to BuildTile
	tile.Status = parseState(build.State)
	if !build.StartedAt.IsZero() {
		tile.StartedAt = ToInt64(build.StartedAt.Unix())
	}
	if !build.FinishedAt.IsZero() {
		tile.FinishedAt = ToInt64(build.FinishedAt.Unix())
	}

	if tile.Status == RunningStatus || tile.Status == QueuedStatus {
		tile.PreviousStatus = parseState(build.PreviousState)
	}

	if tile.Status == RunningStatus {
		tile.Duration = ToInt64(int64(time.Now().Sub(build.StartedAt).Seconds()))

		// Use cached estimated duration
		if estimatedDuration, ok := tu.estimatedDurations[tile.Label]; ok {
			tile.EstimatedDuration = ToInt64(int64(estimatedDuration / time.Second))
		}
	}

	if build.Author.Name != "" || build.Author.AvatarUrl != "" {
		tile.Author = &Author{
			Name:      build.Author.Name,
			AvatarUrl: build.Author.AvatarUrl,
		}
	}

	// Cache Duration when success
	if tile.Status == SuccessStatus {
		tu.estimatedDurations[tile.Label] = build.Duration
	}

	return
}

func parseState(state string) TileStatus {
	switch state {
	case "created":
		return QueuedStatus
	case "received":
		return QueuedStatus
	case "started":
		return RunningStatus
	case "passed":
		return SuccessStatus
	case "failed":
		return FailedStatus
	case "errored":
		return FailedStatus
	default:
		return UnknownStatus
	}
}
