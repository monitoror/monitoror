package usecase

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/monitoror/monitoror/api/config"
	"github.com/monitoror/monitoror/api/config/models"
	"github.com/monitoror/monitoror/api/config/repository"
	coreModels "github.com/monitoror/monitoror/models"
	jenkinsApi "github.com/monitoror/monitoror/monitorables/jenkins/api"
	jenkinsModels "github.com/monitoror/monitoror/monitorables/jenkins/api/models"
	pingApi "github.com/monitoror/monitoror/monitorables/ping/api"
	pingModels "github.com/monitoror/monitoror/monitorables/ping/api/models"
	pingdomApi "github.com/monitoror/monitoror/monitorables/pingdom/api"
	pindomModels "github.com/monitoror/monitoror/monitorables/pingdom/api/models"
	portApi "github.com/monitoror/monitoror/monitorables/port/api"
	portModels "github.com/monitoror/monitoror/monitorables/port/api/models"

	"github.com/jsdidierlaurent/echo-middleware/cache"
	"github.com/stretchr/testify/assert"
)

func initConfigUsecase(repository config.Repository, store cache.Store) *configUsecase {
	usecase := NewConfigUsecase(repository, store, 5000)

	usecase.RegisterTile(pingApi.PingTileType, []coreModels.VariantName{coreModels.DefaultVariant}, MinimalVersion)
	usecase.RegisterTile(portApi.PortTileType, []coreModels.VariantName{coreModels.DefaultVariant}, MinimalVersion)
	usecase.RegisterTile(jenkinsApi.JenkinsBuildTileType, []coreModels.VariantName{coreModels.DefaultVariant}, MinimalVersion)
	usecase.RegisterTile(pingdomApi.PingdomCheckTileType, []coreModels.VariantName{coreModels.DefaultVariant}, MinimalVersion)

	usecase.EnableTile(pingApi.PingTileType, coreModels.DefaultVariant, &pingModels.PingParams{}, "/ping", 1000)
	usecase.EnableTile(portApi.PortTileType, coreModels.DefaultVariant, &portModels.PortParams{}, "/port", 1000)
	usecase.EnableTile(jenkinsApi.JenkinsBuildTileType, coreModels.DefaultVariant, &jenkinsModels.BuildParams{}, "/jenkins/default", 1000)
	usecase.EnableTile(pingdomApi.PingdomCheckTileType, coreModels.DefaultVariant, &pindomModels.CheckParams{}, "/pingdom/default", 1000)

	return usecase.(*configUsecase)
}

func readConfig(input string) (configBag *models.ConfigBag, err error) {
	configBag = &models.ConfigBag{}
	reader := ioutil.NopCloser(strings.NewReader(input))
	configBag.Config, err = repository.ReadConfig(reader)
	return
}

func TestUsecase_Global_Success(t *testing.T) {
	rawConfig := fmt.Sprintf(`
{
	"version" : %q,
  "columns": 4,
  "tiles": [
		{ "type": "PORT", "label": "Monitoror", "params": {"hostname": "localhost", "port": 8080} }
  ]
}
`, CurrentVersion)

	expectRawConfig := fmt.Sprintf(`{"config":{"version":%q,"columns":4,"tiles":[{"type":"PORT","label":"Monitoror","url":"/port?hostname=localhost\u0026port=8080","initialMaxDelay":1000}]}}`, CurrentVersion)

	config, err := readConfig(rawConfig)
	if assert.NoError(t, err) {
		usecase := initConfigUsecase(nil, nil)
		usecase.Verify(config)
		usecase.Hydrate(config)

		assert.Len(t, config.Errors, 0)

		marshal, err := json.Marshal(config)
		assert.NoError(t, err)
		assert.Equal(t, expectRawConfig, string(marshal))
	}
}
