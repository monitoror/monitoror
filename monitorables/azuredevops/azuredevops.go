//+build !faker

package azuredevops

import (
	"net/url"

	uiConfig "github.com/monitoror/monitoror/api/config/usecase"
	coreConfig "github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/monitorables/azuredevops/api"
	azuredevopsDelivery "github.com/monitoror/monitoror/monitorables/azuredevops/api/delivery/http"
	azuredevopsModels "github.com/monitoror/monitoror/monitorables/azuredevops/api/models"
	azuredevopsRepository "github.com/monitoror/monitoror/monitorables/azuredevops/api/repository"
	azuredevopsUsecase "github.com/monitoror/monitoror/monitorables/azuredevops/api/usecase"
	azuredevopsCoreConfig "github.com/monitoror/monitoror/monitorables/azuredevops/config"
	"github.com/monitoror/monitoror/service/store"
)

type Monitorable struct {
	store *store.Store

	config map[string]*azuredevopsCoreConfig.AzureDevOps
}

func NewMonitorable(store *store.Store) *Monitorable {
	monitorable := &Monitorable{}
	monitorable.store = store
	monitorable.config = make(map[string]*azuredevopsCoreConfig.AzureDevOps)

	// Load core config from env
	coreConfig.LoadMonitorableConfig(&monitorable.config, azuredevopsCoreConfig.Default)

	// Register Monitorable Tile in config manager
	store.UIConfigManager.RegisterTile(api.AzureDevOpsBuildTileType, monitorable.GetVariants(), uiConfig.MinimalVersion)
	store.UIConfigManager.RegisterTile(api.AzureDevOpsReleaseTileType, monitorable.GetVariants(), uiConfig.MinimalVersion)

	return monitorable
}

func (m *Monitorable) GetVariants() []string { return coreConfig.GetVariantsFromConfig(m.config) }
func (m *Monitorable) IsValid(variant string) bool {
	conf := m.config[variant]

	if conf.URL == "" {
		return false
	}

	if _, err := url.Parse(conf.URL); err != nil {
		return false
	}

	return conf.Token != ""
}

func (m *Monitorable) Enable(variant string) {
	conf := m.config[variant]

	repository := azuredevopsRepository.NewAzureDevOpsRepository(conf)
	usecase := azuredevopsUsecase.NewAzureDevOpsUsecase(repository)
	delivery := azuredevopsDelivery.NewAzureDevOpsDelivery(usecase)

	// EnableTile route to echo
	azureGroup := m.store.MonitorableRouter.Group("/azuredevops", variant)
	routeBuild := azureGroup.GET("/build", delivery.GetBuild)
	routeRelease := azureGroup.GET("/release", delivery.GetRelease)

	// EnableTile data for config hydration
	m.store.UIConfigManager.EnableTile(api.AzureDevOpsBuildTileType, variant, &azuredevopsModels.BuildParams{}, routeBuild.Path, conf.InitialMaxDelay)
	m.store.UIConfigManager.EnableTile(api.AzureDevOpsReleaseTileType, variant, &azuredevopsModels.ReleaseParams{}, routeRelease.Path, conf.InitialMaxDelay)
}
