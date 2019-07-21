package repository

import (
	"errors"
	"time"

	"github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/monitorable/ping"
	"github.com/monitoror/monitoror/monitorable/ping/models"

	goPing "github.com/sparrc/go-ping"
)

type (
	pingRepository struct {
		config *config.Config
	}
)

func NewPingRepository(config *config.Config) ping.Repository {
	return &pingRepository{config}
}

func (r *pingRepository) ExecutePing(hostname string) (*models.Ping, error) {
	pinger, err := goPing.NewPinger(hostname)
	if err != nil {
		return nil, err
	}

	pinger.Count = r.config.Monitorable.Ping.Count
	pinger.Interval = time.Millisecond * time.Duration(r.config.Monitorable.Ping.Interval)
	pinger.Timeout = time.Millisecond * time.Duration(r.config.Monitorable.Ping.Timeout)
	pinger.SetPrivileged(true) // NEED ROOT PRIVILEGED

	pinger.Run()
	stats := pinger.Statistics()

	if stats.PacketsRecv == 0 {
		return nil, errors.New("ping failed")
	}

	ping := &models.Ping{}
	ping.Min = stats.MinRtt
	ping.Max = stats.MaxRtt
	ping.Average = stats.AvgRtt

	return ping, nil
}
