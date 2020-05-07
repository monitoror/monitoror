package service

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/monitoror/monitoror/service/handlers"
	"github.com/monitoror/monitoror/service/middlewares"
	"github.com/monitoror/monitoror/store"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/color"
)

type (
	Server struct {
		// Echo Server
		*echo.Echo

		// CacheMiddleware using CacheStore to return cached data
		CacheMiddleware *middlewares.CacheMiddleware

		store *store.Store
	}
)

var colorer = color.New()

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Init create echo server with middlewares, ui, routes
func Init(store *store.Store) *Server {
	s := &Server{
		store: store,
	}

	s.setupEchoServer()
	s.setupEchoMiddleware()

	InitUI(s)
	InitApis(s)

	return s
}

func (s *Server) Start() error {
	return s.Echo.Start(fmt.Sprintf("%s:%d", s.store.CoreConfig.Address, s.store.CoreConfig.Port))
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
	if s.store.CoreConfig.Debug {
		s.Use(echoMiddleware.LoggerWithConfig(echoMiddleware.LoggerConfig{
			Format: `[-] ` + colorer.Green("${method}") + ` ${uri} status:${status} latency:` + colorer.Green("${latency_human}") + ` error:"${error}"` + "\n",
		}))
	}

	// Cache
	s.CacheMiddleware = middlewares.NewCacheMiddleware(s.store.CacheStore,
		time.Millisecond*time.Duration(s.store.CoreConfig.DownstreamCacheExpiration),
		time.Millisecond*time.Duration(s.store.CoreConfig.UpstreamCacheExpiration),
	) // Used as Handler wrapper in routes
	s.Use(s.CacheMiddleware.DownstreamStoreMiddleware())

	// CORS
	s.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST},
	}))
}
