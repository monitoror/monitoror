//+build faker

package service

import (
	"math/rand"
	"time"

	"github.com/monitoror/monitoror/monitorable/github"

	"github.com/jsdidierlaurent/echo-middleware/cache"
	. "github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/handlers"
	"github.com/monitoror/monitoror/monitorable/azuredevops"
	_azureDevOpsDelivery "github.com/monitoror/monitoror/monitorable/azuredevops/delivery/http"
	_azureDevOpsModels "github.com/monitoror/monitoror/monitorable/azuredevops/models"
	_azureDevOpsUsecase "github.com/monitoror/monitoror/monitorable/azuredevops/usecase"
	_configDelivery "github.com/monitoror/monitoror/monitorable/config/delivery/http"
	_configRepository "github.com/monitoror/monitoror/monitorable/config/repository"
	_configUsecase "github.com/monitoror/monitoror/monitorable/config/usecase"
	_githubDelivery "github.com/monitoror/monitoror/monitorable/github/delivery/http"
	_githubModels "github.com/monitoror/monitoror/monitorable/github/models"
	_githubUsecase "github.com/monitoror/monitoror/monitorable/github/usecase"
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

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (s *Server) initApis() {
	// Api group definition
	s.api = s.Group("/api/v1")

	// ------------- INFO ------------- //
	s.registerInfo()

	// ------------- CONFIG ------------- //
	s.registerConfig()

	// ------------- TILES ------------- //
	registerTile(s.registerPing, DefaultVariant, true)
	registerTile(s.registerPort, DefaultVariant, true)
	registerTile(s.registerHTTP, DefaultVariant, true)
	registerTile(s.registerPingdom, DefaultVariant, true)
	registerTile(s.registerTravisCI, DefaultVariant, true)
	registerTile(s.registerJenkins, DefaultVariant, true)
	registerTile(s.registerAzureDevOps, DefaultVariant, true)
	registerTile(s.registerGithub, DefaultVariant, true)
}

func (s *Server) registerInfo() {
	infoDelivery := handlers.NewHTTPInfoDelivery()
	s.api.GET("/info", s.cm.UpstreamCacheHandlerWithExpiration(cache.NEVER, infoDelivery.GetInfo))
}

func (s *Server) registerConfig() {
	repository := _configRepository.NewConfigRepository()
	usecase := _configUsecase.NewConfigUsecase(repository, s.store, s.config.DownstreamCacheExpiration)
	delivery := _configDelivery.NewConfigDelivery(usecase)

	s.api.GET("/config", delivery.GetConfig)
	s.configHelper = usecase
}

func (s *Server) registerPing(variant string) {
	usecase := _pingUsecase.NewPingUsecase()
	delivery := _pingDelivery.NewPingDelivery(usecase)

	// Register route to echo
	route := s.api.GET("/ping", delivery.GetPing)

	// Register param and path to config usecase
	s.configHelper.RegisterTile(ping.PingTileType, &_pingModels.PingParams{}, route.Path)
}

func (s *Server) registerPort(variant string) {
	usecase := _portUsecase.NewPortUsecase()
	delivery := _portDelivery.NewPortDelivery(usecase)

	// Register route to echo
	route := s.api.GET("/port", delivery.GetPort)

	// Register param and path to config usecase
	s.configHelper.RegisterTile(port.PortTileType, &_portModels.PortParams{}, route.Path)
}

func (s *Server) registerHTTP(variant string) {
	usecase := _httpUsecase.NewHTTPUsecase()
	delivery := _httpDelivery.NewHTTPDelivery(usecase)

	// Register route to echo
	httpGroup := s.api.Group("/http")
	routeAny := httpGroup.GET("/any", delivery.GetHTTPAny)
	routeRaw := httpGroup.GET("/raw", delivery.GetHTTPRaw)
	routeJson := httpGroup.GET("/formatted", delivery.GetHTTPFormatted)

	// Register data for config hydration
	s.configHelper.RegisterTile(http.HTTPAnyTileType, &_httpModels.HTTPAnyParams{}, routeAny.Path)
	s.configHelper.RegisterTile(http.HTTPRawTileType, &_httpModels.HTTPRawParams{}, routeRaw.Path)
	s.configHelper.RegisterTile(http.HTTPFormattedTileType, &_httpModels.HTTPFormattedParams{}, routeJson.Path)
}

func (s *Server) registerPingdom(variant string) {
	usecase := _pingdomUsecase.NewPingdomUsecase()
	delivery := _pingdomDelivery.NewPingdomDelivery(usecase)

	// Register route to echo
	pingdomGroup := s.api.Group("/pingdom")
	route := pingdomGroup.GET("/check", delivery.GetCheck)

	// Register data for config hydration
	s.configHelper.RegisterTile(pingdom.PingdomCheckTileType, &_pingdomModels.CheckParams{}, route.Path)
}

func (s *Server) registerTravisCI(variant string) {
	usecase := _travisciUsecase.NewTravisCIUsecase()
	delivery := _travisciDelivery.NewTravisCIDelivery(usecase)

	// Register route to echo
	travisCIGroup := s.api.Group("/travisci")
	route := travisCIGroup.GET("/build", delivery.GetBuild)

	// Register param and path to config usecase
	s.configHelper.RegisterTile(travisci.TravisCIBuildTileType, &_travisciModels.BuildParams{}, route.Path)
}

func (s *Server) registerJenkins(variant string) {
	usecase := _jenkinsUsecase.NewJenkinsUsecase()
	delivery := _jenkinsDelivery.NewJenkinsDelivery(usecase)

	// Register route to echo
	jenkinsGroup := s.api.Group("/jenkins")
	route := jenkinsGroup.GET("/build", delivery.GetBuild)

	// Register param and path to config usecase
	s.configHelper.RegisterTile(jenkins.JenkinsBuildTileType, &_jenkinsModels.BuildParams{}, route.Path)
}

func (s *Server) registerAzureDevOps(variant string) {
	usecase := _azureDevOpsUsecase.NewAzureDevOpsUsecase()
	delivery := _azureDevOpsDelivery.NewAzureDevOpsDelivery(usecase)

	// Register route to echo
	azureGroup := s.api.Group("/azuredevops")
	routeBuild := azureGroup.GET("/build", s.cm.UpstreamCacheHandler(delivery.GetBuild))
	routeRelease := azureGroup.GET("/release", s.cm.UpstreamCacheHandler(delivery.GetRelease))

	// Register data for config hydration
	s.configHelper.RegisterTile(azuredevops.AzureDevOpsBuildTileType, &_azureDevOpsModels.BuildParams{}, routeBuild.Path)
	s.configHelper.RegisterTile(azuredevops.AzureDevOpsReleaseTileType, &_azureDevOpsModels.ReleaseParams{}, routeRelease.Path)
}

func (s *Server) registerGithub(variant string) {
	usecase := _githubUsecase.NewGithubUsecase()
	delivery := _githubDelivery.NewGithubDelivery(usecase)

	// Register route to echo
	azureGroup := s.api.Group("/github")
	routeCount := azureGroup.GET("/count", s.cm.UpstreamCacheHandler(delivery.GetCount))
	routeChecks := azureGroup.GET("/checks", s.cm.UpstreamCacheHandler(delivery.GetChecks))

	// Register data for config hydration
	s.configHelper.RegisterTile(github.GithubCountTileType, &_githubModels.CountParams{}, routeCount.Path)
	s.configHelper.RegisterTile(github.GithubChecksTileType, &_githubModels.ChecksParams{}, routeChecks.Path)
}
