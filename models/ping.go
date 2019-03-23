package models

import (
	"time"

	r "github.com/jsdidierlaurent/monitowall/renderings"

	goPing "github.com/sparrc/go-ping"
)

type (
	PingModelImpl interface {
		Ping(host string) *r.HealthCheckResponse
	}

	PingModel struct{}
)

func NewPingModel() *PingModel {
	return &PingModel{}
}

func newResponse() *r.HealthCheckResponse {
	return &r.HealthCheckResponse{
		Type: r.TypePing,
	}
}

func (u *PingModel) Ping(hostname string) (pingResponse *r.HealthCheckResponse) {
	pingResponse = newResponse()
	pingResponse.Label = hostname

	pinger, err := goPing.NewPinger(hostname)
	if err != nil {
		// Lookup failed
		pingResponse.Status = r.FailStatus
		return
	}

	pinger.Count = 2
	pinger.Interval = time.Second
	pinger.Timeout = time.Second * 3
	pinger.SetPrivileged(true) // NEED ROOT PRIVILEGED

	pinger.Run()
	stats := pinger.Statistics()

	if stats.PacketsRecv == 0 {
		pingResponse.Status = r.FailStatus
	} else {
		pingResponse.Status = r.SuccessStatus
		pingResponse.Message = stats.AvgRtt.String()
	}

	return
}
