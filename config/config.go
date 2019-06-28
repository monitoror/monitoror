package config

import (
	"strings"

	"github.com/spf13/viper"
)

const EnvPrefix = "MO"

type (
	//Backend Configuration
	Config struct {
		// --- General Configuration ---
		Port int    `json:"port"` // Default: 8080
		Env  string `json:"env"`  // Default: production

		// --- Cache Configuration ---
		//UpstreamCache is used to respond before executing the request. Avoid overloading services.
		UpstreamCache Cache `json:"upstreamCache"`
		//DownstreamCache is used to respond after executing the request in case of timeout error.
		DownstreamCache Cache `json:"downstreamCache"`

		//Monitorables GetConfig
		Monitorable Monitorable `json:"monitorable"`
	}

	Cache struct {
		Expire          int `json:"expire"`          // In Millisecond
		CleanupInterval int `json:"cleanupInterval"` // In Millisecond
	}

	Monitorable struct {
		Ping     Ping     `json:"ping"`
		Port     Port     `json:"port"`
		Gitlab   Gitlab   `json:"gitlab"`
		Github   Github   `json:"github"`
		TravisCI TravisCI `json:"travisCI"`
		Jenkins  Jenkins  `json:"jenkins"`
	}

	Ping struct {
		Count    int `json:"count"`
		Timeout  int `json:"timeout"`  // In Millisecond
		Interval int `json:"interval"` // In Millisecond
	}

	Port struct {
		Timeout int `json:"timeout"` // In Millisecond
	}

	Gitlab struct {
		Token string `json:"token"`
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

// Load confiuration from configuration file / env / default value
func InitConfig() *Config {
	var config Config

	// Setup Env
	viper.AutomaticEnv()
	viper.SetEnvPrefix(EnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

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

	// --- Gitlab Configuration ---
	viper.SetDefault("Monitorable.Gitlab.Token", "")

	// --- Github Configuration ---
	viper.SetDefault("Monitorable.Github.Token", "")

	// --- TravisCI Configuration ---
	viper.SetDefault("Monitorable.TravisCI.Url", "https://api.travis-ci.org/")
	viper.SetDefault("Monitorable.TravisCI.Timeout", 2000)
	viper.SetDefault("Monitorable.TravisCI.Token", "")

	// --- Jenkins Configuration ---
	viper.SetDefault("Monitorable.Jenkins.Timeout", 2000)
	viper.SetDefault("Monitorable.Jenkins.SSLVerify", true)
	viper.SetDefault("Monitorable.Jenkins.Url", "")
	viper.SetDefault("Monitorable.Jenkins.Login", "")
	viper.SetDefault("Monitorable.Jenkins.Token", "")

	_ = viper.Unmarshal(&config)

	return &config
}

func (t *TravisCI) IsValid() bool {
	return t.Url != ""
}
