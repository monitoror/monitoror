package config

type (
	HTTP struct {
		Timeout     int `validate:"gte=0"` // In Millisecond
		SSLVerify   bool
		Certificate string
		Key         string
		Headers     []string
	}
)

var Default = &HTTP{
	Timeout:   2000,
	SSLVerify: true,
	Headers:   []string{"User-Agent: Monitoror"},
}
