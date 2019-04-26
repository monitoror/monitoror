package service

import (
	"fmt"
	"net/http"

	rice "github.com/GeertJohan/go.rice"

	"github.com/labstack/gommon/color"

	"github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/handlers"
	"github.com/monitoror/monitoror/middlewares"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type Server struct {
	// Echo Server
	*echo.Echo

	// Config
	config *config.Config

	// Middleware
	cm *middlewares.CacheMiddleware
}

var colorer = color.New()

// Init create echo server with middlewares, front, routes
func Init(config *config.Config) *Server {
	server := &Server{config: config}

	server.initEcho()
	server.initMiddleware()
	server.initFront()
	server.initApis()

	return server
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
	register("front", s.config.Env == "production", func() {
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
	})
}

func (s *Server) Start() {
	fmt.Println()
	log.Fatal(s.Echo.Start(fmt.Sprintf(":%d", s.config.Port)))
}

// register route utility function (used for print if route is enabled at start
func register(name string, enabled bool, handle func()) {
	if enabled {
		handle()
		fmt.Printf("⇨ %s: %s\n", name, colorer.Green("enabled"))
	} else {
		fmt.Printf("⇨ %s: %s\n", name, colorer.Red("disabled"))
	}
}
