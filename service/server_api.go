//+build !faker

package service

import (
	"fmt"
	"time"

	"github.com/monitoror/monitoror/pkg/monitoror/utils/system"

	"github.com/monitoror/monitoror/handlers"
	"github.com/monitoror/monitoror/monitorable/azuredevops"
	_azureDevOpsDelivery "github.com/monitoror/monitoror/monitorable/azuredevops/delivery/http"
	_azureDevOpsModels "github.com/monitoror/monitoror/monitorable/azuredevops/models"
	_azureDevOpsRepository "github.com/monitoror/monitoror/monitorable/azuredevops/repository"
	_azureDevOpsUsecase "github.com/monitoror/monitoror/monitorable/azuredevops/usecase"
	_configDelivery "github.com/monitoror/monitoror/monitorable/config/delivery/http"
	_configRepository "github.com/monitoror/monitoror/monitorable/config/repository"
	_configUsecase "github.com/monitoror/monitoror/monitorable/config/usecase"
	"github.com/monitoror/monitoror/monitorable/github"
	_githubDelivery "github.com/monitoror/monitoror/monitorable/github/delivery/http"
	_githubModels "github.com/monitoror/monitoror/monitorable/github/models"
	_githubRepository "github.com/monitoror/monitoror/monitorable/github/repository"
	_githubUsecase "github.com/monitoror/monitoror/monitorable/github/usecase"
	"github.com/monitoror/monitoror/monitorable/http"
	_httpDelivery "github.com/monitoror/monitoror/monitorable/http/delivery/http"
	_httpModels "github.com/monitoror/monitoror/monitorable/http/models"
	_httpRepository "github.com/monitoror/monitoror/monitorable/http/repository"
	_httpUsecase "github.com/monitoror/monitoror/monitorable/http/usecase"
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
	"github.com/monitoror/monitoror/monitorable/pingdom"
	_pingdomDelivery "github.com/monitoror/monitoror/monitorable/pingdom/delivery/http"
	_pingdomModels "github.com/monitoror/monitoror/monitorable/pingdom/models"
	_pingdomRepository "github.com/monitoror/monitoror/monitorable/pingdom/repository"
	_pingdomUsecase "github.com/monitoror/monitoror/monitorable/pingdom/usecase"
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

	"github.com/jsdidierlaurent/echo-middleware/cache"
)

func (s *Server) initApis() {
	// Api group definition
	s.api = s.Group("/api/v1")

	// ------------- INFO ------------- //
	s.registerInfo()

	// ------------- CONFIG ------------- //
	s.registerConfig()

	// ------------- TILES ------------- //
	for variant := range s.config.Monitorable.Ping {
		registerTile(s.registerPing, variant, system.IsRawSocketAvailable())
	}
	for variant := range s.config.Monitorable.Port {
		registerTile(s.registerPort, variant, true)
	}
	for variant := range s.config.Monitorable.HTTP {
		registerTile(s.registerHTTP, variant, true)
	}
	for variant, conf := range s.config.Monitorable.Pingdom {
		registerTile(s.registerPingdom, variant, conf.IsValid())
	}
	for variant, conf := range s.config.Monitorable.TravisCI {
		registerTile(s.registerTravisCI, variant, conf.IsValid())
	}
	for variant, conf := range s.config.Monitorable.Jenkins {
		registerTile(s.registerJenkins, variant, conf.IsValid())
	}
	for variant, conf := range s.config.Monitorable.AzureDevOps {
		registerTile(s.registerAzureDevOps, variant, conf.IsValid())
	}
	for variant, conf := range s.config.Monitorable.Github {
		registerTile(s.registerGithub, variant, conf.IsValid())
	}
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
	pingConfig := s.config.Monitorable.Ping[variant]

	repository := _pingRepository.NewPingRepository(pingConfig)
	usecase := _pingUsecase.NewPingUsecase(repository)
	delivery := _pingDelivery.NewPingDelivery(usecase)

	// Register route to echo
	route := s.api.GET("/ping", s.cm.UpstreamCacheHandler(delivery.GetPing))

	// Register data for config hydration
	s.configHelper.RegisterTile(ping.PingTileType, &_pingModels.PingParams{}, route.Path)
}

func (s *Server) registerPort(variant string) {
	portConfig := s.config.Monitorable.Port[variant]

	repository := _portRepository.NewPortRepository(portConfig)
	usecase := _portUsecase.NewPortUsecase(repository)
	delivery := _portDelivery.NewPortDelivery(usecase)

	// Register route to echo
	route := s.api.GET("/port", s.cm.UpstreamCacheHandler(delivery.GetPort))

	// Register data for config hydration
	s.configHelper.RegisterTile(port.PortTileType, &_portModels.PortParams{}, route.Path)
}

func (s *Server) registerHTTP(variant string) {
	httpConfig := s.config.Monitorable.HTTP[variant]

	repository := _httpRepository.NewHTTPRepository(httpConfig)
	usecase := _httpUsecase.NewHTTPUsecase(repository, s.store, s.config.UpstreamCacheExpiration)
	delivery := _httpDelivery.NewHTTPDelivery(usecase)

	// Register route to echo
	httpGroup := s.api.Group("/http")
	routeAny := httpGroup.GET("/any", s.cm.UpstreamCacheHandler(delivery.GetHTTPAny))
	routeRaw := httpGroup.GET("/raw", s.cm.UpstreamCacheHandler(delivery.GetHTTPRaw))
	routeJSON := httpGroup.GET("/formatted", s.cm.UpstreamCacheHandler(delivery.GetHTTPFormatted))

	// Register data for config hydration
	s.configHelper.RegisterTile(http.HTTPAnyTileType, &_httpModels.HTTPAnyParams{}, routeAny.Path)
	s.configHelper.RegisterTile(http.HTTPRawTileType, &_httpModels.HTTPRawParams{}, routeRaw.Path)
	s.configHelper.RegisterTile(http.HTTPFormattedTileType, &_httpModels.HTTPFormattedParams{}, routeJSON.Path)
}

