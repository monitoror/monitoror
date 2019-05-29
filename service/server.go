package service

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jsdidierlaurent/echo-middleware/cache"

	"github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/handlers"
	"github.com/monitoror/monitoror/middlewares"
	"github.com/monitoror/monitoror/models/tiles"
	_configDelivery "github.com/monitoror/monitoror/monitorable/config/delivery/http"
	_configRepository "github.com/monitoror/monitoror/monitorable/config/repository"
	_configUsecase "github.com/monitoror/monitoror/monitorable/config/usecase"

	rice "github.com/GeertJohan/go.rice"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/color"
	"github.com/labstack/gommon/log"
)

type (
	Server struct {
		// Echo Server
		*echo.Echo

		// Config
		config *config.Config

		// Middleware
		cm *middlewares.CacheMiddleware
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

	//  ----- Default Logger -----
	log.SetPrefix("")
	log.SetHeader("[${level}]")
	log.SetLevel(log.INFO)

	//  ----- Echo Logger -----
	s.Logger.SetPrefix("")
	s.Logger.SetHeader("[${level}]")
	s.Logger.SetLevel(log.INFO)

	// ----- Errors Handler -----
	s.HTTPErrorHandler = handlers.HttpErrorHandler
}

func (s *Server) initMiddleware() {
	// Recover (don't panic 😎)
	s.Use(echoMiddleware.Recover())

	// Log requests
	s.Use(echoMiddleware.LoggerWithConfig(echoMiddleware.LoggerConfig{
		Format: `[-] ` + colorer.Green("${method}") + ` ${uri} status:${status} error:"${error}"` + "\n",
	}))

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
	if s.config.Env != "production" {
		fmt.Printf("⇨ %s: %s\n", "front", colorer.Red("disabled"))
		return
	}

	// Never use constant or variable according to docs : https://github.com/GeertJohan/go.rice#calling-findbox-and-mustfindbox
	frontAssets, err := rice.FindBox("../front/dist")
	if err != nil {
		panic("static front/dist not found. Build them with `cd front && yarn run build` first.")
	}

	assetHandler := http.FileServer(frontAssets.HTTPBox())
	s.GET("/", echo.WrapHandler(assetHandler))
	s.GET("/css/*", echo.WrapHandler(http.StripPrefix("/", assetHandler)))
	s.GET("/js/*", echo.WrapHandler(http.StripPrefix("/", assetHandler)))
	s.GET("/fonts/*", echo.WrapHandler(http.StripPrefix("/", assetHandler)))
	s.GET("/img/*", echo.WrapHandler(http.StripPrefix("/", assetHandler)))

	fmt.Printf("⇨ %s: %s\n", "front", colorer.Green("enabled"))
}

func (s *Server) initApis() {
	v1 := s.Group("/api/v1")

	// ------------- INFO ------------- //
	infoHandler := handlers.HttpInfoHandler(s.config)
	v1.GET("/info", s.cm.UpstreamCacheHandlerWithExpiration(cache.NEVER, infoHandler.GetInfo))

	// ------------- CONFIG ------------- //
	configRepo := _configRepository.NewConfigRepository()
	configUC := _configUsecase.NewConfigUsecase(configRepo)
	configHandler := _configDelivery.NewHttpConfigDelivery(configUC)
	v1.GET("/config", s.cm.UpstreamCacheHandler(configHandler.GetConfig))

	// ------------- PING ------------- //
	s.registerPing(v1, configUC)

	// ------------- PORT ------------- //
	s.registerPort(v1, configUC)

	// ------------- TRAVIS-CI ------------- //
	s.registerTravisCIBuild(v1, configUC)
}

func (s *Server) register(g *echo.Group, valid bool, path string, tileType tiles.TileType, factory func() (handler echo.HandlerFunc)) {
	if !valid {
		fmt.Printf("⇨ %s: %s\n", strings.ToLower(string(tileType)), colorer.Red("disabled"))
		return
	}

	// Registering route
	g.GET(path, factory())

	fmt.Printf("⇨ %s: %s\n", strings.ToLower(string(tileType)), colorer.Green("enabled"))
}
