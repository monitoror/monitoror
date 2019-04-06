package repository

import (
	"fmt"
	"net"
	"time"

	"github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/monitorable/port"
	pkgNet "github.com/monitoror/monitoror/pkg/net"
)

type (
	systemPortRepository struct {
		config *config.Config
		dialer pkgNet.Dialer
	}
)

func NewNetworkPortRepository(conf *config.Config) port.Repository {
	timeout := time.Millisecond * time.Duration(conf.PortConfig.Timeout)
	return &systemPortRepository{conf, &net.Dialer{Timeout: timeout}}
}

func (r *systemPortRepository) CheckPort(hostname string, port int) (err error) {
	target := fmt.Sprintf("%s:%d", hostname, port)

	conn, err := r.dialer.Dial("tcp", target)
	if err != nil {
		return
	}
	if conn != nil {
		_ = conn.Close()
	}

	return nil
}
