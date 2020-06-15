//+build !faker

package usecase

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/AlekSi/pointer"

	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/pingdom/api"
	"github.com/monitoror/monitoror/monitorables/pingdom/api/models"

	"github.com/jsdidierlaurent/echo-middleware/cache"
	uuid "github.com/satori/go.uuid"
)

type (
	pingdomUsecase struct {
		repository api.Repository
		// Used to generate store key by repository
		repositoryUID string

		// Mutex for lock multi access on Pingdom
		sync.Mutex
		// Used for caching result of pingdom (to avoid bursting query limit)
		store           cache.Store
		cacheExpiration int
	}
)

const (
	PingdomChecksTagsByIDStoreKeyPrefix            = "monitoror.pingdom.checksTagsById.store"
	PingdomChecksStoreKeyPrefix                    = "monitoror.pingdom.checks.store"
	PingdomCheckStoreKeyPrefix                     = "monitoror.pingdom.check.store"
	PingdomTransactionChecksTagsByIDStoreKeyPrefix = "monitoror.pingdom.transactionChecksTagsById.store"
	PingdomTransactionChecksStoreKeyPrefix         = "monitoror.pingdom.transactionChecks.store"
	PingdomTransactionCheckStoreKeyPrefix          = "monitoror.pingdom.transactionCheck.store"

	UpCheckStatus     = "up"
	DownCheckStatus   = "down"
	PausedCheckStatus = "paused"

	SuccessfulTransactionCheckStatus = "successful"
	FailingTransactionCheckStatus    = "failing"
	UnknownTransactionCheckStatus    = "unknown"
)

func NewPingdomUsecase(repository api.Repository, store cache.Store, cacheExpiration int) api.Usecase {
	return &pingdomUsecase{
		repository:      repository,
		repositoryUID:   uuid.NewV4().String(),
		store:           store,
		cacheExpiration: cacheExpiration,
	}
}

func (pu *pingdomUsecase) Check(params *models.CheckParams) (*coreModels.Tile, error) {
	return pu.check(false, *params.ID)
}

func (pu *pingdomUsecase) TransactionCheck(params *models.TransactionCheckParams) (*coreModels.Tile, error) {
	return pu.check(true, *params.ID)
}

func (pu *pingdomUsecase) check(transaction bool, checkID int) (*coreModels.Tile, error) {
	tile := coreModels.NewTile(api.PingdomCheckTileType)

	var result models.Check
	var tags string

	// Lookup in store for bulk query for this ID, if found, use it
	if err := pu.store.Get(pu.getTagsByIDStoreKey(transaction, checkID), &tags); err == nil {
		checks, err := pu.loadChecks(transaction, tags)
		if err != nil {
			return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: "unable to find checks"}
		}

		// Find check in array
		for _, tmpCheck := range checks {
			if tmpCheck.ID == checkID {
				result = tmpCheck
				break
			}
		}
	} else // Bulk not found, request single check
	{
		check, err := pu.loadCheck(transaction, checkID)
		if err != nil {
			return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: "unable to find check"}
		}
		result = *check
	}

	// Parse result to tile
	tile.Label = result.Name
	tile.Status = parseCheckStatus(result.Status)

	return tile, nil
}

func (pu *pingdomUsecase) CheckGenerator(params interface{}) ([]uiConfigModels.GeneratedTile, error) {
	cParams := params.(*models.CheckGeneratorParams)
	return pu.checkGenerator(false, cParams.Tags, cParams.SortBy)
}

func (pu *pingdomUsecase) TransactionCheckGenerator(params interface{}) ([]uiConfigModels.GeneratedTile, error) {
	cParams := params.(*models.TransactionCheckGeneratorParams)
	return pu.checkGenerator(true, cParams.Tags, cParams.SortBy)
}

