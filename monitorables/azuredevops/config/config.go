package config

type (
	AzureDevOps struct {
		URL     string `validate:"required,url,http"`
		Token   string `validate:"required"`
		Timeout int    `validate:"gte=0"` // In Millisecond
	}
)

var Default = &AzureDevOps{
	URL:     "",
	Token:   "",
	Timeout: 4000,
}
