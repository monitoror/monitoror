package config

type (
	Youtrack struct {
		URL       string `validate:"required,url,http"`
		Token     string
		Timeout   int `validate:"gte=0"` // In Millisecond
		SSLVerify bool
	}
)

var Default = &Youtrack{
	URL:       "",
	Token:     "",
	Timeout:   5000,
	SSLVerify: true,
}
