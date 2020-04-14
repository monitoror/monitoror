package config

type (
	Port struct {
		Timeout int // In Millisecond
	}
)

var Default = &Port{
	Timeout: 2000,
}
