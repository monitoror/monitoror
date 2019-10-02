//+build !faker

package usecase

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/jsdidierlaurent/echo-middleware/cache"
	"github.com/monitoror/monitoror/config"
	. "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/pingdom"
	"github.com/monitoror/monitoror/monitorable/pingdom/models"
	"github.com/monitoror/monitoror/pkg/monitoror/builder"
	uuid "github.com/satori/go.uuid"
)

type (
	pingdomUsecase struct {
		repository pingdom.Repository
		// Used to generate store key by repository
		repositoryUid string

		// Config
		config *config.Pingdom

		// Mutex for lock multi access on Pingdom
		sync.Mutex
		// Used for caching result of pingdom (to avoid bursting query limit)
		store cache.Store
	}
)

const (
	PingdomChecksStoreKeyPrefix   = "monitoror.pingdom.checks.store"
	PingdomTagsByIdStoreKeyPrefix = "monitoror.pingdom.tagsById.store"
	PingdomCheckStoreKeyPrefix    = "monitoror.pingdom.check.store"

	PausedStatus = "paused"
	UpStatus     = "up"
	DownStatus   = "down"
)

func NewPingdomUsecase(repository pingdom.Repository, config *config.Pingdom, store cache.Store) pingdom.Usecase {
	return &pingdomUsecase{
		repository:    repository,
		repositoryUid: uuid.NewV4().String(),
		config:        config,
		store:         store,
	}
}

func (pu *pingdomUsecase) Check(params *models.CheckParams) (*Tile, error) {
	tile := NewTile(pingdom.PingdomCheckTileType)

	checkId := *params.Id
	var result models.Check
	var tags string

	// Lookup in store for bulk query for this ID, if found, use it
	if err := pu.store.Get(pu.getTagsByIdStoreKey(checkId), &tags); err == nil {
		checks, err := pu.loadChecks(tags)
		if err != nil {
			return nil, &MonitororError{Err: err, Message: "unable to found checks"}
		}

		// Find check in array
		for _, tmpCheck := range checks {
			if tmpCheck.Id == checkId {
				result = tmpCheck
				break
			}
		}
	} else // Bulk not found, request single check
	{
		check, err := pu.loadCheck(checkId)
		if err != nil {
			return nil, &MonitororError{Err: err, Message: "unable to found check"}
		}
		result = *check
	}

	// Parse result to tile
	tile.Label = result.Name
	tile.Status = parseStatus(result.Status)

	return tile, nil
}

func (pu *pingdomUsecase) ListDynamicTile(params interface{}) (results []builder.Result, err error) {
	lcParams := params.(*models.ChecksParams)

	checks, err := pu.loadChecks(lcParams.Tags)
	if err != nil {
		return nil, &MonitororError{Err: err, Message: "unable to list checks"}
	}

	if lcParams.SortBy == "name" {
		sort.SliceStable(checks, func(i, j int) bool { return checks[i].Name < checks[j].Name })
	}

	for _, check := range checks {
		// Adding id -> tags in the store for one minute. This value will be refresh each time we call this route
		// This store will be use to find the best route to call for loading check result.
		_ = pu.store.Set(pu.getTagsByIdStoreKey(check.Id), lcParams.Tags, time.Minute)

		// Build results
		if check.Status != PausedStatus {
			p := make(map[string]interface{})
			p["id"] = check.Id

			results = append(results, builder.Result{
				TileType: pingdom.PingdomCheckTileType,
				Label:    check.Name,
				Params:   p,
			})
		}
	}

	return
}

func (pu *pingdomUsecase) loadCheck(id int) (result *models.Check, err error) {
	// Synchronize to avoid multi call on pingdom api
	pu.Lock()
	defer pu.Unlock()

	// Lookup in cache
	result = &models.Check{}
	key := pu.getCheckStoreKey(id)
	if err = pu.store.Get(key, result); err == nil {
		// Cache found, return
		return
	}

	result, err = pu.repository.GetCheck(id)
	if err != nil {
		return
	}

	// Adding result in store
	_ = pu.store.Set(key, *result, time.Millisecond*time.Duration(pu.config.CacheExpiration))

	return
}

func (pu *pingdomUsecase) loadChecks(tags string) (results []models.Check, err error) {
	// Synchronize to avoid multi call on pingdom api
	pu.Lock()
	defer pu.Unlock()

	// Lookup in cache
	key := pu.getChecksStoreKey(tags)
	if err = pu.store.Get(key, &results); err == nil {
		// Cache found, return
		return
	}

	results, err = pu.repository.GetChecks(tags)
	if err != nil {
		return
	}

	// Adding result in store
	_ = pu.store.Set(key, results, time.Millisecond*time.Duration(pu.config.CacheExpiration))

	return
}

func (pu *pingdomUsecase) getChecksStoreKey(tags string) string {
	return fmt.Sprintf("%s:%s-%s", PingdomChecksStoreKeyPrefix, pu.repositoryUid, tags)
}

func (pu *pingdomUsecase) getCheckStoreKey(id int) string {
	return fmt.Sprintf("%s:%s-%d", PingdomCheckStoreKeyPrefix, pu.repositoryUid, id)
}

func (pu *pingdomUsecase) getTagsByIdStoreKey(id int) string {
	return fmt.Sprintf("%s:%s-%d", PingdomTagsByIdStoreKeyPrefix, pu.repositoryUid, id)
}

func parseStatus(status string) TileStatus {
	switch status {
	case UpStatus:
		return SuccessStatus
	case DownStatus:
		return FailedStatus
	case PausedStatus:
		return DisabledStatus
	default:
		return UnknownStatus
	}
}
