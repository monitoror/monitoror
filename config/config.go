package config

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/viper"
)

const EnvPrefix = "MO"
const DefaultVariant = "default"

type (
	// Config contain backend Configuration
	Config struct {
		// --- General Configuration ---
		Port int    `json:"port"` // Default: 8080
		Env  string `json:"env"`  // Default: production

		// --- Cache Configuration ---
		// UpstreamCacheExpiration is used to respond before executing the request. Avoid overloading services.
		UpstreamCacheExpiration int `json:"upstreamCacheExpiration"`
		// DownstreamCacheExpiration is used to respond after executing the request in case of timeout error.
		DownstreamCacheExpiration int `json:"downstreamCacheExpiration"`

		// Monitorable Config
		Monitorable Monitorable `json:"monitorable"`
	}

	Monitorable struct {
		Ping        Ping                    `json:"ping"`
		Port        Port                    `json:"port"`
		HTTP        HTTP                    `json:"http"`
		Pingdom     map[string]*Pingdom     `json:"pingdom"`     // With variants
		Github      map[string]*Github      `json:"github"`      // With variants
		TravisCI    map[string]*TravisCI    `json:"travisCI"`    // With variants
		Jenkins     map[string]*Jenkins     `json:"jenkins"`     // With variants
		AzureDevOps map[string]*AzureDevOps `json:"azureDevOps"` // With variants
	}

	Ping struct {
		Count    int `json:"count"`
		Timeout  int `json:"timeout"`  // In Millisecond
		Interval int `json:"interval"` // In Millisecond
	}

	Port struct {
		Timeout int `json:"timeout"` // In Millisecond
	}

	HTTP struct {
		Timeout   int  `json:"timeout"` // In Millisecond
		SSLVerify bool `json:"sslVerify"`
	}

	Pingdom struct {
		URL             string `json:"url"`
		Token           string `json:"token"`
		Timeout         int    `json:"timeout"`         // In Millisecond
		CacheExpiration int    `json:"cacheExpiration"` // In Millisecond
	}

	Github struct {
		Token string `json:"token"`
	}

	TravisCI struct {
		URL     string `json:"url"`
		Timeout int    `json:"timeout"` // In Millisecond
		Token   string `json:"token"`
	}

	Jenkins struct {
		URL       string `json:"url"`
		Timeout   int    `json:"timeout"` // In Millisecond
		SSLVerify bool   `json:"sslVerify"`
		Login     string `json:"login"`
		Token     string `json:"token"`
	}

	AzureDevOps struct {
		URL     string `json:"url"`
		Timeout int    `json:"timeout"` // In Millisecond
		Token   string `json:"token"`
	}
)

// InitConfig from configuration file / env / default value
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
	viper.SetDefault("UpstreamCacheExpiration", 10000)
	viper.SetDefault("DownstreamCacheExpiration", 120000)

	// --- Ping Configuration ---
	viper.SetDefault("Monitorable.Ping.Count", 2)
	viper.SetDefault("Monitorable.Ping.Timeout", 1000)
	viper.SetDefault("Monitorable.Ping.Interval", 100)

	// --- Port Configuration ---
	viper.SetDefault("Monitorable.Port.Timeout", 2000)

	// --- HTTP Configuration ---
	viper.SetDefault("Monitorable.HTTP.Timeout", 2000)
	viper.SetDefault("Monitorable.HTTP.SSLVerify", true)
	viper.SetDefault("Monitorable.HTTP.URL", "")

	// --- Pingdom Configuration ---
	for variant := range variants["Pingdom"] {
		viper.SetDefault(fmt.Sprintf("Monitorable.Pingdom.%s.URL", variant), "")
		viper.SetDefault(fmt.Sprintf("Monitorable.Pingdom.%s.Token", variant), "")
		viper.SetDefault(fmt.Sprintf("Monitorable.Pingdom.%s.Timeout", variant), 2000)
		viper.SetDefault(fmt.Sprintf("Monitorable.Pingdom.%s.CacheExpiration", variant), 30000)
	}

	// --- Github Configuration ---
	for variant := range variants["Github"] {
		viper.SetDefault(fmt.Sprintf("Monitorable.Github.%s.Token", variant), "")
	}

	// --- TravisCI Configuration ---
	for variant := range variants["TravisCI"] {
		viper.SetDefault(fmt.Sprintf("Monitorable.TravisCI.%s.URL", variant), "https://api.travis-ci.org/")
		viper.SetDefault(fmt.Sprintf("Monitorable.TravisCI.%s.Timeout", variant), 2000)
		viper.SetDefault(fmt.Sprintf("Monitorable.TravisCI.%s.Token", variant), "")
	}

	// --- Jenkins Configuration ---
	for variant := range variants["Jenkins"] {
		viper.SetDefault(fmt.Sprintf("Monitorable.Jenkins.%s.URL", variant), "")
		viper.SetDefault(fmt.Sprintf("Monitorable.Jenkins.%s.Timeout", variant), 2000)
		viper.SetDefault(fmt.Sprintf("Monitorable.Jenkins.%s.SSLVerify", variant), true)
		viper.SetDefault(fmt.Sprintf("Monitorable.Jenkins.%s.Login", variant), "")
		viper.SetDefault(fmt.Sprintf("Monitorable.Jenkins.%s.Token", variant), "")
	}

	// --- Azure DevOps Configuration ---
	for variant := range variants["AzureDevOps"] {
		viper.SetDefault(fmt.Sprintf("Monitorable.AzureDevOps.%s.URL", variant), "")
		viper.SetDefault(fmt.Sprintf("Monitorable.AzureDevOps.%s.Timeout", variant), 4000)
		viper.SetDefault(fmt.Sprintf("Monitorable.AzureDevOps.%s.Token", variant), "")
	}

	_ = viper.Unmarshal(&config)

	return &config
}

func (t *Pingdom) IsValid() bool {
	// Pingdom url can be empty, plugin will use default value
	if t.URL != "" {
		if _, err := url.Parse(t.URL); err != nil {
			return false
		}
	}

	return t.Token != ""
}

func (t *TravisCI) IsValid() bool {
	if t.URL == "" {
		return false
	}

	if _, err := url.Parse(t.URL); err != nil {
		return false
	}

	return true
}

func (t *Jenkins) IsValid() bool {
	if t.URL == "" {
		return false
	}

	if _, err := url.Parse(t.URL); err != nil {
		return false
	}

	return true
}

func (t *AzureDevOps) IsValid() bool {
	if t.URL == "" {
		return false
	}

	if _, err := url.Parse(t.URL); err != nil {
		return false
	}

	return t.Token != ""
}
