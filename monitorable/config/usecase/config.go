package usecase

import (
	"github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/monitorable/config"
	"github.com/monitoror/monitoror/monitorable/config/models"
	"github.com/monitoror/monitoror/pkg/monitoror/utils"
)

const (
	TypeKey   = "type"
	LabelKey  = "label"
	ParamsKey = "params"
	TilesKey  = "tiles"
	UrlKey    = "url" // Inject by hydrate function

	EmptyTileType tiles.TileType = "EMPTY"
	GroupTileType tiles.TileType = "GROUP"
)

type (
	configUsecase struct {
		repository  config.Repository
		tileConfigs map[tiles.TileType]*TileConfig
	}

	// TileConfig struct is used by Config endpoint to check / hydrate config
	TileConfig struct {
		Path      string
		Validator utils.Validator
	}
)

func NewConfigUsecase(repository config.Repository) config.Usecase {
	return &configUsecase{
		repository:  repository,
		tileConfigs: make(map[tiles.TileType]*TileConfig),
	}
}

func (cu *configUsecase) RegisterTile(tileType tiles.TileType, path string, validator utils.Validator) {
	cu.tileConfigs[tileType] = &TileConfig{
		Path:      path,
		Validator: validator,
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
