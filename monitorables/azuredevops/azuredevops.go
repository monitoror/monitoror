//+build !faker

package azuredevops

import (
	"fmt"
	"net/url"

	uiConfig "github.com/monitoror/monitoror/api/config/usecase"
	coreConfig "github.com/monitoror/monitoror/config"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/azuredevops/api"
	azuredevopsDelivery "github.com/monitoror/monitoror/monitorables/azuredevops/api/delivery/http"
	azuredevopsModels "github.com/monitoror/monitoror/monitorables/azuredevops/api/models"
	azuredevopsRepository "github.com/monitoror/monitoror/monitorables/azuredevops/api/repository"
	azuredevopsUsecase "github.com/monitoror/monitoror/monitorables/azuredevops/api/usecase"
	azuredevopsConfig "github.com/monitoror/monitoror/monitorables/azuredevops/config"
	"github.com/monitoror/monitoror/service/store"
)

type Monitorable struct {
	store *store.Store

	config map[coreModels.Variant]*azuredevopsConfig.AzureDevOps
}

func NewMonitorable(store *store.Store) *Monitorable {
	monitorable := &Monitorable{}
	monitorable.store = store
	monitorable.config = make(map[coreModels.Variant]*azuredevopsConfig.AzureDevOps)

	// Load core config from env
	coreConfig.LoadMonitorableConfig(&monitorable.config, azuredevopsConfig.Default)

	// Register Monitorable Tile in config manager
	store.UIConfigManager.RegisterTile(api.AzureDevOpsBuildTileType, monitorable.GetVariants(), uiConfig.MinimalVersion)
	store.UIConfigManager.RegisterTile(api.AzureDevOpsReleaseTileType, monitorable.GetVariants(), uiConfig.MinimalVersion)

	return monitorable
}

func (m *Monitorable) GetDisplayName() string {
	return "Azure DevOps"
}

func (m *Monitorable) GetVariants() []coreModels.Variant {
	return coreConfig.GetVariantsFromConfig(m.config)
}

func (m *Monitorable) Validate(variant coreModels.Variant) (bool, error) {
	conf := m.config[variant]

	// No configuration set
	if conf.URL == "" && conf.Token == "" {
		return false, nil
	}

	// Error in URL
	if _, err := url.Parse(conf.URL); err != nil {
		return false, fmt.Errorf(`%s contains invalid URL: "%s"`, coreConfig.GetEnvFromMonitorableVariable(conf, variant, "URL"), conf.URL)
	}

	// Error in Token
	if conf.Token == "" {
		return false, fmt.Errorf(`%s is required, no value founds`, coreConfig.GetEnvFromMonitorableVariable(conf, variant, "TOKEN"))
	}

	return true, nil
}

func (m *Monitorable) Enable(variant coreModels.Variant) {
	conf := m.config[variant]

	repository := azuredevopsRepository.NewAzureDevOpsRepository(conf)
	usecase := azuredevopsUsecase.NewAzureDevOpsUsecase(repository)
	delivery := azuredevopsDelivery.NewAzureDevOpsDelivery(usecase)

	// EnableTile route to echo
	azureGroup := m.store.MonitorableRouter.Group("/azuredevops", variant)
	routeBuild := azureGroup.GET("/build", delivery.GetBuild)
	routeRelease := azureGroup.GET("/release", delivery.GetRelease)

	// EnableTile data for config hydration
	m.store.UIConfigManager.EnableTile(api.AzureDevOpsBuildTileType, variant,
		&azuredevopsModels.BuildParams{}, routeBuild.Path, conf.InitialMaxDelay)
	m.store.UIConfigManager.EnableTile(api.AzureDevOpsReleaseTileType, variant,
		&azuredevopsModels.ReleaseParams{}, routeRelease.Path, conf.InitialMaxDelay)
}
