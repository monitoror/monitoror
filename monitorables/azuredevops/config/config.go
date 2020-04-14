package config

type (
	AzureDevOps struct {
		URL     string
		Timeout int // In Millisecond
		Token   string
	}
)

var Default = &AzureDevOps{
	URL:     "",
	Timeout: 4000,
	Token:   "",
}
