package service

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/handlers"
	"github.com/monitoror/monitoror/middlewares"
	monitorableConfig "github.com/monitoror/monitoror/monitorable/config"

	rice "github.com/GeertJohan/go.rice"
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

		// Config
		store        cache.Store
		config       *config.Config
		configHelper monitorableConfig.Helper

		// Middleware
		cm *middlewares.CacheMiddleware

		// Groups
		api *echo.Group
	}
)

var colorer = color.New()

// Init create echo server with middlewares, ui, routes
func Init(config *config.Config) *Server {
	server := &Server{
		config: config,
	}

	server.initEcho()
	server.initMiddleware()
	server.initUI()
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
	s.HTTPErrorHandler = handlers.HTTPErrorHandler
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

func (s *Server) initUI() {
	loadUI := s.config.Env == "production"
	defer logStatus("UI", loadUI)

	if !loadUI {
		return
	}

	// Remove LocateAppend and LocateFS because we don't use it and it cause tests issue
	riceConfig := rice.Config{
		LocateOrder: []rice.LocateMethod{rice.LocateEmbedded},
	}
	// Never use constant or variable according to docs : https://github.com/GeertJohan/go.rice#calling-findbox-and-mustfindbox
	uiAssets, err := riceConfig.FindBox("../ui/dist")
	if err != nil {
		panic("static ui/dist not found. Build them with `cd ui && yarn build` first.")
	}

	assetHandler := http.FileServer(uiAssets.HTTPBox())
	s.GET("/", echo.WrapHandler(assetHandler))
	s.GET("/favicon*", echo.WrapHandler(assetHandler))
	s.GET("/css/*", echo.WrapHandler(http.StripPrefix("/", assetHandler)))
	s.GET("/js/*", echo.WrapHandler(http.StripPrefix("/", assetHandler)))
	s.GET("/fonts/*", echo.WrapHandler(http.StripPrefix("/", assetHandler)))
	s.GET("/img/*", echo.WrapHandler(http.StripPrefix("/", assetHandler)))
}

// registerTile is a decorator function to print enable/disable log and call handler when enabled
func registerTile(handler func(string), variant string, enabled bool) {
	var nameWithVariant string
	name := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
	name = name[strings.LastIndex(name, ".")+1:]
	name = strings.Replace(name, "register", "", 1)
	name = strings.Replace(name, "-fm", "", 1)
	name = strings.ToUpper(name)

	if variant != config.DefaultVariant && variant != "" {
		nameWithVariant = fmt.Sprintf("%v (variant: %s)", name, variant)
	} else {
		nameWithVariant = fmt.Sprintf("%v", name)
	}

	logStatus(nameWithVariant, enabled)
	if enabled {
		handler(variant)
	}
}

func logStatus(name interface{}, enabled bool) {
	status := colorer.Green("enabled")
	if !enabled {
		status = colorer.Red("disabled")
	}
	fmt.Printf("⇨ %s: %s\n", name, status)
}
