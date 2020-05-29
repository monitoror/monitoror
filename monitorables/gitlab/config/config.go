package config

type (
	Gitlab struct {
		URL     string `validate:"required,url,http"`
		Token   string `validate:"required"`
		Timeout int    `validate:"gte=0"` // In Millisecond
	}
)

var Default = &Gitlab{
	URL:     "https://gitlab.com/",
	Token:   "",
	Timeout: 5000,
}
