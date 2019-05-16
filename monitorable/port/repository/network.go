package repository

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/monitorable/port"
	pkgNet "github.com/monitoror/monitoror/pkg/net"
)

type (
	portRepository struct {
		config *config.Config
		dialer pkgNet.Dialer
	}
)

func NewPortRepository(conf *config.Config) port.Repository {
	timeout := time.Millisecond * time.Duration(conf.Monitorable.Port.Timeout)
	return &portRepository{conf, &net.Dialer{Timeout: timeout}}
}

func (r *portRepository) OpenSocket(ctx context.Context, hostname string, port int) (err error) {
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
