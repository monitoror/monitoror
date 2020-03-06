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
	"github.com/monitoror/monitoror/monitorable/gitlab"
	_gitlabDelivery "github.com/monitoror/monitoror/monitorable/gitlab/delivery/http"
	_gitlabModels "github.com/monitoror/monitoror/monitorable/gitlab/models"
	_gitlabRepository "github.com/monitoror/monitoror/monitorable/gitlab/repository"
	_gitlabUsecase "github.com/monitoror/monitoror/monitorable/gitlab/usecase"
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
	for variant, conf := range s.config.Monitorable.Gitlab {
		registerTile(s.registerGitlab, variant, conf.IsValid())
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
	s.configHelper.RegisterTile(ping.PingTileType, &_pingModels.PingParams{}, route.Path, pingConfig.InitialMaxDelay)
}

func (s *Server) registerPort(variant string) {
	portConfig := s.config.Monitorable.Port[variant]

	repository := _portRepository.NewPortRepository(portConfig)
	usecase := _portUsecase.NewPortUsecase(repository)
	delivery := _portDelivery.NewPortDelivery(usecase)

	// Register route to echo
	route := s.api.GET("/port", s.cm.UpstreamCacheHandler(delivery.GetPort))

	// Register data for config hydration
	s.configHelper.RegisterTile(port.PortTileType, &_portModels.PortParams{}, route.Path, portConfig.InitialMaxDelay)
}

func (s *Server) registerHTTP(variant string) {
	httpConfig := s.config.Monitorable.HTTP[variant]

	repository := _httpRepository.NewHTTPRepository(httpConfig)
	usecase := _httpUsecase.NewHTTPUsecase(repository, s.store, s.config.UpstreamCacheExpiration)
	delivery := _httpDelivery.NewHTTPDelivery(usecase)

	// Register route to echo
	httpGroup := s.api.Group("/http")
	routeStatus := httpGroup.GET("/status", s.cm.UpstreamCacheHandler(delivery.GetHTTPStatus))
	routeRaw := httpGroup.GET("/raw", s.cm.UpstreamCacheHandler(delivery.GetHTTPRaw))
	routeJSON := httpGroup.GET("/formatted", s.cm.UpstreamCacheHandler(delivery.GetHTTPFormatted))

	// Register data for config hydration
	s.configHelper.RegisterTile(http.HTTPStatusTileType,
		&_httpModels.HTTPStatusParams{}, routeStatus.Path, httpConfig.InitialMaxDelay)
	s.configHelper.RegisterTile(http.HTTPRawTileType,
		&_httpModels.HTTPRawParams{}, routeRaw.Path, httpConfig.InitialMaxDelay)
	s.configHelper.RegisterTile(http.HTTPFormattedTileType,
		&_httpModels.HTTPFormattedParams{}, routeJSON.Path, httpConfig.InitialMaxDelay)
}

func (s *Server) registerPingdom(variant string) {
	pingdomConfig := s.config.Monitorable.Pingdom[variant]

	repository := _pingdomRepository.NewPingdomRepository(pingdomConfig)
	usecase := _pingdomUsecase.NewPingdomUsecase(repository, pingdomConfig, s.store)
	delivery := _pingdomDelivery.NewPingdomDelivery(usecase)

	// Register route to echo
	pingdomGroup := s.api.Group(fmt.Sprintf("/pingdom/%s", variant))
	route := pingdomGroup.GET("/check", s.cm.UpstreamCacheHandler(delivery.GetCheck))

	// Register data for config hydration
	s.configHelper.RegisterTileWithConfigVariant(pingdom.PingdomCheckTileType,
		variant, &_pingdomModels.CheckParams{}, route.Path, pingdomConfig.InitialMaxDelay)
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
	s.configHelper.RegisterTileWithConfigVariant(travisci.TravisCIBuildTileType,
		variant, &_travisciModels.BuildParams{}, route.Path, travisCIConfig.InitialMaxDelay)
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
		variant, &_jenkinsModels.BuildParams{}, route.Path, jenkinsConfig.InitialMaxDelay)
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
		variant, &_azureDevOpsModels.BuildParams{}, routeBuild.Path, azureDevOpsConfig.InitialMaxDelay)
	s.configHelper.RegisterTileWithConfigVariant(azuredevops.AzureDevOpsReleaseTileType,
		variant, &_azureDevOpsModels.ReleaseParams{}, routeRelease.Path, azureDevOpsConfig.InitialMaxDelay)
}

func (s *Server) registerGithub(variant string) {
	githubConfig := s.config.Monitorable.Github[variant]
	// Custom UpstreamCacheExpiration only for count because github has no-cache for this query and the rate limit is 30req/Hour
	countCacheExpiration := time.Millisecond * time.Duration(githubConfig.CountCacheExpiration)

	repository := _githubRepository.NewGithubRepository(githubConfig)
	usecase := _githubUsecase.NewGithubUsecase(repository)
	delivery := _githubDelivery.NewGithubDelivery(usecase)

	// Register route to echo
	azureGroup := s.api.Group(fmt.Sprintf("/github/%s", variant))
	routeCount := azureGroup.GET("/count", s.cm.UpstreamCacheHandlerWithExpiration(countCacheExpiration, delivery.GetCount))
	routeChecks := azureGroup.GET("/checks", s.cm.UpstreamCacheHandler(delivery.GetChecks))

	// Register data for config hydration
	s.configHelper.RegisterTileWithConfigVariant(github.GithubCountTileType,
		variant, &_githubModels.CountParams{}, routeCount.Path, githubConfig.InitialMaxDelay)
	s.configHelper.RegisterTileWithConfigVariant(github.GithubChecksTileType,
		variant, &_githubModels.ChecksParams{}, routeChecks.Path, githubConfig.InitialMaxDelay)
	s.configHelper.RegisterDynamicTileWithConfigVariant(github.GithubPullRequestTileType,
		variant, &_githubModels.PullRequestParams{}, usecase)
}

func (s *Server) registerGitlab(variant string) {
	gitlabConfig := s.config.Monitorable.Gitlab[variant]
	// Custom UpstreamCacheExpiration only for count because gitlab has no-cache for this query and the rate limit is 30req/Hour
	countCacheExpiration := time.Millisecond * time.Duration(gitlabConfig.CountCacheExpiration)

	repository := _gitlabRepository.NewGitlabRepository(gitlabConfig)
	usecase := _gitlabUsecase.NewGitlabUsecase(repository)
	delivery := _gitlabDelivery.NewGitlabDelivery(usecase)

	// Register route to echo
	gitlabGroup := s.api.Group(fmt.Sprintf("/gitlab/%s", variant))
	routeCount := gitlabGroup.GET("/count", s.cm.UpstreamCacheHandlerWithExpiration(countCacheExpiration, delivery.GetCount))
	routePipelines := gitlabGroup.GET("/pipelines", s.cm.UpstreamCacheHandler(delivery.GetPipelines))

	// Register data for config hydration
	s.configHelper.RegisterTileWithConfigVariant(gitlab.GitlabCountTileType,
		variant, &_gitlabModels.CountParams{}, routeCount.Path, gitlabConfig.InitialMaxDelay)
	s.configHelper.RegisterTileWithConfigVariant(gitlab.GitlabPipelinesTileType,
		variant, &_gitlabModels.PipelinesParams{}, routePipelines.Path, gitlabConfig.InitialMaxDelay)
	s.configHelper.RegisterDynamicTileWithConfigVariant(gitlab.GitlabPullRequestTileType,
		variant, &_gitlabModels.MergeRequestParams{}, usecase)
}
