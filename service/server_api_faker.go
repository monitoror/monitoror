//+build faker

package service

import (
	"github.com/monitoror/monitoror/handlers"
	"github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/pkg/monitoror/utils"

	_configDelivery "github.com/monitoror/monitoror/monitorable/config/delivery/http"
	_configRepository "github.com/monitoror/monitoror/monitorable/config/repository"
	_configUsecase "github.com/monitoror/monitoror/monitorable/config/usecase"
	"github.com/monitoror/monitoror/monitorable/ping"
	_pingDelivery "github.com/monitoror/monitoror/monitorable/ping/delivery/http"
	_pingModel "github.com/monitoror/monitoror/monitorable/ping/models"
	_pingUsecase "github.com/monitoror/monitoror/monitorable/ping/usecase"
	"github.com/monitoror/monitoror/monitorable/port"
	_portDelivery "github.com/monitoror/monitoror/monitorable/port/delivery/http"
	_portModel "github.com/monitoror/monitoror/monitorable/port/models"
	_portUsecase "github.com/monitoror/monitoror/monitorable/port/usecase"
	"github.com/monitoror/monitoror/monitorable/travisci"
	_travisciDelivery "github.com/monitoror/monitoror/monitorable/travisci/delivery/http"
	_travisciModel "github.com/monitoror/monitoror/monitorable/travisci/models"
	_travisciUsecase "github.com/monitoror/monitoror/monitorable/travisci/usecase"
)

func (s *Server) initApis() {
	v1 := s.Group("/api/v1")

	// Init monitorableParams for config verification
	monitorableParams := make(map[tiles.TileType]utils.Validator)

	// ------------- INFO ------------- //
	infoHandler := handlers.HttpInfoHandler(s.config)
	v1.GET("/info", infoHandler.GetInfo)

	// ------------- PING ------------- //
	register("ping", true, func() {
		pingUC := _pingUsecase.NewPingUsecase()
		pingHandler := _pingDelivery.NewHttpPingHandler(pingUC)

		monitorableParams[ping.PingTileType] = &_pingModel.PingParams{}

		v1.GET("/ping", pingHandler.GetPing)
	})

	// ------------- PORT ------------- //
	register("port", true, func() {
		portUC := _portUsecase.NewPortUsecase()
		portHandler := _portDelivery.NewHttpPortHandler(portUC)

		monitorableParams[port.PortTileType] = &_portModel.PortParams{}

		v1.GET("/port", portHandler.GetPort)
	})

	// ------------- TRAVIS CI ------------- //
	register("travis-ci", s.config.Monitorable.TravisCI.Url != "", func() {
		travisciUC := _travisciUsecase.NewTravisCIUsecase()
		travisciHandler := _travisciDelivery.NewHttpTravisCIHandler(travisciUC)

		monitorableParams[travisci.TravisCIBuildTileType] = &_travisciModel.BuildParams{}

		v1.GET("/travisci/build", travisciHandler.GetTravisCIBuild)
	})

	// ------------- CONFIG ------------- //
	configRepo := _configRepository.NewConfigRepository()
	configUC := _configUsecase.NewConfigUsecase(monitorableParams, configRepo)
	configHandler := _configDelivery.NewHttpConfigHandler(configUC)
	v1.GET("/config", configHandler.GetConfig)

}
