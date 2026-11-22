//go:build faker

package gitea

import (
	"github.com/monitoror/monitoror/api/config/versions"
	"github.com/monitoror/monitoror/internal/pkg/monitorable"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/gitea/api"
	giteaDelivery "github.com/monitoror/monitoror/monitorables/gitea/api/delivery/http"
	giteaModels "github.com/monitoror/monitoror/monitorables/gitea/api/models"
	giteaUsecase "github.com/monitoror/monitoror/monitorables/gitea/api/usecase"
	"github.com/monitoror/monitoror/registry"
	"github.com/monitoror/monitoror/store"
)

type Monitorable struct {
	monitorable.DefaultMonitorableFaker

	store *store.Store

	// Config tile settings
	countTileEnabler       registry.TileEnabler
	checksTileEnabler      registry.TileEnabler
	pullRequestTileEnabler registry.TileEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store

	// Register Monitorable Tile in config manager
	m.countTileEnabler = store.Registry.RegisterTile(api.GiteaCountTileType, versions.Version2001, m.GetVariantsNames())
	m.checksTileEnabler = store.Registry.RegisterTile(api.GiteaChecksTileType, versions.Version2001, m.GetVariantsNames())
	m.pullRequestTileEnabler = store.Registry.RegisterTile(api.GiteaPullRequestTileType, versions.Version2001, m.GetVariantsNames())

	return m
}

func (m *Monitorable) GetDisplayName() string {
	return "Gitea"
}

func (m *Monitorable) Enable(variantName coreModels.VariantName) {
	usecase := giteaUsecase.NewGiteaUsecase()
	delivery := giteaDelivery.NewGiteaDelivery(usecase)

	// EnableTile route to echo
	routeGroup := m.store.MonitorableRouter.Group("/gitea", variantName)
	routeCount := routeGroup.GET("/count", delivery.GetCount)
	routeChecks := routeGroup.GET("/checks", delivery.GetChecks)
	routePullRequest := routeGroup.GET("/pullrequest", delivery.GetPullRequest)

	// EnableTile data for config hydration
	m.countTileEnabler.Enable(variantName, &giteaModels.CountParams{}, routeCount.Path)
	m.checksTileEnabler.Enable(variantName, &giteaModels.ChecksParams{}, routeChecks.Path)
	m.pullRequestTileEnabler.Enable(variantName, &giteaModels.PullRequestParams{}, routePullRequest.Path)
}
