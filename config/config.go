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
		Port int    // Default: 8080
		Env  string // Default: production

		// --- Cache Configuration ---
		// UpstreamCacheExpiration is used to respond before executing the request. Avoid overloading services.
		UpstreamCacheExpiration int
		// DownstreamCacheExpiration is used to respond after executing the request in case of timeout error.
		DownstreamCacheExpiration int

		// Monitorable Config
		Monitorable Monitorable
	}

	Monitorable struct {
		Ping        map[string]*Ping
		Port        map[string]*Port
		HTTP        map[string]*HTTP
		Pingdom     map[string]*Pingdom
		TravisCI    map[string]*TravisCI
		Jenkins     map[string]*Jenkins
		AzureDevOps map[string]*AzureDevOps
		Github      map[string]*Github
	}

	Ping struct {
		Count    int
		Timeout  int // In Millisecond
		Interval int // In Millisecond
	}

	Port struct {
		Timeout int // In Millisecond
	}

	HTTP struct {
		Timeout   int // In Millisecond
		SSLVerify bool
	}

	Pingdom struct {
		URL             string
		Token           string
		Timeout         int // In Millisecond
		CacheExpiration int // In Millisecond
	}

	TravisCI struct {
		URL         string
		Timeout     int // In Millisecond
		Token       string
		GithubToken string
	}

	Jenkins struct {
		URL       string
		Timeout   int // In Millisecond
		SSLVerify bool
		Login     string
		Token     string
	}

	AzureDevOps struct {
		URL     string
		Timeout int // In Millisecond
		Token   string
	}

	Github struct {
		Timeout              int // In Millisecond
		Token                string
		IssueCacheExpiration int // In Millisecond
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
	for variant := range variants["Ping"] {
		viper.SetDefault(fmt.Sprintf("Monitorable.Ping.%s.Count", variant), 2)
		viper.SetDefault(fmt.Sprintf("Monitorable.Ping.%s.Timeout", variant), 1000)
		viper.SetDefault(fmt.Sprintf("Monitorable.Ping.%s.Interval", variant), 100)
	}

	// --- Port Configuration ---
	for variant := range variants["Port"] {
		viper.SetDefault(fmt.Sprintf("Monitorable.Port.%s.Timeout", variant), 2000)
	}

	// --- HTTP Configuration ---
	for variant := range variants["HTTP"] {
		viper.SetDefault(fmt.Sprintf("Monitorable.HTTP.%s.Timeout", variant), 2000)
		viper.SetDefault(fmt.Sprintf("Monitorable.HTTP.%s.SSLVerify", variant), true)
		viper.SetDefault(fmt.Sprintf("Monitorable.HTTP.%s.URL", variant), "")
	}

	// --- Pingdom Configuration ---
	for variant := range variants["Pingdom"] {
		viper.SetDefault(fmt.Sprintf("Monitorable.Pingdom.%s.URL", variant), "https://api.pingdom.com/api/3.1")
		viper.SetDefault(fmt.Sprintf("Monitorable.Pingdom.%s.Token", variant), "")
		viper.SetDefault(fmt.Sprintf("Monitorable.Pingdom.%s.Timeout", variant), 2000)
		viper.SetDefault(fmt.Sprintf("Monitorable.Pingdom.%s.CacheExpiration", variant), 30000)
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

	// --- Github Configuration ---
	for variant := range variants["Github"] {
		viper.SetDefault(fmt.Sprintf("Monitorable.Github.%s.Timeout", variant), 5000)
		viper.SetDefault(fmt.Sprintf("Monitorable.Github.%s.Token", variant), "")
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

func (g *Github) IsValid() bool {
	return g.Token != ""
}
