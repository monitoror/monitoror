//+build faker

package azuredevops

import (
	uiConfig "github.com/monitoror/monitoror/api/config/usecase"
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
}

func NewMonitorable(store *store.Store) *Monitorable {
	monitorable := &Monitorable{}
	monitorable.store = store

	// Register Monitorable Tile in config manager
	store.UIConfigManager.RegisterTile(api.AzureDevOpsBuildTileType, monitorable.GetVariants(), uiConfig.MinimalVersion)
	store.UIConfigManager.RegisterTile(api.AzureDevOpsReleaseTileType, monitorable.GetVariants(), uiConfig.MinimalVersion)

	return monitorable
}

func (m *Monitorable) GetDisplayName() string {
	return "Azure DevOps (faker)"
}

func (m *Monitorable) Enable(variant coreModels.VariantName) {
	usecase := azuredevopsUsecase.NewAzureDevOpsUsecase()
	delivery := azuredevopsDelivery.NewAzureDevOpsDelivery(usecase)

	// EnableTile route to echo
	routeGroup := m.store.MonitorableRouter.Group("/azuredevops", variant)
	routeBuild := routeGroup.GET("/build", delivery.GetBuild)
	routeRelease := routeGroup.GET("/release", delivery.GetRelease)

	// EnableTile data for config hydration
	m.store.UIConfigManager.EnableTile(api.AzureDevOpsBuildTileType, variant,
		&azuredevopsModels.BuildParams{}, routeBuild.Path, coreConfig.DefaultInitialMaxDelay)
	m.store.UIConfigManager.EnableTile(api.AzureDevOpsReleaseTileType, variant,
		&azuredevopsModels.ReleaseParams{}, routeRelease.Path, coreConfig.DefaultInitialMaxDelay)
}
