package usecase

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/monitorable/ping"
	_pingModels "github.com/monitoror/monitoror/monitorable/ping/models"
	"github.com/monitoror/monitoror/monitorable/port"
	_portModels "github.com/monitoror/monitoror/monitorable/port/models"

	"github.com/monitoror/monitoror/monitorable/config"

	"github.com/monitoror/monitoror/monitorable/config/repository"

	"github.com/stretchr/testify/assert"
)

func initHydrateUsecase() config.Usecase {
	usecase := &configUsecase{
		monitorableConfigs: make(map[tiles.TileType]*MonitorableConfig),
	}

	usecase.Register(ping.PingTileType, "/ping", &_pingModels.PingParams{})
	usecase.Register(port.PortTileType, "/port", &_portModels.PortParams{})

	return usecase
}

func TestUsecase_Hydrate(t *testing.T) {
	usecase := initHydrateUsecase()
	reader := ioutil.NopCloser(strings.NewReader("{}"))
	config, err := repository.GetConfig(reader)
	assert.NoError(t, err)

	//TODO

	err = usecase.Hydrate(config)
	assert.NoError(t, err)
}
