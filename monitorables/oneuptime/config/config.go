package config

type (
	Oneuptime struct {
		URL             string `validate:"required,url,http"`
		ProjectID       string `validate:"required"`
		Token           string `validate:"required"`
		Timeout         int    `validate:"gte=0"` // In Millisecond
		CacheExpiration int    `validate:"gte=0"` // In Millisecond
		SSLVerify       bool
	}
)

var Default = &Oneuptime{
	URL:             "https://oneuptime.com",
	ProjectID:       "",
	Token:           "",
	Timeout:         10000,
	CacheExpiration: 30000,
	SSLVerify:       true,
}
