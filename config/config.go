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
		Port int    // Default: 8080
		Env  string // Default: production

		// --- Cache Configuration ---
		//UpstreamCache is used to respond before executing the request. Avoid overloading services.
		UpstreamCache Cache
		//DownstreamCache is used to respond after executing the request in case of timeout error.
		DownstreamCache Cache

		//Monitorables Config
		Monitorable Monitorable
	}

	Cache struct {
		Expire          int // In Millisecond
		CleanupInterval int // In Millisecond
	}

	Monitorable struct {
		Ping     Ping
		Port     Port
		Gitlab   Gitlab
		Github   Github
		TravisCI TravisCI
	}

	Ping struct {
		Count    int
		Timeout  int // In Millisecond
		Interval int // In Millisecond
	}

	Port struct {
		Timeout int // In Millisecond
	}

	Gitlab struct {
		Token string
	}

	Github struct {
		Token string
	}

	TravisCI struct {
		Token   string
		Timeout int // In Millisecond
		Url     string
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

	// --- TravisCI Configuration ---
	viper.SetDefault("Monitorable.TravisCI.Timeout", 2000)
	viper.SetDefault("Monitorable.TravisCI.Url", "https://api.travis-ci.org/")

	// Read Configuration
	_ = viper.Unmarshal(&config)

	return &config
}

func (t *TravisCI) IsValid() bool {
	return t.Url != ""
}
