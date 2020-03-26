package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"

	pingConfig "github.com/monitoror/monitoror/monitorables/ping/config"
	"github.com/monitoror/monitoror/pkg/system"
)

// /!\ this is an integration test /!\
// Note : It may be necessary to separate them from unit tests

func TestRepository_Ping_Error(t *testing.T) {
	pingRepository := NewPingRepository(pingConfig.Default)

	ping, err := pingRepository.ExecutePing("0.0.0.0")

	// I can't mock ping library, so i just test this repository
	if system.IsRawSocketAvailable() {
		assert.NoError(t, err)
		assert.NotNil(t, ping)
	} else {
		assert.Error(t, err)
		assert.Nil(t, ping)
	}
}