func (s *Server) registerPingdom(variant string) {
	pingdomVariant := s.config.Monitorable.Pingdom[variant]

	repository := _pingdomRepository.NewPingdomRepository(pingdomVariant)
	usecase := _pingdomUsecase.NewPingdomUsecase(repository, pingdomVariant, s.store)
	delivery := _pingdomDelivery.NewPingdomDelivery(usecase)

	// Register route to echo
	pingdomGroup := s.api.Group(fmt.Sprintf("/pingdom/%s", variant))
	route := pingdomGroup.GET("/check", s.cm.UpstreamCacheHandler(delivery.GetCheck))

	// Register data for config hydration
	s.configHelper.RegisterTileWithConfigVariant(pingdom.PingdomCheckTileType,
		variant, &_pingdomModels.CheckParams{}, route.Path)
	s.configHelper.RegisterDynamicTileWithConfigVariant(pingdom.PingdomChecksTileType,
		variant, &_pingdomModels.ChecksParams{}, usecase)
}

func (s *Server) registerTravisCI(variant string) {
	travisCIConfig := s.config.Monitorable.TravisCI[variant]

	repository := _travisciRepository.NewTravisCIRepository(travisCIConfig)
	usecase := _travisciUsecase.NewTravisCIUsecase(repository)
	delivery := _travisciDelivery.NewTravisCIDelivery(usecase)

	// Register route to echo
	travisCIGroup := s.api.Group(fmt.Sprintf("/travisci/%s", variant))
	route := travisCIGroup.GET("/build", s.cm.UpstreamCacheHandler(delivery.GetBuild))

	// Register data for config hydration
	s.configHelper.RegisterTileWithConfigVariant(travisci.TravisCIBuildTileType, variant, &_travisciModels.BuildParams{}, route.Path)
}

func (s *Server) registerJenkins(variant string) {
	jenkinsConfig := s.config.Monitorable.Jenkins[variant]

	repository := _jenkinsRepository.NewJenkinsRepository(jenkinsConfig)
	usecase := _jenkinsUsecase.NewJenkinsUsecase(repository)
	delivery := _jenkinsDelivery.NewJenkinsDelivery(usecase)

	// Register route to echo
	jenkinsGroup := s.api.Group(fmt.Sprintf("/jenkins/%s", variant))
	route := jenkinsGroup.GET("/build", s.cm.UpstreamCacheHandler(delivery.GetBuild))

	// Register data for config hydration
	s.configHelper.RegisterTileWithConfigVariant(jenkins.JenkinsBuildTileType,
		variant, &_jenkinsModels.BuildParams{}, route.Path)
	s.configHelper.RegisterDynamicTileWithConfigVariant(jenkins.JenkinsMultiBranchTileType,
		variant, &_jenkinsModels.MultiBranchParams{}, usecase)
}

func (s *Server) registerAzureDevOps(variant string) {
	azureDevOpsConfig := s.config.Monitorable.AzureDevOps[variant]

	repository := _azureDevOpsRepository.NewAzureDevOpsRepository(azureDevOpsConfig)
	usecase := _azureDevOpsUsecase.NewAzureDevOpsUsecase(repository)
	delivery := _azureDevOpsDelivery.NewAzureDevOpsDelivery(usecase)

	// Register route to echo
	azureGroup := s.api.Group(fmt.Sprintf("/azuredevops/%s", variant))
	routeBuild := azureGroup.GET("/build", s.cm.UpstreamCacheHandler(delivery.GetBuild))
	routeRelease := azureGroup.GET("/release", s.cm.UpstreamCacheHandler(delivery.GetRelease))

	// Register data for config hydration
	s.configHelper.RegisterTileWithConfigVariant(azuredevops.AzureDevOpsBuildTileType,
		variant, &_azureDevOpsModels.BuildParams{}, routeBuild.Path)
	s.configHelper.RegisterTileWithConfigVariant(azuredevops.AzureDevOpsReleaseTileType,
		variant, &_azureDevOpsModels.ReleaseParams{}, routeRelease.Path)
}

func (s *Server) registerGithub(variant string) {
	githubConfig := s.config.Monitorable.Github[variant]
	// Custom UpstreamCacheExpiration only for Issues because github has no-cache for this query and the rate limit is 30req/Minutes
	issueCacheExpiration := time.Millisecond * time.Duration(githubConfig.IssueCacheExpiration)

	repository := _githubRepository.NewGithubRepository(githubConfig)
	usecase := _githubUsecase.NewGithubUsecase(repository)
	delivery := _githubDelivery.NewGithubDelivery(usecase)

	// Register route to echo
	azureGroup := s.api.Group(fmt.Sprintf("/github/%s", variant))
	routeIssues := azureGroup.GET("/issues", s.cm.UpstreamCacheHandlerWithExpiration(issueCacheExpiration, delivery.GetIssues))
	routeChecks := azureGroup.GET("/checks", s.cm.UpstreamCacheHandler(delivery.GetChecks))

	// Register data for config hydration
	s.configHelper.RegisterTileWithConfigVariant(github.GithubIssuesTileType,
		variant, &_githubModels.IssuesParams{}, routeIssues.Path)
	s.configHelper.RegisterTileWithConfigVariant(github.GithubChecksTileType,
		variant, &_githubModels.ChecksParams{}, routeChecks.Path)
	s.configHelper.RegisterDynamicTileWithConfigVariant(github.GithubPullRequestTileType,
		variant, &_githubModels.PullRequestParams{}, usecase)
}