func (pu *pingdomUsecase) checkGenerator(transaction bool, tags, sortBy string) ([]uiConfigModels.GeneratedTile, error) {
	checks, err := pu.loadChecks(transaction, tags)
	if err != nil {
		return nil, &coreModels.MonitororError{Err: err, Message: "unable to list checks"}
	}

	if sortBy == "name" {
		sort.SliceStable(checks, func(i, j int) bool { return checks[i].Name < checks[j].Name })
	}

	var results []uiConfigModels.GeneratedTile
	for _, check := range checks {
		// Adding id -> tags in the store for one minute. This value will be refresh each time we call this route
		// This store will be use to find the best route to call for loading check result.
		_ = pu.store.Set(pu.getTagsByIDStoreKey(transaction, check.ID), tags, time.Minute)

		if check.Status == PausedCheckStatus || check.Status == UnknownTransactionCheckStatus {
			continue
		}

		// Build results
		var p interface{}
		if transaction {
			p = models.TransactionCheckParams{
				ID: pointer.ToInt(check.ID),
			}
		} else {
			p = models.CheckParams{
				ID: pointer.ToInt(check.ID),
			}
		}

		results = append(results, uiConfigModels.GeneratedTile{
			Label:  check.Name,
			Params: p,
		})
	}

	return results, err
}

func (pu *pingdomUsecase) loadCheck(transaction bool, id int) (result *models.Check, err error) {
	// Synchronize to avoid multi call on pingdom api
	pu.Lock()
	defer pu.Unlock()

	// Lookup in cache
	result = &models.Check{}
	key := pu.getCheckStoreKey(transaction, id)
	if err = pu.store.Get(key, result); err == nil {
		// Cache found, return
		return
	}

	if transaction {
		result, err = pu.repository.GetTransactionCheck(id)
	} else {
		result, err = pu.repository.GetCheck(id)
	}
	if err != nil {
		return
	}

	// Adding result in store
	_ = pu.store.Set(key, *result, time.Millisecond*time.Duration(pu.cacheExpiration))

	return
}

func (pu *pingdomUsecase) loadChecks(transaction bool, tags string) (results []models.Check, err error) {
	// Synchronize to avoid multi call on pingdom api
	pu.Lock()
	defer pu.Unlock()

	// Lookup in cache
	key := pu.getChecksStoreKey(transaction, tags)
	if err = pu.store.Get(key, &results); err == nil {
		// Cache found, return
		return
	}

	if transaction {
		results, err = pu.repository.GetTransactionChecks(tags)
	} else {
		results, err = pu.repository.GetChecks(tags)
	}
	if err != nil {
		return
	}

	// Adding result in store
	_ = pu.store.Set(key, results, time.Millisecond*time.Duration(pu.cacheExpiration))

	return
}

func (pu *pingdomUsecase) getTagsByIDStoreKey(transaction bool, id int) string {
	prefix := PingdomChecksTagsByIDStoreKeyPrefix
	if transaction {
		prefix = PingdomTransactionChecksTagsByIDStoreKeyPrefix
	}
	return fmt.Sprintf("%s:%s-%d", prefix, pu.repositoryUID, id)
}

func (pu *pingdomUsecase) getChecksStoreKey(transaction bool, tags string) string {
	prefix := PingdomChecksStoreKeyPrefix
	if transaction {
		prefix = PingdomTransactionChecksStoreKeyPrefix
	}
	return fmt.Sprintf("%s:%s-%s", prefix, pu.repositoryUID, tags)
}

func (pu *pingdomUsecase) getCheckStoreKey(transaction bool, id int) string {
	prefix := PingdomCheckStoreKeyPrefix
	if transaction {
		prefix = PingdomTransactionCheckStoreKeyPrefix
	}
	return fmt.Sprintf("%s:%s-%d", prefix, pu.repositoryUID, id)
}

func parseCheckStatus(status string) coreModels.TileStatus {
	switch status {
	case UpCheckStatus, SuccessfulTransactionCheckStatus:
		return coreModels.SuccessStatus
	case DownCheckStatus, FailingTransactionCheckStatus:
		return coreModels.FailedStatus
	case PausedCheckStatus, UnknownTransactionCheckStatus:
		return coreModels.DisabledStatus
	default:
		return coreModels.UnknownStatus
	}
}
