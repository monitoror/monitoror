package model

import (
	"time"

	goPing "github.com/sparrc/go-ping"
)

type (
	PingModelImpl interface {
		Ping(host string) Ping
	}

	Ping struct {
		Status  string `json:"status"`
		Label   string `json:"label"`
		Message string `json:"message,omitempty"`
	}

	PingModel struct{}
)

func NewPingModel() *PingModel {
	return &PingModel{}
}

func (u *PingModel) Ping(hostname string) Ping {
	pinger, err := goPing.NewPinger(hostname)
	if err != nil {
		// Lookup failed
		return Ping{Status: "FAILURE", Label: hostname}
	}

	pinger.Count = 2
	pinger.Interval = time.Second
	pinger.Timeout = time.Second * 3
	pinger.SetPrivileged(true) // NEED ROOT PRIVILEGED

	pinger.Run()
	stats := pinger.Statistics()

	var result Ping
	if stats.PacketsRecv == 0 {
		result = Ping{Status: "FAILURE", Label: hostname}
	} else {
		result = Ping{Status: "SUCCESS", Label: hostname, Message: stats.AvgRtt.String()}
	}

	return result
}
