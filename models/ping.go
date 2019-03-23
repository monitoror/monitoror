package models

import (
	"time"

	. "github.com/jsdidierlaurent/monitowall/renderings"

	goPing "github.com/sparrc/go-ping"
)

type (
	PingModelImpl interface {
		Ping(host string) *HealthCheckResponse
	}

	PingModel struct{}
)

func NewPingModel() *PingModel {
	return &PingModel{}
}

func newResponse() *HealthCheckResponse {
	return &HealthCheckResponse{
		Type: TypePing,
	}
}

func (u *PingModel) Ping(hostname string) (pingResponse *HealthCheckResponse) {
	pingResponse = newResponse()
	pingResponse.Label = hostname

	pinger, err := goPing.NewPinger(hostname)
	if err != nil {
		// Lookup failed
		pingResponse.Status = FailStatus
		return
	}

	pinger.Count = 2
	pinger.Interval = time.Second
	pinger.Timeout = time.Second * 3
	pinger.SetPrivileged(true) // NEED ROOT PRIVILEGED

	pinger.Run()
	stats := pinger.Statistics()

	if stats.PacketsRecv == 0 {
		pingResponse.Status = FailStatus
	} else {
		pingResponse.Status = SuccessStatus
		pingResponse.Message = stats.AvgRtt.String()
	}

	return
}
