//+build !faker

package service

import (
	"github.com/jsdidierlaurent/echo-middleware/cache"
	"github.com/monitoror/monitoror/handlers"
	"github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/pkg/monitoror/utils"

	_configDelivery "github.com/monitoror/monitoror/monitorable/config/delivery/http"
	_configRepository "github.com/monitoror/monitoror/monitorable/config/repository"
	_configUsecase "github.com/monitoror/monitoror/monitorable/config/usecase"
	"github.com/monitoror/monitoror/monitorable/ping"
	_pingDelivery "github.com/monitoror/monitoror/monitorable/ping/delivery/http"
	_pingModel "github.com/monitoror/monitoror/monitorable/ping/models"
	_pingRepository "github.com/monitoror/monitoror/monitorable/ping/repository"
	_pingUsecase "github.com/monitoror/monitoror/monitorable/ping/usecase"
	"github.com/monitoror/monitoror/monitorable/port"
	_portDelivery "github.com/monitoror/monitoror/monitorable/port/delivery/http"
	_portModel "github.com/monitoror/monitoror/monitorable/port/models"
	_portRepository "github.com/monitoror/monitoror/monitorable/port/repository"
	_portUsecase "github.com/monitoror/monitoror/monitorable/port/usecase"
	"github.com/monitoror/monitoror/monitorable/travisci"
	_travisciDelivery "github.com/monitoror/monitoror/monitorable/travisci/delivery/http"
	_travisciModel "github.com/monitoror/monitoror/monitorable/travisci/models"
	_travisciRepository "github.com/monitoror/monitoror/monitorable/travisci/repository"
	_travisciUsecase "github.com/monitoror/monitoror/monitorable/travisci/usecase"
)

func (s *Server) initApis() {
	v1 := s.Group("/api/v1")

	// Init monitorableParams for config verification
	monitorableParams := make(map[tiles.TileType]utils.Validator)

	// ------------- INFO ------------- //
	infoHandler := handlers.HttpInfoHandler(s.config)
	v1.GET("/info", s.cm.UpstreamCacheHandlerWithExpiration(cache.NEVER, infoHandler.GetInfo))

	// ------------- PING ------------- //
	register("ping", true, func() {
		pingRepo := _pingRepository.NewPingRepository(s.config)
		pingUC := _pingUsecase.NewPingUsecase(pingRepo)
		pingHandler := _pingDelivery.NewHttpPingHandler(pingUC)

		monitorableParams[ping.PingTileType] = &_pingModel.PingParams{}

		v1.GET("/ping", s.cm.UpstreamCacheHandler(pingHandler.GetPing))
	})

	// ------------- PORT ------------- //
	register("port", true, func() {
		portRepo := _portRepository.NewPortRepository(s.config)
		portUC := _portUsecase.NewPortUsecase(portRepo)
		portHandler := _portDelivery.NewHttpPortHandler(portUC)

		monitorableParams[port.PortTileType] = &_portModel.PortParams{}

		v1.GET("/port", s.cm.UpstreamCacheHandler(portHandler.GetPort))
	})

	// ------------- TRAVIS CI ------------- //
	register("travis-ci", s.config.Monitorable.TravisCI.Url != "", func() {
		travisciRepo := _travisciRepository.NewTravisCIRepository(s.config)
		travisciUC := _travisciUsecase.NewTravisCIUsecase(s.config, travisciRepo)
		travisciHandler := _travisciDelivery.NewHttpTravisCIHandler(travisciUC)

		monitorableParams[travisci.TravisCIBuildTileType] = &_travisciModel.BuildParams{}

		v1.GET("/travisci/build", s.cm.UpstreamCacheHandler(travisciHandler.GetTravisCIBuild))
	})

	// ------------- CONFIG ------------- //
	configRepo := _configRepository.NewConfigRepository()
	configUC := _configUsecase.NewConfigUsecase(monitorableParams, configRepo)
	configHandler := _configDelivery.NewHttpConfigHandler(configUC)
	v1.GET("/config", s.cm.UpstreamCacheHandler(configHandler.GetConfig))

}
