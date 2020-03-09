//+build faker

package service

import (
	"math/rand"
	"time"

	"github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/handlers"
	"github.com/monitoror/monitoror/monitorable/azuredevops"
	azureDevOpsDelivery "github.com/monitoror/monitoror/monitorable/azuredevops/delivery/http"
	azureDevOpsModels "github.com/monitoror/monitoror/monitorable/azuredevops/models"
	azureDevOpsUsecase "github.com/monitoror/monitoror/monitorable/azuredevops/usecase"
	configDelivery "github.com/monitoror/monitoror/monitorable/config/delivery/http"
	configRepository "github.com/monitoror/monitoror/monitorable/config/repository"
	configUsecase "github.com/monitoror/monitoror/monitorable/config/usecase"
	"github.com/monitoror/monitoror/monitorable/github"
	githubDelivery "github.com/monitoror/monitoror/monitorable/github/delivery/http"
	githubModels "github.com/monitoror/monitoror/monitorable/github/models"
	githubUsecase "github.com/monitoror/monitoror/monitorable/github/usecase"
	"github.com/monitoror/monitoror/monitorable/http"
	httpDelivery "github.com/monitoror/monitoror/monitorable/http/delivery/http"
	httpModels "github.com/monitoror/monitoror/monitorable/http/models"
	httpUsecase "github.com/monitoror/monitoror/monitorable/http/usecase"
	"github.com/monitoror/monitoror/monitorable/jenkins"
	jenkinsDelivery "github.com/monitoror/monitoror/monitorable/jenkins/delivery/http"
	jenkinsModels "github.com/monitoror/monitoror/monitorable/jenkins/models"
	jenkinsUsecase "github.com/monitoror/monitoror/monitorable/jenkins/usecase"
	"github.com/monitoror/monitoror/monitorable/ping"
	pingDelivery "github.com/monitoror/monitoror/monitorable/ping/delivery/http"
	pingModels "github.com/monitoror/monitoror/monitorable/ping/models"
	pingUsecase "github.com/monitoror/monitoror/monitorable/ping/usecase"
	"github.com/monitoror/monitoror/monitorable/pingdom"
	pingdomDelivery "github.com/monitoror/monitoror/monitorable/pingdom/delivery/http"
	pingdomModels "github.com/monitoror/monitoror/monitorable/pingdom/models"
	pingdomUsecase "github.com/monitoror/monitoror/monitorable/pingdom/usecase"
	"github.com/monitoror/monitoror/monitorable/port"
	portDelivery "github.com/monitoror/monitoror/monitorable/port/delivery/http"
	portModels "github.com/monitoror/monitoror/monitorable/port/models"
	portUsecase "github.com/monitoror/monitoror/monitorable/port/usecase"
	"github.com/monitoror/monitoror/monitorable/travisci"
	travisciDelivery "github.com/monitoror/monitoror/monitorable/travisci/delivery/http"
	travisciModels "github.com/monitoror/monitoror/monitorable/travisci/models"
	travisciUsecase "github.com/monitoror/monitoror/monitorable/travisci/usecase"

	"github.com/jsdidierlaurent/echo-middleware/cache"
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
	registerTile(s.registerPing, config.DefaultVariant, true)
	registerTile(s.registerPort, config.DefaultVariant, true)
	registerTile(s.registerHTTP, config.DefaultVariant, true)
	registerTile(s.registerPingdom, config.DefaultVariant, true)
	registerTile(s.registerTravisCI, config.DefaultVariant, true)
	registerTile(s.registerJenkins, config.DefaultVariant, true)
	registerTile(s.registerAzureDevOps, config.DefaultVariant, true)
	registerTile(s.registerGithub, config.DefaultVariant, true)
}

func (s *Server) registerInfo() {
	infoDelivery := handlers.NewHTTPInfoDelivery()
	s.api.GET("/info", s.cm.UpstreamCacheHandlerWithExpiration(cache.NEVER, infoDelivery.GetInfo))
}

func (s *Server) registerConfig() {
	repository := configRepository.NewConfigRepository()
	usecase := configUsecase.NewConfigUsecase(repository, s.store, s.config.DownstreamCacheExpiration)
	delivery := configDelivery.NewConfigDelivery(usecase)

	s.api.GET("/config", delivery.GetConfig)
	s.configHelper = usecase
}

func (s *Server) registerPing(variant string) {
	usecase := pingUsecase.NewPingUsecase()
	delivery := pingDelivery.NewPingDelivery(usecase)

	// Register route to echo
	route := s.api.GET("/ping", delivery.GetPing)

	// Register param and path to config usecase
	s.configHelper.RegisterTile(ping.PingTileType, &pingModels.PingParams{}, route.Path, config.DefaultInitialMaxDelay)
}

