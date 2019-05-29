//+build faker

package service

import (
	"github.com/labstack/echo/v4"
	"github.com/monitoror/monitoror/monitorable/config"
	"github.com/monitoror/monitoror/monitorable/ping"
	_pingDelivery "github.com/monitoror/monitoror/monitorable/ping/delivery/http"
	_pingModels "github.com/monitoror/monitoror/monitorable/ping/models"
	_pingUsecase "github.com/monitoror/monitoror/monitorable/ping/usecase"
	"github.com/monitoror/monitoror/monitorable/port"
	_portDelivery "github.com/monitoror/monitoror/monitorable/port/delivery/http"
	_portModels "github.com/monitoror/monitoror/monitorable/port/models"
	_portUsecase "github.com/monitoror/monitoror/monitorable/port/usecase"
	"github.com/monitoror/monitoror/monitorable/travisci"
	_travisciDelivery "github.com/monitoror/monitoror/monitorable/travisci/delivery/http"
	_travisciModels "github.com/monitoror/monitoror/monitorable/travisci/models"
	_travisciUsecase "github.com/monitoror/monitoror/monitorable/travisci/usecase"
)

func (s *Server) registerPing(g *echo.Group, registerer config.Regiterer) {
	path := "/ping"
	tileType := ping.PingTileType

	factory := func() (handler echo.HandlerFunc) {
		usecase := _pingUsecase.NewPingUsecase()
		delivery := _pingDelivery.NewHttpPingDelivery(usecase)

		// Registering route param
		registerer.Register(tileType, path, &_pingModels.PingParams{})

		handler = delivery.GetPing
		return
	}

	s.register(g, true, path, tileType, factory)
}

func (s *Server) registerPort(g *echo.Group, registerer config.Regiterer) {
	path := "/port"
	tileType := port.PortTileType

	factory := func() (handler echo.HandlerFunc) {
		usecase := _portUsecase.NewPortUsecase()
		delivery := _portDelivery.NewHttpPortDelivery(usecase)

		// Registering route param
		registerer.Register(tileType, path, &_portModels.PortParams{})

		handler = delivery.GetPort
		return
	}

	s.register(g, true, path, tileType, factory)
}

func (s *Server) registerTravisCIBuild(g *echo.Group, registerer config.Regiterer) {
	valid := s.config.Monitorable.TravisCI.Url != ""
	path := "/travisci/build"
	tileType := travisci.TravisCIBuildTileType

	factory := func() (handler echo.HandlerFunc) {
		usecase := _travisciUsecase.NewTravisCIUsecase()
		delivery := _travisciDelivery.NewHttpTravisCIDelivery(usecase)

		// Registering route param
		registerer.Register(tileType, path, &_travisciModels.BuildParams{})

		handler = delivery.GetTravisCIBuild
		return
	}

	s.register(g, valid, path, tileType, factory)
}
