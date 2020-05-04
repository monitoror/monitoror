package service

import (
	"github.com/jsdidierlaurent/echo-middleware/cache"
	configDelivery "github.com/monitoror/monitoror/api/config/delivery/http"
	configRepository "github.com/monitoror/monitoror/api/config/repository"
	configUsecase "github.com/monitoror/monitoror/api/config/usecase"
	"github.com/monitoror/monitoror/api/info"
	"github.com/monitoror/monitoror/monitorables"
	"github.com/monitoror/monitoror/service/router"
)

func InitApis(s *Server) {
	// API group
	apiGroup := s.Group("/api/v1")

	// ------------- INFO ------------- //
	infoDelivery := info.NewHTTPInfoDelivery()
	apiGroup.GET("/info", s.store.CacheMiddleware.UpstreamCacheHandlerWithExpiration(cache.NEVER, infoDelivery.GetInfo))

	// ------------- CONFIG ------------- //
	confRepository := configRepository.NewConfigRepository()
	confUsecase := configUsecase.NewConfigUsecase(confRepository, s.store)
	confDelivery := configDelivery.NewConfigDelivery(confUsecase)
	apiGroup.GET("/configs", s.store.CacheMiddleware.UpstreamCacheHandler(confDelivery.GetConfigList))
	apiGroup.GET("/configs/:config", s.store.CacheMiddleware.UpstreamCacheHandler(confDelivery.GetConfig))

	// ---------------------------------- //
	s.store.MonitorableRouter = router.NewMonitorableRouter(apiGroup, s.store.CacheMiddleware)
	// ---------------------------------- //

	// ------------- MONITORABLES ------------- //
	monitorableManager := monitorables.NewMonitorableManager(s.store)
	monitorableManager.RegisterMonitorables()
	monitorableManager.EnableMonitorables()
}