func (s *Server) registerPort(variant string) {
	usecase := portUsecase.NewPortUsecase()
	delivery := portDelivery.NewPortDelivery(usecase)

	// Register route to echo
	route := s.api.GET("/port", delivery.GetPort)

	// Register param and path to config usecase
	s.configHelper.RegisterTile(port.PortTileType, &portModels.PortParams{}, route.Path, config.DefaultInitialMaxDelay)
}

func (s *Server) registerHTTP(variant string) {
	usecase := httpUsecase.NewHTTPUsecase()
	delivery := httpDelivery.NewHTTPDelivery(usecase)

	// Register route to echo
	httpGroup := s.api.Group("/http")
	routeStatus := httpGroup.GET("/status", delivery.GetHTTPStatus)
	routeRaw := httpGroup.GET("/raw", delivery.GetHTTPRaw)
	routeJson := httpGroup.GET("/formatted", delivery.GetHTTPFormatted)

	// Register data for config hydration
	s.configHelper.RegisterTile(http.HTTPStatusTileType, &httpModels.HTTPStatusParams{}, routeStatus.Path, config.DefaultInitialMaxDelay)
	s.configHelper.RegisterTile(http.HTTPRawTileType, &httpModels.HTTPRawParams{}, routeRaw.Path, config.DefaultInitialMaxDelay)
	s.configHelper.RegisterTile(http.HTTPFormattedTileType, &httpModels.HTTPFormattedParams{}, routeJson.Path, config.DefaultInitialMaxDelay)
}

func (s *Server) registerPingdom(variant string) {
	usecase := pingdomUsecase.NewPingdomUsecase()
	delivery := pingdomDelivery.NewPingdomDelivery(usecase)

	// Register route to echo
	pingdomGroup := s.api.Group("/pingdom")
	route := pingdomGroup.GET("/check", delivery.GetCheck)

	// Register data for config hydration
	s.configHelper.RegisterTile(pingdom.PingdomCheckTileType, &pingdomModels.CheckParams{}, route.Path, config.DefaultInitialMaxDelay)
}

func (s *Server) registerTravisCI(variant string) {
	usecase := travisciUsecase.NewTravisCIUsecase()
	delivery := travisciDelivery.NewTravisCIDelivery(usecase)

	// Register route to echo
	travisCIGroup := s.api.Group("/travisci")
	route := travisCIGroup.GET("/build", delivery.GetBuild)

	// Register param and path to config usecase
	s.configHelper.RegisterTile(travisci.TravisCIBuildTileType, &travisciModels.BuildParams{}, route.Path, config.DefaultInitialMaxDelay)
}

func (s *Server) registerJenkins(variant string) {
	usecase := jenkinsUsecase.NewJenkinsUsecase()
	delivery := jenkinsDelivery.NewJenkinsDelivery(usecase)

	// Register route to echo
	jenkinsGroup := s.api.Group("/jenkins")
	route := jenkinsGroup.GET("/build", delivery.GetBuild)

	// Register param and path to config usecase
	s.configHelper.RegisterTile(jenkins.JenkinsBuildTileType, &jenkinsModels.BuildParams{}, route.Path, config.DefaultInitialMaxDelay)
}

func (s *Server) registerAzureDevOps(variant string) {
	usecase := azureDevOpsUsecase.NewAzureDevOpsUsecase()
	delivery := azureDevOpsDelivery.NewAzureDevOpsDelivery(usecase)

	// Register route to echo
	azureGroup := s.api.Group("/azuredevops")
	routeBuild := azureGroup.GET("/build", s.cm.UpstreamCacheHandler(delivery.GetBuild))
	routeRelease := azureGroup.GET("/release", s.cm.UpstreamCacheHandler(delivery.GetRelease))

	// Register data for config hydration
	s.configHelper.RegisterTile(azuredevops.AzureDevOpsBuildTileType, &azureDevOpsModels.BuildParams{}, routeBuild.Path, config.DefaultInitialMaxDelay)
	s.configHelper.RegisterTile(azuredevops.AzureDevOpsReleaseTileType, &azureDevOpsModels.ReleaseParams{}, routeRelease.Path, config.DefaultInitialMaxDelay)
}

func (s *Server) registerGithub(variant string) {
	usecase := githubUsecase.NewGithubUsecase()
	delivery := githubDelivery.NewGithubDelivery(usecase)

	// Register route to echo
	azureGroup := s.api.Group("/github")
	routeCount := azureGroup.GET("/count", s.cm.UpstreamCacheHandler(delivery.GetCount))
	routeChecks := azureGroup.GET("/checks", s.cm.UpstreamCacheHandler(delivery.GetChecks))

	// Register data for config hydration
	s.configHelper.RegisterTile(github.GithubCountTileType, &githubModels.CountParams{}, routeCount.Path, config.DefaultInitialMaxDelay)
	s.configHelper.RegisterTile(github.GithubChecksTileType, &githubModels.ChecksParams{}, routeChecks.Path, config.DefaultInitialMaxDelay)
}
