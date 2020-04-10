package config

type (
	HTTP struct {
		Timeout   int // In Millisecond
		SSLVerify bool
	}
)

var Default = &HTTP{
	Timeout:   2000,
	SSLVerify: true,
}
