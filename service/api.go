package service

import (
	"github.com/jsdidierlaurent/echo-middleware/cache"
	"github.com/monitoror/monitoror/handlers"
	"github.com/monitoror/monitoror/monitorable/config"
	configDelivery "github.com/monitoror/monitoror/monitorable/config/delivery/http"
	configRepository "github.com/monitoror/monitoror/monitorable/config/repository"
	configUsecase "github.com/monitoror/monitoror/monitorable/config/usecase"
	"github.com/monitoror/monitoror/service/monitorable"
	"github.com/monitoror/monitoror/service/router"
)

type monitorableManager struct {
	server *Server

	apiRouter     router.MonitorableRouter
	configManager config.Manager
}

func InitApis(s *Server) {
	// API group
	apiGroup := s.Group("/api/v1")

	// ------------- INFO ------------- //
	infoDelivery := handlers.NewHTTPInfoDelivery()
	apiGroup.GET("/info", s.cm.UpstreamCacheHandlerWithExpiration(cache.NEVER, infoDelivery.GetInfo))

	// ------------- CONFIG ------------- //
	confRepository := configRepository.NewConfigRepository()
	confUsecase := configUsecase.NewConfigUsecase(confRepository, s.store, s.config.DownstreamCacheExpiration)
	confDelivery := configDelivery.NewConfigDelivery(confUsecase)

	apiGroup.GET("/config", s.cm.UpstreamCacheHandler(confDelivery.GetConfig))

	// ---------------------------------- //
	m := newMonitorableManager(s, router.NewMonitorableRouter(apiGroup, s.cm), confUsecase)

	// ------------- AZURE DEVOPS ------------- //
	m.registerMonitorable(monitorable.NewAzureDevOpsMonitorable(s.config.Monitorable.AzureDevOps))
	// ------------- GITHUB ------------- //
	m.registerMonitorable(monitorable.NewGithubMonitorable(s.config.Monitorable.Github))
	// ------------- HTTP ------------- //
	m.registerMonitorable(monitorable.NewHTTPMonitorable(s.config.Monitorable.HTTP, s.store, s.config.UpstreamCacheExpiration))
	// ------------- JENKINS ------------- //
	m.registerMonitorable(monitorable.NewJenkinsMonitorable(s.config.Monitorable.Jenkins))
	// ------------- PING ------------- //
	m.registerMonitorable(monitorable.NewPingMonitorable(s.config.Monitorable.Ping))
	// ------------- PINGDOM ------------- //
	m.registerMonitorable(monitorable.NewPingdomMonitorable(s.config.Monitorable.Pingdom, s.store))
	// ------------- PORT ------------- //
	m.registerMonitorable(monitorable.NewPortMonitorable(s.config.Monitorable.Port))
	// ------------- TRAVIS CI ------------- //
	m.registerMonitorable(monitorable.NewTravisCIMonitorable(s.config.Monitorable.TravisCI))
}

func newMonitorableManager(s *Server, r router.MonitorableRouter, m config.Manager) *monitorableManager {
	return &monitorableManager{s, r, m}
}

func (m *monitorableManager) registerMonitorable(monitorable monitorable.Monitorable) {
	// TODO: Store monitorable and execute all RegisterTile after (for improve core startup log)

	for _, variant := range monitorable.GetVariants() {
		monitorable.Register(variant, m.apiRouter, m.configManager)
	}
}
