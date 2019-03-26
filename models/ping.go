package models

import (
	"fmt"
	"time"

	goPing "github.com/sparrc/go-ping"
)

type (
	PingModel interface {
		Ping(host string) (response string, err error)
	}

	PingModelImpl struct{}
)

func NewPingModel() *PingModelImpl {
	return &PingModelImpl{}
}

func (u *PingModelImpl) Ping(hostname string) (response string, err error) {
	pinger, err := goPing.NewPinger(hostname)
	if err != nil {
		return
	}

	pinger.Count = 2
	pinger.Interval = time.Second
	pinger.Timeout = time.Second * 3
	pinger.SetPrivileged(true) // NEED ROOT PRIVILEGED

	pinger.Run()
	stats := pinger.Statistics()

	if stats.PacketsRecv == 0 {
		err = fmt.Errorf("ping failed")
	} else {
		response = stats.AvgRtt.String()
	}
	return
}
