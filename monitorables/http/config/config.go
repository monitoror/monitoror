package config

type (
	HTTP struct {
		Timeout   int `validate:"gte=0"` // In Millisecond
		SSLVerify bool
	}
)

var Default = &HTTP{
	Timeout:   2000,
	SSLVerify: true,
}
