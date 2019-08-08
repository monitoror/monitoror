package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

const EnvPrefix = "MO"
const DefaultVariant = "default"

type (
	// Backend Configuration
	Config struct {
		// --- General Configuration ---
		Port int    `json:"port"` // Default: 8080
		Env  string `json:"env"`  // Default: production

		// --- Cache Configuration ---
		// UpstreamCache is used to respond before executing the request. Avoid overloading services.
		UpstreamCache Cache `json:"upstreamCache"`
		// DownstreamCache is used to respond after executing the request in case of timeout error.
		DownstreamCache Cache `json:"downstreamCache"`

		// Monitorable Config
		Monitorable Monitorable `json:"monitorable"`
	}

	Cache struct {
		Expire          int `json:"expire"`          // In Millisecond
		CleanupInterval int `json:"cleanupInterval"` // In Millisecond
	}

	Monitorable struct {
		Ping     Ping                 `json:"ping"`
		Port     Port                 `json:"port"`
		Github   map[string]*Github   `json:"github"`
		TravisCI map[string]*TravisCI `json:"travisCI"`
		Jenkins  map[string]*Jenkins  `json:"jenkins"`
	}

	Ping struct {
		Count    int `json:"count"`
		Timeout  int `json:"timeout"`  // In Millisecond
		Interval int `json:"interval"` // In Millisecond
	}

	Port struct {
		Timeout int `json:"timeout"` // In Millisecond
	}

	Github struct {
		Token string `json:"token"`
	}

	TravisCI struct {
		Url     string `json:"url"`
		Timeout int    `json:"timeout"` // In Millisecond
		Token   string `json:"token"`
	}

	Jenkins struct {
		Url       string `json:"url"`
		Timeout   int    `json:"timeout"` // In Millisecond
		SSLVerify bool   `json:"sslVerify"`
		Login     string `json:"login"`
		Token     string `json:"token"`
	}
)

// Load configuration from configuration file / env / default value
func InitConfig() *Config {
	var config Config

	// Setup Env
	viper.AutomaticEnv()
	viper.SetEnvPrefix(EnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Transform Env and define Label for setting default value
	variants := initEnvAndVariant()

	// Setup default values
	// --- General Configuration ---
	viper.SetDefault("Port", 8080)
	viper.SetDefault("Env", "production")

	// --- Cache Configuration ---
	viper.SetDefault("UpstreamCache.Expire", 10000)
	viper.SetDefault("UpstreamCache.CleanupInterval", 1000)
	viper.SetDefault("DownstreamCache.Expire", 120000)
	viper.SetDefault("DownstreamCache.CleanupInterval", 10000)

	// --- Ping Configuration ---
	viper.SetDefault("Monitorable.Ping.Count", 2)
	viper.SetDefault("Monitorable.Ping.Timeout", 1000)
	viper.SetDefault("Monitorable.Ping.Interval", 100)

	// --- Port Configuration ---
	viper.SetDefault("Monitorable.Port.Timeout", 2000)

	// --- Github Configuration ---
	for variant := range variants["Github"] {
		viper.SetDefault(fmt.Sprintf("Monitorable.Github.%s.Token", variant), "")
	}

	// --- TravisCI Configuration ---
	for variant := range variants["TravisCI"] {
		viper.SetDefault(fmt.Sprintf("Monitorable.TravisCI.%s.Url", variant), "https://api.travis-ci.org/")
		viper.SetDefault(fmt.Sprintf("Monitorable.TravisCI.%s.Timeout", variant), 2000)
		viper.SetDefault(fmt.Sprintf("Monitorable.TravisCI.%s.Token", variant), "")
	}

	// --- Jenkins Configuration ---
	for variant := range variants["Jenkins"] {
		viper.SetDefault(fmt.Sprintf("Monitorable.Jenkins.%s.Timeout", variant), 2000)
		viper.SetDefault(fmt.Sprintf("Monitorable.Jenkins.%s.SSLVerify", variant), true)
		viper.SetDefault(fmt.Sprintf("Monitorable.Jenkins.%s.Url", variant), "")
		viper.SetDefault(fmt.Sprintf("Monitorable.Jenkins.%s.Login", variant), "")
		viper.SetDefault(fmt.Sprintf("Monitorable.Jenkins.%s.Token", variant), "")
	}

	_ = viper.Unmarshal(&config)

	return &config
}

func (t *TravisCI) IsValid() bool {
	return t.Url != ""
}

func (t *Jenkins) IsValid() bool {
	return t.Url != ""
}
