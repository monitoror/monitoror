package config

type (
	Jenkins struct {
		URL       string
		Timeout   int // In Millisecond
		SSLVerify bool
		Login     string
		Token     string
	}
)

var Default = &Jenkins{
	URL:       "",
	Timeout:   2000,
	SSLVerify: true,
	Login:     "",
	Token:     "",
}
