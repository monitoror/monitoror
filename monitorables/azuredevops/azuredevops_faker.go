//+build faker

package azuredevops

import (
	uiConfig "github.com/monitoror/monitoror/api/config"
	"github.com/monitoror/monitoror/api/config/versions"
	coreConfig "github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/internal/pkg/monitorable"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/azuredevops/api"
	azuredevopsDelivery "github.com/monitoror/monitoror/monitorables/azuredevops/api/delivery/http"
	azuredevopsModels "github.com/monitoror/monitoror/monitorables/azuredevops/api/models"
	azuredevopsUsecase "github.com/monitoror/monitoror/monitorables/azuredevops/api/usecase"
	"github.com/monitoror/monitoror/service/store"
)

type Monitorable struct {
	monitorable.DefaultMonitorableFaker
	store *store.Store

	// Config tile settings
	buildTileSetting   uiConfig.TileEnabler
	releaseTileSetting uiConfig.TileEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store

	// Register Monitorable Tile in config manager
	m.buildTileSetting = store.TileSettingManager.Register(api.AzureDevOpsBuildTileType, versions.MinimalVersion, m.GetVariantNames())
	m.releaseTileSetting = store.TileSettingManager.Register(api.AzureDevOpsReleaseTileType, versions.MinimalVersion, m.GetVariantNames())

	return m
}

func (m *Monitorable) GetDisplayName() string {
	return "Azure DevOps (faker)"
}

func (m *Monitorable) Enable(variantName coreModels.VariantName) {
	usecase := azuredevopsUsecase.NewAzureDevOpsUsecase()
	delivery := azuredevopsDelivery.NewAzureDevOpsDelivery(usecase)

	// EnableTile route to echo
	routeGroup := m.store.MonitorableRouter.Group("/azuredevops", variantName)
	routeBuild := routeGroup.GET("/build", delivery.GetBuild)
	routeRelease := routeGroup.GET("/release", delivery.GetRelease)

	// EnableTile data for config hydration
	m.buildTileSetting.Enable(variant, &azuredevopsModels.BuildParams{}, routeBuild.Path, coreConfig.DefaultInitialMaxDelay)
	m.releaseTileSetting.Enable(variant, &azuredevopsModels.ReleaseParams{}, routeRelease.Path, coreConfig.DefaultInitialMaxDelay)
}
