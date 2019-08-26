//+build faker

package service

import (
	"github.com/monitoror/monitoror/monitorable/config"
	_configDelivery "github.com/monitoror/monitoror/monitorable/config/delivery/http"
	_configRepository "github.com/monitoror/monitoror/monitorable/config/repository"
	_configUsecase "github.com/monitoror/monitoror/monitorable/config/usecase"
	"github.com/monitoror/monitoror/monitorable/jenkins"
	_jenkinsDelivery "github.com/monitoror/monitoror/monitorable/jenkins/delivery/http"
	_jenkinsModels "github.com/monitoror/monitoror/monitorable/jenkins/models"
	_jenkinsUsecase "github.com/monitoror/monitoror/monitorable/jenkins/usecase"
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

func (s *Server) registerConfig() config.Helper {
	repository := _configRepository.NewConfigRepository()
	usecase := _configUsecase.NewConfigUsecase(repository)
	delivery := _configDelivery.NewHttpConfigDelivery(usecase)

	s.v1.GET("/config", delivery.GetConfig)

	return usecase
}

func (s *Server) registerPing(configHelper config.Helper) {
	defer logStatus(ping.PingTileType, true)

	usecase := _pingUsecase.NewPingUsecase()
	delivery := _pingDelivery.NewHttpPingDelivery(usecase)

	// Register route to echo
	route := s.v1.GET("/ping", delivery.GetPing)

	// Register param and path to config usecase
	configHelper.RegisterTile(ping.PingTileType, &_pingModels.PingParams{}, route.Path)
}

func (s *Server) registerPort(configHelper config.Helper) {
	defer logStatus(port.PortTileType, true)

	usecase := _portUsecase.NewPortUsecase()
	delivery := _portDelivery.NewHttpPortDelivery(usecase)

	// Register route to echo
	route := s.v1.GET("/port", delivery.GetPort)

	// Register param and path to config usecase
	configHelper.RegisterTile(port.PortTileType, &_portModels.PortParams{}, route.Path)
}

func (s *Server) registerTravisCI(configHelper config.Helper) {
	defer logStatus(travisci.TravisCIBuildTileType, true)

	usecase := _travisciUsecase.NewTravisCIUsecase()
	delivery := _travisciDelivery.NewHttpTravisCIDelivery(usecase)

	// Register route to echo
	travisCIGroup := s.v1.Group("/travisci")
	route := travisCIGroup.GET("/build", delivery.GetBuild)

	// Register param and path to config usecase
	configHelper.RegisterTile(travisci.TravisCIBuildTileType, &_travisciModels.BuildParams{}, route.Path)
}

func (s *Server) registerJenkins(configHelper config.Helper) {
	defer logStatus(jenkins.JenkinsBuildTileType, true)

	usecase := _jenkinsUsecase.NewJenkinsUsecase()
	delivery := _jenkinsDelivery.NewHttpJenkinsDelivery(usecase)

	// Register route to echo
	jenkinsGroup := s.v1.Group("/jenkins")
	route := jenkinsGroup.GET("/build", delivery.GetBuild)

	// Register param and path to config usecase
	configHelper.RegisterTile(jenkins.JenkinsBuildTileType, &_jenkinsModels.BuildParams{}, route.Path)
}
