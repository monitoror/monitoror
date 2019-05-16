package usecase

import (
	"github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/monitorable/config"
	"github.com/monitoror/monitoror/monitorable/config/models"
	"github.com/monitoror/monitoror/pkg/monitoror/utils"
)

type (
	configUsecase struct {
		monitorableParams map[tiles.TileType]utils.Validator
		repository        config.Repository
	}
)

func NewConfigUsecase(monitorableParams map[tiles.TileType]utils.Validator, repository config.Repository) config.Usecase {
	return &configUsecase{monitorableParams, repository}
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
