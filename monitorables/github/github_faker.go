//+build faker

package github

import (
	uiConfig "github.com/monitoror/monitoror/api/config/usecase"
	coreConfig "github.com/monitoror/monitoror/config"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/github/api"
	githubDelivery "github.com/monitoror/monitoror/monitorables/github/api/delivery/http"
	githubModels "github.com/monitoror/monitoror/monitorables/github/api/models"
	githubUsecase "github.com/monitoror/monitoror/monitorables/github/api/usecase"
	"github.com/monitoror/monitoror/service/store"
)

type Monitorable struct {
	store *store.Store
}

func NewMonitorable(store *store.Store) *Monitorable {
	monitorable := &Monitorable{}
	monitorable.store = store

	// Register Monitorable Tile in config manager
	store.UIConfigManager.RegisterTile(api.GithubCountTileType, monitorable.GetVariants(), uiConfig.MinimalVersion)
	store.UIConfigManager.RegisterTile(api.GithubChecksTileType, monitorable.GetVariants(), uiConfig.MinimalVersion)

	return monitorable
}

func (m *Monitorable) GetDisplayName() string                { return "GitHub (faker)" }
func (m *Monitorable) GetVariants() []string                 { return []coreModels.Variant{coreConfig.DefaultVariant} }
func (m *Monitorable) Validate(variant string) (bool, error) { return true, nil }

func (m *Monitorable) Enable(variant string) {
	usecase := githubUsecase.NewGithubUsecase()
	delivery := githubDelivery.NewGithubDelivery(usecase)

	// EnableTile route to echo
	githubGroup := m.store.MonitorableRouter.Group("/github", variant)
	routeCount := githubGroup.GET("/count", delivery.GetCount)
	routeChecks := githubGroup.GET("/checks", delivery.GetChecks)

	// EnableTile data for config hydration
	m.store.UIConfigManager.EnableTile(api.GithubCountTileType, variant,
		&githubModels.CountParams{}, routeCount.Path, coreConfig.DefaultInitialMaxDelay)
	m.store.UIConfigManager.EnableTile(api.GithubChecksTileType, variant,
		&githubModels.ChecksParams{}, routeChecks.Path, coreConfig.DefaultInitialMaxDelay)
}
