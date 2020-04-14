//+build !faker

package azuredevops

import (
	"fmt"
	"net/url"

	"github.com/monitoror/monitoror/api/config/versions"
	pkgMonitorable "github.com/monitoror/monitoror/internal/pkg/monitorable"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/azuredevops/api"
	azuredevopsDelivery "github.com/monitoror/monitoror/monitorables/azuredevops/api/delivery/http"
	azuredevopsModels "github.com/monitoror/monitoror/monitorables/azuredevops/api/models"
	azuredevopsRepository "github.com/monitoror/monitoror/monitorables/azuredevops/api/repository"
	azuredevopsUsecase "github.com/monitoror/monitoror/monitorables/azuredevops/api/usecase"
	azuredevopsConfig "github.com/monitoror/monitoror/monitorables/azuredevops/config"
	"github.com/monitoror/monitoror/service/registry"
	"github.com/monitoror/monitoror/service/store"
)

type Monitorable struct {
	store *store.Store

	config map[coreModels.VariantName]*azuredevopsConfig.AzureDevOps

	// Config tile settings
	buildTileEnabler   registry.TileEnabler
	releaseTileEnabler registry.TileEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store
	m.config = make(map[coreModels.VariantName]*azuredevopsConfig.AzureDevOps)

	// Load core config from env
	pkgMonitorable.LoadConfig(&m.config, azuredevopsConfig.Default)

	// Register Monitorable Tile in config manager
	m.buildTileEnabler = store.Registry.RegisterTile(api.AzureDevOpsBuildTileType, versions.MinimalVersion, m.GetVariantNames())
	m.releaseTileEnabler = store.Registry.RegisterTile(api.AzureDevOpsReleaseTileType, versions.MinimalVersion, m.GetVariantNames())

	return m
}

func (m *Monitorable) GetDisplayName() string {
	return "Azure DevOps"
}

func (m *Monitorable) GetVariantNames() []coreModels.VariantName {
	return pkgMonitorable.GetVariants(m.config)
}

func (m *Monitorable) Validate(variantName coreModels.VariantName) (bool, error) {
	conf := m.config[variantName]

	// No configuration set
	if conf.URL == "" && conf.Token == "" {
		return false, nil
	}

	// Error in URL
	if _, err := url.Parse(conf.URL); err != nil {
		return false, fmt.Errorf(`%s contains invalid URL: "%s"`, pkgMonitorable.BuildMonitorableEnvKey(conf, variantName, "URL"), conf.URL)
	}

	// Error in Token
	if conf.Token == "" {
		return false, fmt.Errorf(`%s is required, no value found`, pkgMonitorable.BuildMonitorableEnvKey(conf, variantName, "TOKEN"))
	}

	return true, nil
}

func (m *Monitorable) Enable(variantName coreModels.VariantName) {
	conf := m.config[variantName]

	repository := azuredevopsRepository.NewAzureDevOpsRepository(conf)
	usecase := azuredevopsUsecase.NewAzureDevOpsUsecase(repository)
	delivery := azuredevopsDelivery.NewAzureDevOpsDelivery(usecase)

	// EnableTile route to echo
	routeGroup := m.store.MonitorableRouter.Group("/azuredevops", variantName)
	routeBuild := routeGroup.GET("/build", delivery.GetBuild)
	routeRelease := routeGroup.GET("/release", delivery.GetRelease)

	// EnableTile data for config hydration
	m.buildTileEnabler.Enable(variantName, &azuredevopsModels.BuildParams{}, routeBuild.Path)
	m.releaseTileEnabler.Enable(variantName, &azuredevopsModels.ReleaseParams{}, routeRelease.Path)
}
