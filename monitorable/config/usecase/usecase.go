package usecase

import (
	"time"

	"github.com/jsdidierlaurent/echo-middleware/cache"

	"github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/models"
	monitorableConfig "github.com/monitoror/monitoror/monitorable/config"
	"github.com/monitoror/monitoror/pkg/monitoror/builder"
	"github.com/monitoror/monitoror/pkg/monitoror/utils"
)

// Versions
const (
	CurrentVersion = Version5

	Version1 = 1
	Version2 = 2
	Version3 = 3
	Version4 = 4
	Version5 = 5
)

const (
	EmptyTileType models.TileType = "EMPTY"
	GroupTileType models.TileType = "GROUP"

	DynamicTileStoreKeyPrefix = "monitoror.config.dynamicTile.key"
)

var SupportedVersions = map[int]bool{
	Version5: true,
}

type (
	configUsecase struct {
		repository monitorableConfig.Repository

		tileConfigs        map[models.TileType]map[string]*TileConfig
		dynamicTileConfigs map[models.TileType]map[string]*DynamicTileConfig

		// jobs cache. used in case of timeout
		dynamicTileStore          cache.Store
		downstreamStoreExpiration time.Duration
	}

	// TileConfig struct is used by GetConfig endpoint to check / hydrate config
	TileConfig struct {
		Validator utils.Validator
		Path      string
	}

	// DynamicTileConfig struct is used by GetConfig endpoint to check / hydrate config
	DynamicTileConfig struct {
		Validator utils.Validator
		Builder   builder.DynamicTileBuilder
	}
)

func NewConfigUsecase(repository monitorableConfig.Repository, store cache.Store, downstreamStoreExpiration int) monitorableConfig.Usecase {
	tileConfigs := make(map[models.TileType]map[string]*TileConfig)

	// Used for authorized type
	tileConfigs[EmptyTileType] = nil
	tileConfigs[GroupTileType] = nil

	dynamicTileConfigs := make(map[models.TileType]map[string]*DynamicTileConfig)

	return &configUsecase{
		repository:                repository,
		tileConfigs:               tileConfigs,
		dynamicTileConfigs:        dynamicTileConfigs,
		dynamicTileStore:          store,
		downstreamStoreExpiration: time.Millisecond * time.Duration(downstreamStoreExpiration),
	}
}

func (cu *configUsecase) RegisterTile(tileType models.TileType, clientConfigValidator utils.Validator, path string) {
	cu.RegisterTileWithConfigVariant(tileType, config.DefaultVariant, clientConfigValidator, path)
}

func (cu *configUsecase) RegisterTileWithConfigVariant(tileType models.TileType, variant string, clientConfigValidator utils.Validator, path string) {
	value, exists := cu.tileConfigs[tileType]
	if !exists {
		value = make(map[string]*TileConfig)
		cu.tileConfigs[tileType] = value
	}

	value[variant] = &TileConfig{
		Path:      path,
		Validator: clientConfigValidator,
	}
}

func (cu *configUsecase) RegisterDynamicTile(tileType models.TileType, clientConfigValidator utils.Validator, builder builder.DynamicTileBuilder) {
	cu.RegisterDynamicTileWithConfigVariant(tileType, config.DefaultVariant, clientConfigValidator, builder)
}

func (cu *configUsecase) RegisterDynamicTileWithConfigVariant(tileType models.TileType, variant string, clientConfigValidator utils.Validator, builder builder.DynamicTileBuilder) {
	// Used for authorized type
	cu.tileConfigs[tileType] = nil

	value, exists := cu.dynamicTileConfigs[tileType]
	if !exists {
		value = make(map[string]*DynamicTileConfig)
	}

	value[variant] = &DynamicTileConfig{
		Validator: clientConfigValidator,
		Builder:   builder,
	}
	cu.dynamicTileConfigs[tileType] = value
}
