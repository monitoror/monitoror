package config

type (
	Ping struct {
		Count    int `validate:"gte=1"`
		Timeout  int `validate:"gte=0"` // In Millisecond
		Interval int `validate:"gte=0"` // In Millisecond
	}
)

var Default = &Ping{
	Count:    2,
	Timeout:  1000,
	Interval: 100,
}
