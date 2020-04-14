package service

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/monitoror/monitoror/cli"
	"github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/pkg/system"
	"github.com/monitoror/monitoror/service/handlers"
	"github.com/monitoror/monitoror/service/middlewares"
	"github.com/monitoror/monitoror/service/registry"
	"github.com/monitoror/monitoror/service/store"

	"github.com/jsdidierlaurent/echo-middleware/cache"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/color"
	"github.com/labstack/gommon/log"
)

type (
	Server struct {
		// Echo Server
		*echo.Echo

		store *store.Store
	}
)

var colorer = color.New()

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Init create echo server with middlewares, ui, routes
func Init(config *config.Config, cli cli.CLI) *Server {
	s := &Server{
		store: &store.Store{
			CoreConfig: config,
			Cli:        cli,
			Registry:   registry.NewRegistry(),
		},
	}

	s.setupEchoServer()
	s.setupEchoMiddleware()

	InitUI(s)
	InitApis(s)

	return s
}

func (s *Server) Start() {
	s.store.Cli.PrintServerStartup(system.GetNetworkIP(), s.store.CoreConfig.Port)
	log.Fatal(s.Echo.Start(fmt.Sprintf(":%d", s.store.CoreConfig.Port)))
}

func (s *Server) setupEchoServer() {
	s.Echo = echo.New()
	s.HideBanner = true
	s.HidePort = true

	// ----- Errors Handler -----
	s.HTTPErrorHandler = handlers.HTTPErrorHandler
}

func (s *Server) setupEchoMiddleware() {
	// Recover (don't panic 😎)
	s.Use(echoMiddleware.Recover())

	// Log requests
	if s.store.CoreConfig.Env != "production" {
		s.Use(echoMiddleware.LoggerWithConfig(echoMiddleware.LoggerConfig{
			Format: `[-] ` + colorer.Green("${method}") + ` ${uri} status:${status} latency:` + colorer.Green("${latency_human}") + ` error:"${error}"` + "\n",
		}))
	}

	// Cache
	s.store.CacheStore = cache.NewGoCacheStore(time.Minute*5, time.Second) // Default value, always override
	s.store.CacheMiddleware = middlewares.NewCacheMiddleware(s.store.CacheStore,
		time.Millisecond*time.Duration(s.store.CoreConfig.DownstreamCacheExpiration),
		time.Millisecond*time.Duration(s.store.CoreConfig.UpstreamCacheExpiration),
	) // Used as Handler wrapper in routes
	s.Use(s.store.CacheMiddleware.DownstreamStoreMiddleware())

	// CORS
	s.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST},
	}))
}
