//+build !faker

package pingdom

import (
	"fmt"
	"net/url"

	coreModels "github.com/monitoror/monitoror/models"

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

	config map[coreModels.Variant]*pingdomConfig.Pingdom
}

func NewMonitorable(store *store.Store) *Monitorable {
	monitorable := &Monitorable{}
	monitorable.store = store
	monitorable.config = make(map[coreModels.Variant]*pingdomConfig.Pingdom)

	// Load core config from env
	coreConfig.LoadMonitorableConfig(&monitorable.config, pingdomConfig.Default)

	// Register Monitorable Tile in config manager
	store.UIConfigManager.RegisterTile(api.PingdomCheckTileType, monitorable.GetVariants(), uiConfig.MinimalVersion)
	store.UIConfigManager.RegisterTile(api.PingdomChecksTileType, monitorable.GetVariants(), uiConfig.MinimalVersion)

	return monitorable
}

func (m *Monitorable) GetDisplayName() string {
	return "Pingdom"
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

	repository := pingdomRepository.NewPingdomRepository(conf)
	usecase := pingdomUsecase.NewPingdomUsecase(repository, m.store.CacheStore, conf.CacheExpiration)
	delivery := pingdomDelivery.NewPingdomDelivery(usecase)

	// EnableTile route to echo
	route := m.store.MonitorableRouter.Group("/pingdom", variant).GET("/pingdom", delivery.GetCheck)

	// EnableTile data for config hydration
	m.store.UIConfigManager.EnableTile(api.PingdomCheckTileType, variant,
		&pingdomModels.ChecksParams{}, route.Path, conf.InitialMaxDelay)
	m.store.UIConfigManager.EnableDynamicTile(api.PingdomChecksTileType, variant,
		&pingdomModels.ChecksParams{}, usecase.Checks)
}
