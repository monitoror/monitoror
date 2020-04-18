package config

type (
	Github struct {
		URL                  string `validate:"required,url,http"`
		Token                string `validate:"required"`
		Timeout              int    `validate:"gte=0"` // In Millisecond
		CountCacheExpiration int    `validate:"gte=0"` // In Millisecond
	}
)

var Default = &Github{
	URL:                  "https://api.github.com/",
	Token:                "",
	Timeout:              5000,
	CountCacheExpiration: 30000,
}
