//go:build faker

package youtrack

import (
	"github.com/monitoror/monitoror/api/config/versions"
	"github.com/monitoror/monitoror/internal/pkg/monitorable"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/youtrack/api"
	youtrackDelivery "github.com/monitoror/monitoror/monitorables/youtrack/api/delivery/http"
	youtrackModels "github.com/monitoror/monitoror/monitorables/youtrack/api/models"
	youtrackUsecase "github.com/monitoror/monitoror/monitorables/youtrack/api/usecase"
	"github.com/monitoror/monitoror/registry"
	"github.com/monitoror/monitoror/store"
)

type Monitorable struct {
	monitorable.DefaultMonitorableFaker

	store *store.Store

	// Config tile settings
	countIssuesTileEnabler registry.TileEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store

	// Register Monitorable Tile in config manager
	m.countIssuesTileEnabler = store.Registry.RegisterTile(api.YoutrackCountIssuesTileType, versions.Version2002, m.GetVariantsNames())

	return m
}

func (m *Monitorable) GetDisplayName() string {
	return "YouTrack"
}

func (m *Monitorable) Enable(variantName coreModels.VariantName) {
	usecase := youtrackUsecase.NewYoutrackUsecase()
	delivery := youtrackDelivery.NewYoutrackDelivery(usecase)

	// EnableTile route to echo
	routeGroup := m.store.MonitorableRouter.Group("/youtrack", variantName)
	routeCountIssues := routeGroup.GET("/count-issues", delivery.GetCountIssues)

	// EnableTile data for config hydration
	m.countIssuesTileEnabler.Enable(variantName, &youtrackModels.IssuesCountParams{}, routeCountIssues.Path)
}
