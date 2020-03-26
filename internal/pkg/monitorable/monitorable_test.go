package monitorable

import (
	"os"
	"testing"

	coreModels "github.com/monitoror/monitoror/models"
	pingConfig "github.com/monitoror/monitoror/monitorables/ping/config"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	_ = os.Setenv("MO_MONITORABLE_PING_TIMEOUT", "1337")
	_ = os.Setenv("MO_MONITORABLE_PING_TEST_TIMEOUT", "1337")

	conf := make(map[coreModels.VariantName]pingConfig.Ping)
	LoadConfig(&conf, pingConfig.Default)

	assert.Len(t, conf, 2)
	assert.Equal(t, 1337, conf["default"].Timeout)
	assert.Equal(t, 1337, conf["test"].Timeout)

}

func TestGetVariants(t *testing.T) {
	conf := map[coreModels.VariantName]pingConfig.Ping{
		"test": {},
	}

	variants := GetVariants(conf)
	assert.Len(t, variants, 1)
	assert.Equal(t, coreModels.VariantName("test"), variants[0])
	assert.Panics(t, func() { GetVariants("test") })
}

func TestGetEnvName(t *testing.T) {
	assert.Equal(t, "MO_MONITORABLE_PING_TEST_TEST", GetEnvName(pingConfig.Default, "test", "test"))
	assert.Equal(t, "MO_MONITORABLE_PING_TEST", GetEnvName(pingConfig.Default, coreModels.DefaultVariant, "test"))
	assert.Panics(t, func() {
		GetEnvName("test", "test", "test")
	})
}
