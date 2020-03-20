//+build faker

package jenkins

import (
	uiConfig "github.com/monitoror/monitoror/api/config/usecase"
	coreConfig "github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/monitorables/jenkins/api"
	jenkinsDelivery "github.com/monitoror/monitoror/monitorables/jenkins/api/delivery/http"
	jenkinsModels "github.com/monitoror/monitoror/monitorables/jenkins/api/models"
	jenkinsUsecase "github.com/monitoror/monitoror/monitorables/jenkins/api/usecase"
	"github.com/monitoror/monitoror/service/store"
)

type Monitorable struct {
	store *store.Store
}

func NewMonitorable(store *store.Store) *Monitorable {
	monitorable := &Monitorable{}
	monitorable.store = store

	// Register Monitorable Tile in config manager
	store.UIConfigManager.RegisterTile(api.JenkinsBuildTileType, monitorable.GetVariants(), uiConfig.MinimalVersion)

	return monitorable
}

func (m *Monitorable) GetVariants() []string       { return []string{coreConfig.DefaultVariant} }
func (m *Monitorable) IsValid(variant string) bool { return true }

func (m *Monitorable) Enable(variant string) {
	usecase := jenkinsUsecase.NewJenkinsUsecase()
	delivery := jenkinsDelivery.NewJenkinsDelivery(usecase)

	// EnableTile route to echo
	jenkinsGroup := m.store.MonitorableRouter.Group("/jenkins", variant)
	route := jenkinsGroup.GET("/build", delivery.GetBuild)

	// EnableTile data for config hydration
	m.store.UIConfigManager.EnableTile(api.JenkinsBuildTileType, variant, &jenkinsModels.BuildParams{}, route.Path, coreConfig.DefaultInitialMaxDelay)
}
