//+build faker

package monitorable

import (
	"github.com/jsdidierlaurent/echo-middleware/cache"
	"github.com/monitoror/monitoror/config"
	monitorableConfig "github.com/monitoror/monitoror/monitorable/config"
	"github.com/monitoror/monitoror/monitorable/http"
	httpDelivery "github.com/monitoror/monitoror/monitorable/http/delivery/http"
	httpModels "github.com/monitoror/monitoror/monitorable/http/models"
	httpUsecase "github.com/monitoror/monitoror/monitorable/http/usecase"
	"github.com/monitoror/monitoror/service/router"
)

type httpMonitorable struct{}

func NewHTTPMonitorable(_ map[string]*config.HTTP, _ cache.Store, _ int) Monitorable {
	return &httpMonitorable{}
}

func (m *httpMonitorable) GetHelp() string       { return "" }
func (m *httpMonitorable) GetVariants() []string { return []string{config.DefaultVariant} }

func (m *httpMonitorable) Register(variant string, router router.MonitorableRouter, configManager monitorableConfig.Manager) bool {
	usecase := httpUsecase.NewHTTPUsecase()
	delivery := httpDelivery.NewHTTPDelivery(usecase)

	// RegisterTile route to echo
	httpGroup := router.Group("/http", variant)
	routeStatus := httpGroup.GET("/status", delivery.GetHTTPStatus)
	routeRaw := httpGroup.GET("/raw", delivery.GetHTTPRaw)
	routeJSON := httpGroup.GET("/formatted", delivery.GetHTTPFormatted)

	// RegisterTile data for config hydration
	configManager.RegisterTile(http.HTTPStatusTileType, variant, &httpModels.HTTPStatusParams{}, routeStatus.Path, config.DefaultInitialMaxDelay)
	configManager.RegisterTile(http.HTTPRawTileType, variant, &httpModels.HTTPRawParams{}, routeRaw.Path, config.DefaultInitialMaxDelay)
	configManager.RegisterTile(http.HTTPFormattedTileType, variant, &httpModels.HTTPFormattedParams{}, routeJSON.Path, config.DefaultInitialMaxDelay)

	return true
}
