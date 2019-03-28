package repository

import (
	"errors"
	"time"

	"github.com/jsdidierlaurent/monitowall/config"
	"github.com/jsdidierlaurent/monitowall/monitorable/ping"
	"github.com/jsdidierlaurent/monitowall/monitorable/ping/model"

	goPing "github.com/sparrc/go-ping"
)

type (
	systemPingRepository struct {
		config *config.Config
	}
)

func NewNetworkPingRepository(config *config.Config) ping.Repository {
	return &systemPingRepository{config}
}

func (r *systemPingRepository) Ping(hostname string) (*model.Ping, error) {
	pinger, err := goPing.NewPinger(hostname)
	if err != nil {
		return nil, err
	}

	pinger.Count = r.config.Ping.Count
	pinger.Interval = time.Millisecond * time.Duration(r.config.Ping.Timeout)
	pinger.Timeout = time.Millisecond * time.Duration(r.config.Ping.Interval)
	pinger.SetPrivileged(true) // NEED ROOT PRIVILEGED

	pinger.Run()
	stats := pinger.Statistics()

	if stats.PacketsRecv == 0 {
		return nil, errors.New("ping failed")
	}

	ping := &model.Ping{}
	ping.Min = stats.MinRtt
	ping.Max = stats.MaxRtt
	ping.Average = stats.AvgRtt

	return ping, nil
}
