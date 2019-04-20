//+build faker

package service

import (
	"github.com/monitoror/monitoror/handlers"

	_pingDelivery "github.com/monitoror/monitoror/monitorable/ping/delivery/http"
	_pingUsecase "github.com/monitoror/monitoror/monitorable/ping/usecase"
	_portDelivery "github.com/monitoror/monitoror/monitorable/port/delivery/http"
	_portUsecase "github.com/monitoror/monitoror/monitorable/port/usecase"
	_travisciDelivery "github.com/monitoror/monitoror/monitorable/travisci/delivery/http"
	_travisciUsecase "github.com/monitoror/monitoror/monitorable/travisci/usecase"
)

func (s *Server) initApis() {
	v1 := s.Group("/api/v1")

	// Register route utils function

	// ------------- INFO ------------- //
	register("info", true, func() {
		infoHandler := handlers.HttpInfoHandler(s.config)
		v1.GET("/info", infoHandler.GetInfo)
	})

	// ------------- PING ------------- //
	register("ping", true, func() {
		pingUC := _pingUsecase.NewPingUsecase()
		pingHandler := _pingDelivery.NewHttpPingHandler(pingUC)

		v1.GET("/ping", pingHandler.GetPing)
	})

	// ------------- PORT ------------- //
	register("port", true, func() {
		portUC := _portUsecase.NewPortUsecase()
		portHandler := _portDelivery.NewHttpPortHandler(portUC)

		v1.GET("/port", portHandler.GetPort)
	})

	// ------------- TRAVIS CI ------------- //
	register("travis-ci", s.config.Monitorable.TravisCI.Url != "", func() {
		travisciUC := _travisciUsecase.NewTravisCIUsecase()
		travisciHandler := _travisciDelivery.NewHttpTravisCIHandler(travisciUC)

		v1.GET("/travisci/build", travisciHandler.GetTravisCIBuild)
	})
}
