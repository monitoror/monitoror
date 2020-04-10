//+build !faker

package pingdom

import (
	"fmt"
	"net/url"

	uiConfig "github.com/monitoror/monitoror/api/config"
	"github.com/monitoror/monitoror/api/config/versions"
	pkgMonitorable "github.com/monitoror/monitoror/internal/pkg/monitorable"

	coreModels "github.com/monitoror/monitoror/models"

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

	config map[coreModels.VariantName]*pingdomConfig.Pingdom

	// Config tile settings
	checkTileSetting          uiConfig.TileEnabler
	checkGeneratorTileSetting uiConfig.TileGeneratorEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store
	m.config = make(map[coreModels.VariantName]*pingdomConfig.Pingdom)

	// Load core config from env
	pkgMonitorable.LoadConfig(&m.config, pingdomConfig.Default)

	// Register Monitorable Tile in config manager
	m.checkTileSetting = store.TileSettingManager.Register(api.PingdomCheckTileType, versions.MinimalVersion, m.GetVariantNames())
	m.checkGeneratorTileSetting = store.TileSettingManager.RegisterGenerator(api.PingdomCheckTileType, versions.MinimalVersion, m.GetVariantNames())

	return m
}

func (m *Monitorable) GetDisplayName() string {
	return "Pingdom"
}

func (m *Monitorable) GetVariantNames() []coreModels.VariantName {
	return pkgMonitorable.GetVariants(m.config)
}

func (m *Monitorable) Validate(variantName coreModels.VariantName) (bool, error) {
	conf := m.config[variantName]

	// No configuration set
	if conf.URL == pingdomConfig.Default.URL && conf.Token == "" {
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

	repository := pingdomRepository.NewPingdomRepository(conf)
	usecase := pingdomUsecase.NewPingdomUsecase(repository, m.store.CacheStore, conf.CacheExpiration)
	delivery := pingdomDelivery.NewPingdomDelivery(usecase)

	// EnableTile route to echo
	routeGroup := m.store.MonitorableRouter.Group("/pingdom", variantName)
	route := routeGroup.GET("/pingdom", delivery.GetCheck)

	// EnableTile data for config hydration
	m.checkTileSetting.Enable(variant, &pingdomModels.CheckParams{}, route.Path, conf.InitialMaxDelay)
	m.checkGeneratorTileSetting.Enable(variant, &pingdomModels.CheckGeneratorParams{}, usecase.CheckGenerator)
}
