//+build faker

package github

import (
	uiConfig "github.com/monitoror/monitoror/api/config"
	"github.com/monitoror/monitoror/api/config/versions"
	coreConfig "github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/internal/pkg/monitorable"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/github/api"
	githubDelivery "github.com/monitoror/monitoror/monitorables/github/api/delivery/http"
	githubModels "github.com/monitoror/monitoror/monitorables/github/api/models"
	githubUsecase "github.com/monitoror/monitoror/monitorables/github/api/usecase"
	"github.com/monitoror/monitoror/service/store"
)

type Monitorable struct {
	monitorable.DefaultMonitorableFaker

	store *store.Store

	// Config tile settings
	countTypeSetting  uiConfig.TileEnabler
	checksTypeSetting uiConfig.TileEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store

	// Register Monitorable Tile in config manager
	m.countTypeSetting = store.TileSettingManager.Register(api.GithubCountTileType, versions.MinimalVersion, m.GetVariantNames())
	m.checksTypeSetting = store.TileSettingManager.Register(api.GithubChecksTileType, versions.MinimalVersion, m.GetVariantNames())

	return m
}

func (m *Monitorable) GetDisplayName() string {
	return "GitHub (faker)"
}

func (m *Monitorable) Enable(variantName coreModels.VariantName) {
	usecase := githubUsecase.NewGithubUsecase()
	delivery := githubDelivery.NewGithubDelivery(usecase)

	// EnableTile route to echo
	routeGroup := m.store.MonitorableRouter.Group("/github", variantName)
	routeCount := routeGroup.GET("/count", delivery.GetCount)
	routeChecks := routeGroup.GET("/checks", delivery.GetChecks)

	// EnableTile data for config hydration
	m.countTypeSetting.Enable(variant, &githubModels.CountParams{}, routeCount.Path, coreConfig.DefaultInitialMaxDelay)
	m.checksTypeSetting.Enable(variant, &githubModels.ChecksParams{}, routeChecks.Path, coreConfig.DefaultInitialMaxDelay)
}
