//+build !faker

package pingdom

import (
	"net/url"

	uiConfig "github.com/monitoror/monitoror/api/config/usecase"
	coreConfig "github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/monitorables/pingdom/api"
	pingdomDelivery "github.com/monitoror/monitoror/monitorables/pingdom/api/delivery/http"
	pingdomModels "github.com/monitoror/monitoror/monitorables/pingdom/api/models"
	pingdomRepository "github.com/monitoror/monitoror/monitorables/pingdom/api/repository"
	pingdomUsecase "github.com/monitoror/monitoror/monitorables/pingdom/api/usecase"
	pingdomConfig "github.com/monitoror/monitoror/monitorables/pingdom/config"
	"github.com/monitoror/monitoror/service/store"
)

type Monitorable struct {
	store *store.Store

	config map[string]*pingdomConfig.Pingdom
}

func NewMonitorable(store *store.Store) *Monitorable {
	monitorable := &Monitorable{}
	monitorable.store = store
	monitorable.config = make(map[string]*pingdomConfig.Pingdom)

	// Load core config from env
	coreConfig.LoadMonitorableConfig(&monitorable.config, pingdomConfig.Default)

	// Register Monitorable Tile in config manager
	store.UIConfigManager.RegisterTile(api.PingdomCheckTileType, monitorable.GetVariants(), uiConfig.MinimalVersion)
	store.UIConfigManager.RegisterTile(api.PingdomChecksTileType, monitorable.GetVariants(), uiConfig.MinimalVersion)

	return monitorable
}

func (m *Monitorable) GetVariants() []string { return coreConfig.GetVariantsFromConfig(m.config) }
func (m *Monitorable) IsValid(variant string) bool {
	conf := m.config[variant]

	// Pingdom url can be empty, plugin will use default value
	if conf.URL != "" {
		if _, err := url.Parse(conf.URL); err != nil {
			return false
		}
	}

	return conf.Token != ""
}

func (m *Monitorable) Enable(variant string) {
	conf := m.config[variant]

	repository := pingdomRepository.NewPingdomRepository(conf)
	usecase := pingdomUsecase.NewPingdomUsecase(repository, m.store.CacheStore, conf.CacheExpiration)
	delivery := pingdomDelivery.NewPingdomDelivery(usecase)

	// EnableTile route to echo
	route := m.store.MonitorableRouter.Group("/pingdom", variant).GET("/pingdom", delivery.GetCheck)

	// EnableTile data for config hydration
	m.store.UIConfigManager.EnableTile(api.PingdomCheckTileType, variant, &pingdomModels.ChecksParams{}, route.Path, conf.InitialMaxDelay)
	m.store.UIConfigManager.EnableDynamicTile(api.PingdomChecksTileType, variant, &pingdomModels.ChecksParams{}, usecase.Checks)
}
