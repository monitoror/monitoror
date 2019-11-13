//+build faker

package service

import (
	"github.com/monitoror/monitoror/monitorable/azuredevops"
	_azureDevOpsDelivery "github.com/monitoror/monitoror/monitorable/azuredevops/delivery/http"
	_azureDevOpsModels "github.com/monitoror/monitoror/monitorable/azuredevops/models"
	_azureDevOpsUsecase "github.com/monitoror/monitoror/monitorable/azuredevops/usecase"
	"github.com/monitoror/monitoror/monitorable/config"
	_configDelivery "github.com/monitoror/monitoror/monitorable/config/delivery/http"
	_configRepository "github.com/monitoror/monitoror/monitorable/config/repository"
	_configUsecase "github.com/monitoror/monitoror/monitorable/config/usecase"
	"github.com/monitoror/monitoror/monitorable/http"
	_httpDelivery "github.com/monitoror/monitoror/monitorable/http/delivery/http"
	_httpModels "github.com/monitoror/monitoror/monitorable/http/models"
	_httpUsecase "github.com/monitoror/monitoror/monitorable/http/usecase"
	"github.com/monitoror/monitoror/monitorable/jenkins"
	_jenkinsDelivery "github.com/monitoror/monitoror/monitorable/jenkins/delivery/http"
	_jenkinsModels "github.com/monitoror/monitoror/monitorable/jenkins/models"
	_jenkinsUsecase "github.com/monitoror/monitoror/monitorable/jenkins/usecase"
	"github.com/monitoror/monitoror/monitorable/ping"
	_pingDelivery "github.com/monitoror/monitoror/monitorable/ping/delivery/http"
	_pingModels "github.com/monitoror/monitoror/monitorable/ping/models"
	_pingUsecase "github.com/monitoror/monitoror/monitorable/ping/usecase"
	"github.com/monitoror/monitoror/monitorable/pingdom"
	_pingdomDelivery "github.com/monitoror/monitoror/monitorable/pingdom/delivery/http"
	_pingdomModels "github.com/monitoror/monitoror/monitorable/pingdom/models"
	_pingdomUsecase "github.com/monitoror/monitoror/monitorable/pingdom/usecase"
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
	usecase := _configUsecase.NewConfigUsecase(repository, s.store, s.config.DownstreamCacheExpiration)
	delivery := _configDelivery.NewConfigDelivery(usecase)

	s.v1.GET("/config", delivery.GetConfig)

	return usecase
}

func (s *Server) registerPing(configHelper config.Helper) {
	defer logStatus(ping.PingTileType, true)

	usecase := _pingUsecase.NewPingUsecase()
	delivery := _pingDelivery.NewPingDelivery(usecase)

	// Register route to echo
	route := s.v1.GET("/ping", delivery.GetPing)

	// Register param and path to config usecase
	configHelper.RegisterTile(ping.PingTileType, &_pingModels.PingParams{}, route.Path)
}

func (s *Server) registerPort(configHelper config.Helper) {
	defer logStatus(port.PortTileType, true)

	usecase := _portUsecase.NewPortUsecase()
	delivery := _portDelivery.NewPortDelivery(usecase)

	// Register route to echo
	route := s.v1.GET("/port", delivery.GetPort)

	// Register param and path to config usecase
	configHelper.RegisterTile(port.PortTileType, &_portModels.PortParams{}, route.Path)
}

func (s *Server) registerHTTP(configHelper config.Helper) {
	defer logStatus("HTTP", true)

	usecase := _httpUsecase.NewHTTPUsecase()
	delivery := _httpDelivery.NewHTTPDelivery(usecase)

	// Register route to echo
	httpGroup := s.v1.Group("/http")
	routeAny := httpGroup.GET("/any", delivery.GetHTTPAny)
	routeRaw := httpGroup.GET("/raw", delivery.GetHTTPRaw)
	routeJson := httpGroup.GET("/json", delivery.GetHTTPJson)
	routeYaml := httpGroup.GET("/yaml", delivery.GetHTTPYaml)

	// Register data for config hydration
	configHelper.RegisterTile(http.HTTPAnyTileType, &_httpModels.HTTPAnyParams{}, routeAny.Path)
	configHelper.RegisterTile(http.HTTPRawTileType, &_httpModels.HTTPRawParams{}, routeRaw.Path)
	configHelper.RegisterTile(http.HTTPJsonTileType, &_httpModels.HTTPJsonParams{}, routeJson.Path)
	configHelper.RegisterTile(http.HTTPYamlTileType, &_httpModels.HTTPYamlParams{}, routeYaml.Path)
}

func (s *Server) registerPingdom(configHelper config.Helper) {
	defer logStatus(pingdom.PingdomCheckTileType, true)

	usecase := _pingdomUsecase.NewPingdomUsecase()
	delivery := _pingdomDelivery.NewPingdomDelivery(usecase)

	// Register route to echo
	pingdomGroup := s.v1.Group("/pingdom")
	route := pingdomGroup.GET("/check", delivery.GetCheck)

	// Register data for config hydration
	configHelper.RegisterTile(pingdom.PingdomCheckTileType, &_pingdomModels.CheckParams{}, route.Path)
}

func (s *Server) registerTravisCI(configHelper config.Helper) {
	defer logStatus(travisci.TravisCIBuildTileType, true)

	usecase := _travisciUsecase.NewTravisCIUsecase()
	delivery := _travisciDelivery.NewTravisCIDelivery(usecase)

	// Register route to echo
	travisCIGroup := s.v1.Group("/travisci")
	route := travisCIGroup.GET("/build", delivery.GetBuild)

	// Register param and path to config usecase
	configHelper.RegisterTile(travisci.TravisCIBuildTileType, &_travisciModels.BuildParams{}, route.Path)
}

func (s *Server) registerJenkins(configHelper config.Helper) {
	defer logStatus(jenkins.JenkinsBuildTileType, true)

	usecase := _jenkinsUsecase.NewJenkinsUsecase()
	delivery := _jenkinsDelivery.NewJenkinsDelivery(usecase)

	// Register route to echo
	jenkinsGroup := s.v1.Group("/jenkins")
	route := jenkinsGroup.GET("/build", delivery.GetBuild)

	// Register param and path to config usecase
	configHelper.RegisterTile(jenkins.JenkinsBuildTileType, &_jenkinsModels.BuildParams{}, route.Path)
}

func (s *Server) registerAzureDevOps(configHelper config.Helper) {
	defer logStatus("AZURE-DEVOPS", true)

	usecase := _azureDevOpsUsecase.NewAzureDevOpsUsecase()
	delivery := _azureDevOpsDelivery.NewAzureDevOpsDelivery(usecase)

	// Register route to echo
	azureGroup := s.v1.Group("/azuredevops")
	routeBuild := azureGroup.GET("/build", s.cm.UpstreamCacheHandler(delivery.GetBuild))
	routeRelease := azureGroup.GET("/release", s.cm.UpstreamCacheHandler(delivery.GetRelease))

	// Register data for config hydration
	configHelper.RegisterTile(azuredevops.AzureDevOpsBuildTileType, &_azureDevOpsModels.BuildParams{}, routeBuild.Path)
	configHelper.RegisterTile(azuredevops.AzureDevOpsReleaseTileType, &_azureDevOpsModels.ReleaseParams{}, routeRelease.Path)
}
