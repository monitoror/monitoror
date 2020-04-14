package config

type (
	Ping struct {
		Count    int
		Timeout  int // In Millisecond
		Interval int // In Millisecond
	}
)

var Default = &Ping{
	Count:    2,
	Timeout:  1000,
	Interval: 100,
}
