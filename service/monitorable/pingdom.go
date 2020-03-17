//+build !faker

package monitorable

import (
	"net/url"

	"github.com/jsdidierlaurent/echo-middleware/cache"

	"github.com/monitoror/monitoror/config"
	monitorableConfig "github.com/monitoror/monitoror/monitorable/config"
	"github.com/monitoror/monitoror/monitorable/pingdom"
	pingdomDelivery "github.com/monitoror/monitoror/monitorable/pingdom/delivery/http"
	pingdomModels "github.com/monitoror/monitoror/monitorable/pingdom/models"
	pingdomRepository "github.com/monitoror/monitoror/monitorable/pingdom/repository"
	pingdomUsecase "github.com/monitoror/monitoror/monitorable/pingdom/usecase"
	"github.com/monitoror/monitoror/service/router"
)

type pingdomMonitorable struct {
	config map[string]*config.Pingdom

	// Used for caching result of pingdom (to avoid bursting query limit)
	store cache.Store
}

func NewPingdomMonitorable(config map[string]*config.Pingdom, store cache.Store) Monitorable {
	return &pingdomMonitorable{config: config, store: store}
}

func (m *pingdomMonitorable) GetHelp() string       { return "HEEEELLLPPPP" }
func (m *pingdomMonitorable) GetVariants() []string { return config.GetVariantsFromConfig(m.config) }
func (m *pingdomMonitorable) isEnabled(variant string) bool {
	conf := m.config[variant]

	// Pingdom url can be empty, plugin will use default value
	if conf.URL != "" {
		if _, err := url.Parse(conf.URL); err != nil {
			return false
		}
	}

	return conf.Token != ""
}

func (m *pingdomMonitorable) Register(variant string, router router.MonitorableRouter, configManager monitorableConfig.Manager) bool {
	enabled := m.isEnabled(variant)
	if enabled {
		conf := m.config[variant]

		repository := pingdomRepository.NewPingdomRepository(conf)
		usecase := pingdomUsecase.NewPingdomUsecase(repository, conf, m.store)
		delivery := pingdomDelivery.NewPingdomDelivery(usecase)

		// RegisterTile route to echo
		route := router.Group("/pingdom", variant).GET("/check", delivery.GetCheck)

		// RegisterTile data for config hydration
		configManager.RegisterTile(pingdom.PingdomCheckTileType, variant, &pingdomModels.CheckParams{}, route.Path, conf.InitialMaxDelay)
		configManager.RegisterDynamicTile(pingdom.PingdomChecksTileType, variant, &pingdomModels.ChecksParams{}, usecase.Checks)
	} else {
		// RegisterTile data for config verify
		configManager.DisableTile(pingdom.PingdomCheckTileType, variant)
		configManager.DisableTile(pingdom.PingdomChecksTileType, variant)
	}

	return enabled
}
