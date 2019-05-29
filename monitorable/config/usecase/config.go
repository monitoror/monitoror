package usecase

import (
	"github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/monitorable/config"
	"github.com/monitoror/monitoror/monitorable/config/models"
	"github.com/monitoror/monitoror/pkg/monitoror/utils"
)

type (
	configUsecase struct {
		monitorableConfigs map[tiles.TileType]*MonitorableConfig
		repository         config.Repository
	}

	MonitorableConfig struct {
		Path            string
		ConfigValidator utils.Validator
	}
)

func NewConfigUsecase(repository config.Repository) config.Usecase {
	return &configUsecase{
		monitorableConfigs: make(map[tiles.TileType]*MonitorableConfig),
		repository:         repository,
	}
}

func (cu *configUsecase) Register(tileType tiles.TileType, path string, configValidator utils.Validator) {
	cu.monitorableConfigs[tileType] = &MonitorableConfig{
		Path:            path,
		ConfigValidator: configValidator,
	}
}

//Config load and parse Config
func (cu *configUsecase) Config(params *models.ConfigParams) (config *models.Config, err error) {
	if params.Url != "" {
		config, err = cu.repository.GetConfigFromUrl(params.Url)
	} else if params.Path != "" {
		config, err = cu.repository.GetConfigFromPath(params.Path)
	}

	return
}
