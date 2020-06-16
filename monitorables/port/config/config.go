package config

type (
	Port struct {
		Timeout int `validate:"gte=0"` // In Millisecond
	}
)

var Default = &Port{
	Timeout: 2000,
}
