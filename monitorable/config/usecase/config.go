package usecase

import (
	. "github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/monitorable/config"
	"github.com/monitoror/monitoror/monitorable/config/models"
	. "github.com/monitoror/monitoror/pkg/monitoror/validator"
)

// Versions
const (
	CurrentVersion = Version2

	Version1 = 1
	Version2 = 2
)

const (
	// Custom Tile Type
	EmptyTileType   tiles.TileType = "EMPTY"
	GroupTileType   tiles.TileType = "GROUP"
	DynamicTileType tiles.TileType = "DYNAMIC"
)

var SupportedVersions = map[int]bool{
	Version1: true,
	Version2: true,
}

type (
	configUsecase struct {
		repository  config.Repository
		tileConfigs map[tiles.TileType]map[string]*TileConfig
	}

	// TileConfig struct is used by GetConfig endpoint to check / hydrate config
	TileConfig struct {
		Path      string
		Validator Validator
	}
)

func NewConfigUsecase(repository config.Repository) config.Usecase {
	return &configUsecase{
		repository:  repository,
		tileConfigs: make(map[tiles.TileType]map[string]*TileConfig),
	}
}

func (cu *configUsecase) RegisterTile(tileType tiles.TileType, path string, validator Validator) {
	cu.RegisterTileWithConfigVariant(tileType, DefaultVariant, path, validator)
}

func (cu *configUsecase) RegisterTileWithConfigVariant(tileType tiles.TileType, variant, path string, validator Validator) {
	value, exists := cu.tileConfigs[tileType]
	if !exists {
		value = make(map[string]*TileConfig)
		cu.tileConfigs[tileType] = value
	}

	value[variant] = &TileConfig{
		Path:      path,
		Validator: validator,
	}
}

// GetConfig load and parse GetConfig
func (cu *configUsecase) GetConfig(params *models.ConfigParams) (config *models.Config, err error) {
	if params.Url != "" {
		config, err = cu.repository.GetConfigFromUrl(params.Url)
	} else if params.Path != "" {
		config, err = cu.repository.GetConfigFromPath(params.Path)
	}

	if err != nil {
		return
	}

	// Set config to CurrentVersion if config isn't set
	if config.Version == 0 {
		config.Version = CurrentVersion
	}

	return
}
