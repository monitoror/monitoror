package config

type (
	Pingdom struct {
		URL             string `validate:"required,url,http"`
		Token           string `validate:"required"`
		Timeout         int    `validate:"gte=0"` // In Millisecond
		CacheExpiration int    `validate:"gte=0"` // In Millisecond
	}
)

var Default = &Pingdom{
	URL:             "https://api.pingdom.com/api/3.1",
	Token:           "",
	Timeout:         2000,
	CacheExpiration: 30000,
}
