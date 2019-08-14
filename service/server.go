package service

import (
	"fmt"
	"net/http"

	"github.com/labstack/gommon/log"

	"github.com/jsdidierlaurent/echo-middleware/cache"

	rice "github.com/GeertJohan/go.rice"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/color"
	"github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/handlers"
	"github.com/monitoror/monitoror/middlewares"
)

type (
	Server struct {
		// Echo Server
		*echo.Echo

		// GetConfig
		config *config.Config

		// Middleware
		cm *middlewares.CacheMiddleware

		// Groups
		api *echo.Group
		v1  *echo.Group
	}
)

var colorer = color.New()

// Init create echo server with middlewares, front, routes
func Init(config *config.Config) *Server {
	server := &Server{
		config: config,
	}

	server.initEcho()
	server.initMiddleware()
	server.initFront()
	server.initApis()

	return server
}

func (s *Server) Start() {
	fmt.Println()
	log.Fatal(s.Echo.Start(fmt.Sprintf(":%d", s.config.Port)))
}

func (s *Server) initEcho() {
	s.Echo = echo.New()
	s.HideBanner = true

	// ----- Errors Handler -----
	s.HTTPErrorHandler = handlers.HttpErrorHandler
}

func (s *Server) initMiddleware() {
	// Recover (don't panic 😎)
	s.Use(echoMiddleware.Recover())

	// Log requests
	if s.config.Env != "production" {
		s.Use(echoMiddleware.LoggerWithConfig(echoMiddleware.LoggerConfig{
			Format: `[-] ` + colorer.Green("${method}") + ` ${uri} status:${status} latency:` + colorer.Green("${latency_human}") + ` error:"${error}"` + "\n",
		}))
	}

	// Cache
	s.cm = middlewares.NewCacheMiddleware(s.config) // Used as Handler wrapper in routes
	s.Use(s.cm.DownstreamStoreMiddleware())

	// CORS
	s.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST},
	}))
}

func (s *Server) initFront() {
	loadFront := s.config.Env == "production"
	defer logStatus("FRONT", loadFront)

	if !loadFront {
		return
	}

	// Never use constant or variable according to docs : https://github.com/GeertJohan/go.rice#calling-findbox-and-mustfindbox
	frontAssets, err := rice.FindBox("../front/dist")
	if err != nil {
		panic("static front/dist not found. GetBuildStatus them with `cd front && yarn run build` first.")
	}

	assetHandler := http.FileServer(frontAssets.HTTPBox())
	s.GET("/", echo.WrapHandler(assetHandler))
	s.GET("/css/*", echo.WrapHandler(http.StripPrefix("/", assetHandler)))
	s.GET("/js/*", echo.WrapHandler(http.StripPrefix("/", assetHandler)))
	s.GET("/fonts/*", echo.WrapHandler(http.StripPrefix("/", assetHandler)))
	s.GET("/img/*", echo.WrapHandler(http.StripPrefix("/", assetHandler)))
}

func (s *Server) initApis() {
	// Api group definition
	s.api = s.Group("/api")

	// V1
	s.v1 = s.api.Group("/v1")

	// ------------- INFO ------------- //
	infoDelivery := handlers.NewHttpInfoDelivery()
	s.v1.GET("/info", s.cm.UpstreamCacheHandlerWithExpiration(cache.NEVER, infoDelivery.GetInfo))

	// ------------- CONFIG ------------- //
	configHelper := s.registerConfig()

	// ------------- PING ------------- //
	s.registerPing(configHelper)

	// ------------- PORT ------------- //
	s.registerPort(configHelper)

	// ------------- TRAVIS-CI ------------- //
	s.registerTravisCI(configHelper)

	// ------------- JENKINS ------------- //
	s.registerJenkins(configHelper)
}

func logStatus(name interface{}, enabled bool) {
	status := colorer.Green("enabled")
	if !enabled {
		status = colorer.Red("disabled")
	}
	fmt.Printf("⇨ %s: %s\n", name, status)
}

func logStatusWithConfigVariant(name interface{}, variant string, enabled bool) {
	var nameWithVariant string
	if variant != config.DefaultVariant && variant != "" {
		nameWithVariant = fmt.Sprintf("%v (%s)", name, variant)
	} else {
		nameWithVariant = fmt.Sprintf("%v", name)
	}

	logStatus(nameWithVariant, enabled)
}
