//+build !faker

package service

import (
	"fmt"

	"github.com/monitoror/monitoror/monitorable/config"
	_configDelivery "github.com/monitoror/monitoror/monitorable/config/delivery/http"
	_configRepository "github.com/monitoror/monitoror/monitorable/config/repository"
	_configUsecase "github.com/monitoror/monitoror/monitorable/config/usecase"

	"github.com/monitoror/monitoror/monitorable/jenkins"
	_jenkinsDelivery "github.com/monitoror/monitoror/monitorable/jenkins/delivery/http"
	_jenkinsModels "github.com/monitoror/monitoror/monitorable/jenkins/models"
	_jenkinsRepository "github.com/monitoror/monitoror/monitorable/jenkins/repository"
	_jenkinsUsecase "github.com/monitoror/monitoror/monitorable/jenkins/usecase"
	"github.com/monitoror/monitoror/monitorable/ping"
	_pingDelivery "github.com/monitoror/monitoror/monitorable/ping/delivery/http"
	_pingModels "github.com/monitoror/monitoror/monitorable/ping/models"
	_pingRepository "github.com/monitoror/monitoror/monitorable/ping/repository"
	_pingUsecase "github.com/monitoror/monitoror/monitorable/ping/usecase"
	"github.com/monitoror/monitoror/monitorable/port"
	_portDelivery "github.com/monitoror/monitoror/monitorable/port/delivery/http"
	_portModels "github.com/monitoror/monitoror/monitorable/port/models"
	_portRepository "github.com/monitoror/monitoror/monitorable/port/repository"
	_portUsecase "github.com/monitoror/monitoror/monitorable/port/usecase"
	"github.com/monitoror/monitoror/monitorable/travisci"
	_travisciDelivery "github.com/monitoror/monitoror/monitorable/travisci/delivery/http"
	_travisciModels "github.com/monitoror/monitoror/monitorable/travisci/models"
	_travisciRepository "github.com/monitoror/monitoror/monitorable/travisci/repository"
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

	repository := _pingRepository.NewPingRepository(s.config)
	usecase := _pingUsecase.NewPingUsecase(repository)
	delivery := _pingDelivery.NewHttpPingDelivery(usecase)

	// Register route to echo
	route := s.v1.GET("/ping", s.cm.UpstreamCacheHandler(delivery.GetPing))

	// Register data for config hydration
	configHelper.RegisterTile(ping.PingTileType, &_pingModels.PingParams{}, route.Path)
}

func (s *Server) registerPort(configHelper config.Helper) {
	defer logStatus(port.PortTileType, true)

	repository := _portRepository.NewPortRepository(s.config)
	usecase := _portUsecase.NewPortUsecase(repository)
	delivery := _portDelivery.NewHttpPortDelivery(usecase)

	// Register route to echo
	route := s.v1.GET("/port", s.cm.UpstreamCacheHandler(delivery.GetPort))

	// Register data for config hydration
	configHelper.RegisterTile(port.PortTileType, &_portModels.PortParams{}, route.Path)
}

//noinspection ALL
func (s *Server) registerTravisCI(configHelper config.Helper) {
	for variant, travisCIConf := range s.config.Monitorable.TravisCI {
		// Associate github config
		githubConf := s.config.Monitorable.Github[variant]

		defer logStatusWithConfigVariant("TRAVISCI", variant, travisCIConf.IsValid())
		if !travisCIConf.IsValid() {
			continue
		}

		repository := _travisciRepository.NewTravisCIRepository(travisCIConf, githubConf)
		usecase := _travisciUsecase.NewTravisCIUsecase(repository)
		delivery := _travisciDelivery.NewHttpTravisCIDelivery(usecase)

		// Register route to echo
		travisCIGroup := s.v1.Group(fmt.Sprintf("/travisci/%s", variant))
		route := travisCIGroup.GET("/build", s.cm.UpstreamCacheHandler(delivery.GetBuild))

		// Register data for config hydration
		configHelper.RegisterTileWithConfigVariant(travisci.TravisCIBuildTileType, variant, &_travisciModels.BuildParams{}, route.Path)
	}
}

//noinspection ALL
func (s *Server) registerJenkins(configHelper config.Helper) {
	for variant, jenkinsConf := range s.config.Monitorable.Jenkins {
		defer logStatusWithConfigVariant("JENKINS", variant, jenkinsConf.IsValid())
		if !jenkinsConf.IsValid() {
			continue
		}

		repository := _jenkinsRepository.NewJenkinsRepository(jenkinsConf)
		usecase := _jenkinsUsecase.NewJenkinsUsecase(repository, s.config.DownstreamCache)
		delivery := _jenkinsDelivery.NewHttpJenkinsDelivery(usecase)

		// Register route to echo
		jenkinsGroup := s.v1.Group(fmt.Sprintf("/jenkins/%s", variant))
		route := jenkinsGroup.GET("/build", s.cm.UpstreamCacheHandler(delivery.GetBuild))

		// Register data for config hydration
		configHelper.RegisterTileWithConfigVariant(jenkins.JenkinsBuildTileType,
			variant, &_jenkinsModels.BuildParams{}, route.Path)
		configHelper.RegisterDynamicTileWithConfigVariant(jenkins.JenkinsMultiBranchTileType,
			variant, &_jenkinsModels.MultiBranchParams{}, usecase)
	}
}
