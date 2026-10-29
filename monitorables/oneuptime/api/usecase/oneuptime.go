//go:build !faker

package usecase

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/oneuptime/api"
	"github.com/monitoror/monitoror/monitorables/oneuptime/api/models"

	"github.com/jsdidierlaurent/echo-middleware/cache"
	uuid "github.com/satori/go.uuid"
)

type (
	oneuptimeUsecase struct {
		repository api.Repository
		// Used to generate store key by repository
		repositoryUID string

		// Mutex for lock multi access on OneUptime
		sync.Mutex
		// Used for caching monitors list, one API call refreshes every tile
		store           cache.Store
		cacheExpiration int
	}
)

const (
	OneuptimeMonitorStoreKeyPrefix  = "monitoror.oneuptime.monitor.store"
	OneuptimeMonitorsStoreKeyPrefix = "monitoror.oneuptime.monitors.store"

	OneuptimeLabelsSeparator = ","
)

func NewOneuptimeUsecase(repository api.Repository, store cache.Store, cacheExpiration int) api.Usecase {
	return &oneuptimeUsecase{
		repository:      repository,
		repositoryUID:   uuid.NewV4().String(),
		store:           store,
		cacheExpiration: cacheExpiration,
	}
}

func (ou *oneuptimeUsecase) Monitor(params *models.MonitorParams) (*coreModels.Tile, error) {
	tile := coreModels.NewTile(api.OneuptimeMonitorTileType)

	var result models.Monitor
	if err := ou.store.Get(ou.getMonitorStoreKey(params.ID), &result); err != nil {
		if _, err := ou.loadMonitors(); err != nil {
			return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: "unable to load monitors"}
		}

		if err := ou.store.Get(ou.getMonitorStoreKey(params.ID), &result); err != nil {
			return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: "unable to find monitor"}
		}
	}

	tile.Label = result.Name
	tile.Status = parseMonitorStatus(&result)

	return tile, nil
}

func (ou *oneuptimeUsecase) MonitorGenerator(params interface{}) ([]uiConfigModels.GeneratedTile, error) {
	gParams := params.(*models.MonitorGeneratorParams)

	monitors, err := ou.loadMonitors()
	if err != nil {
		return nil, &coreModels.MonitororError{Err: err, Message: "unable to list monitors"}
	}

	if gParams.SortBy == "name" {
		sort.SliceStable(monitors, func(i, j int) bool { return monitors[i].Name < monitors[j].Name })
	}

	var results []uiConfigModels.GeneratedTile
	for _, monitor := range monitors {
		if monitor.IsDisabled {
			continue
		}

		if gParams.Labels != "" && !monitor.MatchOneLabel(strings.Split(gParams.Labels, OneuptimeLabelsSeparator)) {
			continue
		}

		results = append(results, uiConfigModels.GeneratedTile{
			Label:  monitor.Name,
			Params: models.MonitorParams{ID: monitor.ID},
		})
	}

	return results, nil
}

func (ou *oneuptimeUsecase) loadMonitors() (results []models.Monitor, err error) {
	// Synchronize to avoid multi call on OneUptime API
	ou.Lock()
	defer ou.Unlock()

	key := ou.getMonitorsStoreKey()
	if err = ou.store.Get(key, &results); err == nil {
		return
	}

	results, err = ou.repository.GetMonitors()
	if err != nil {
		return
	}

	expiration := time.Millisecond * time.Duration(ou.cacheExpiration)
	_ = ou.store.Set(key, results, expiration)
	for _, monitor := range results {
		_ = ou.store.Set(ou.getMonitorStoreKey(monitor.ID), monitor, expiration)
	}

	return
}

func (ou *oneuptimeUsecase) getMonitorStoreKey(id string) string {
	return fmt.Sprintf("%s:%s-%s", OneuptimeMonitorStoreKeyPrefix, ou.repositoryUID, id)
}

func (ou *oneuptimeUsecase) getMonitorsStoreKey() string {
	return fmt.Sprintf("%s:%s", OneuptimeMonitorsStoreKeyPrefix, ou.repositoryUID)
}

func parseMonitorStatus(monitor *models.Monitor) coreModels.TileStatus {
	switch {
	case monitor.IsDisabled:
		return coreModels.DisabledStatus
	case monitor.IsOperational:
		return coreModels.SuccessStatus
	case monitor.IsOffline:
		return coreModels.FailedStatus
	default:
		// Degraded or any custom status of the project
		return coreModels.WarningStatus
	}
}
