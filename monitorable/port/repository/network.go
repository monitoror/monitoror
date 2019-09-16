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
	portRepository struct {
		config *config.Port
		dialer pkgNet.Dialer
	}
)

func NewPortRepository(conf *config.Port) port.Repository {
	timeout := time.Millisecond * time.Duration(conf.Timeout)
	return &portRepository{conf, &net.Dialer{Timeout: timeout}}
}

func (r *portRepository) OpenSocket(hostname string, port int) (err error) {
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
