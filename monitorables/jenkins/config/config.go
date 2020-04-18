package config

type (
	Jenkins struct {
		URL       string `validate:"required,url,http"`
		Login     string
		Token     string
		Timeout   int `validate:"gte=0"` // In Millisecond
		SSLVerify bool
	}
)

var Default = &Jenkins{
	URL:       "",
	Login:     "",
	Token:     "",
	Timeout:   2000,
	SSLVerify: true,
}
