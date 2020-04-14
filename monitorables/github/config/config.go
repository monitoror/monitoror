package config

type (
	Github struct {
		URL                  string
		Timeout              int // In Millisecond
		Token                string
		CountCacheExpiration int // In Millisecond
	}
)

var Default = &Github{
	URL:                  "https://api.github.com/",
	Timeout:              5000,
	Token:                "",
	CountCacheExpiration: 30000,
}
