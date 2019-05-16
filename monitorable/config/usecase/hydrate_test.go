package usecase

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/monitoror/monitoror/monitorable/config/repository"

	"github.com/stretchr/testify/assert"

	"github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/monitorable/ping"
	_pingModels "github.com/monitoror/monitoror/monitorable/ping/models"
	"github.com/monitoror/monitoror/monitorable/port"
	_portModels "github.com/monitoror/monitoror/monitorable/port/models"
	"github.com/monitoror/monitoror/pkg/monitoror/utils"
)

func initHydrateUsecase() *configUsecase {
	usecase := &configUsecase{}

	usecase.monitorableParams = make(map[tiles.TileType]utils.Validator)
	usecase.monitorableParams[ping.PingTileType] = &_pingModels.PingParams{}
	usecase.monitorableParams[port.PortTileType] = &_portModels.PortParams{}

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
