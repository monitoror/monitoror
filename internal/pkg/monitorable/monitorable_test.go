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

func TestValidateConfig(t *testing.T) {
	pingConf := &pingConfig.Ping{}
	errors := ValidateConfig(pingConf, "test")
	assert.NotEmpty(t, errors)
	assert.Equal(t, `Invalid "MO_MONITORABLE_PING_TEST_COUNT" field. Must be greater or equal to 1.`, errors[0].Error())

	pingConf.Count = 2
	errors = ValidateConfig(pingConf, "test")
	assert.Empty(t, errors)
}

func TestGetVariantsNames(t *testing.T) {
	conf := map[coreModels.VariantName]pingConfig.Ping{
		"test": {},
	}

	variants := GetVariantsNames(conf)
	assert.Len(t, variants, 1)
	assert.Equal(t, coreModels.VariantName("test"), variants[0])
	assert.Panics(t, func() { GetVariantsNames("test") })
}

func TestBuildMonitorableEnvKey(t *testing.T) {
	assert.Equal(t, "MO_MONITORABLE_PING_TEST_TEST", buildMonitorableEnvKey(pingConfig.Default, "test", "test"))
	assert.Equal(t, "MO_MONITORABLE_PING_TEST", buildMonitorableEnvKey(pingConfig.Default, coreModels.DefaultVariantName, "test"))
	assert.Panics(t, func() {
		buildMonitorableEnvKey("test", "test", "test")
	})
}
