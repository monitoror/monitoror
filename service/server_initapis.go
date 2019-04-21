//+build !faker

package service

import (
	"github.com/jsdidierlaurent/echo-middleware/cache"
	"github.com/monitoror/monitoror/handlers"

	_pingDelivery "github.com/monitoror/monitoror/monitorable/ping/delivery/http"
	_pingRepository "github.com/monitoror/monitoror/monitorable/ping/repository"
	_pingUsecase "github.com/monitoror/monitoror/monitorable/ping/usecase"
	_portDelivery "github.com/monitoror/monitoror/monitorable/port/delivery/http"
	_portRepository "github.com/monitoror/monitoror/monitorable/port/repository"
	_portUsecase "github.com/monitoror/monitoror/monitorable/port/usecase"
	_travisciDelivery "github.com/monitoror/monitoror/monitorable/travisci/delivery/http"
	_travisciRepository "github.com/monitoror/monitoror/monitorable/travisci/repository"
	_travisciUsecase "github.com/monitoror/monitoror/monitorable/travisci/usecase"
)

func (s *Server) initApis() {
	v1 := s.Group("/api/v1")

	// Register route utils function

	// ------------- INFO ------------- //
	register("info", true, func() {
		infoHandler := handlers.HttpInfoHandler(s.config)

		v1.GET("/info", s.cm.UpstreamCacheHandlerWithExpiration(cache.NEVER, infoHandler.GetInfo))
	})

	// ------------- PING ------------- //
	register("ping", true, func() {
		pingRepo := _pingRepository.NewPingRepository(s.config)
		pingUC := _pingUsecase.NewPingUsecase(pingRepo)
		pingHandler := _pingDelivery.NewHttpPingHandler(pingUC)

		v1.GET("/ping", s.cm.UpstreamCacheHandler(pingHandler.GetPing))
	})

	// ------------- PORT ------------- //
	register("port", true, func() {
		portRepo := _portRepository.NewPortRepository(s.config)
		portUC := _portUsecase.NewPortUsecase(portRepo)
		portHandler := _portDelivery.NewHttpPortHandler(portUC)

		v1.GET("/port", s.cm.UpstreamCacheHandler(portHandler.GetPort))
	})

	// ------------- TRAVIS CI ------------- //
	register("travis-ci", s.config.Monitorable.TravisCI.Url != "", func() {
		travisciRepo := _travisciRepository.NewTravisCIRepository(s.config)
		travisciUC := _travisciUsecase.NewTravisCIUsecase(s.config, travisciRepo)
		travisciHandler := _travisciDelivery.NewHttpTravisCIHandler(travisciUC)

		v1.GET("/travisci/build", s.cm.UpstreamCacheHandler(travisciHandler.GetTravisCIBuild))
	})
}
