package service

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jsdidierlaurent/echo-middleware/cache"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/color"
	"github.com/labstack/gommon/log"
	"github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/handlers"
	"github.com/monitoror/monitoror/middlewares"
)

type (
	Server struct {
		// Echo Server
		*echo.Echo

		// Config
		config *config.Config
		store  cache.Store

		// Middleware
		cm *middlewares.CacheMiddleware
	}
)

var colorer = color.New()

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Init create echo server with middlewares, ui, routes
func Init(config *config.Config) *Server {
	server := &Server{
		config: config,
	}

	server.setupEchoServer()
	server.setupEchoMiddleware()

	InitUI(server)
	InitApis(server)

	return server
}

func (s *Server) Start() {
	fmt.Println()
	log.Fatal(s.Echo.Start(fmt.Sprintf(":%d", s.config.Port)))
}

func (s *Server) setupEchoServer() {
	s.Echo = echo.New()
	s.HideBanner = true

	// ----- Errors Handler -----
	s.HTTPErrorHandler = handlers.HTTPErrorHandler
}

func (s *Server) setupEchoMiddleware() {
	// Recover (don't panic 😎)
	s.Use(echoMiddleware.Recover())

	// Log requests
	if s.config.Env != "production" {
		s.Use(echoMiddleware.LoggerWithConfig(echoMiddleware.LoggerConfig{
			Format: `[-] ` + colorer.Green("${method}") + ` ${uri} status:${status} latency:` + colorer.Green("${latency_human}") + ` error:"${error}"` + "\n",
		}))
	}

	// Cache
	s.store = cache.NewGoCacheStore(time.Minute*5, time.Second) // Default value, always override
	s.cm = middlewares.NewCacheMiddleware(s.store,
		time.Millisecond*time.Duration(s.config.DownstreamCacheExpiration),
		time.Millisecond*time.Duration(s.config.UpstreamCacheExpiration),
	) // Used as Handler wrapper in routes
	s.Use(s.cm.DownstreamStoreMiddleware())

	// CORS
	s.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST},
	}))
}
