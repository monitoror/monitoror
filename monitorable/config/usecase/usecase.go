package usecase

import (
	"time"

	"github.com/jsdidierlaurent/echo-middleware/cache"

	. "github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/config"
	. "github.com/monitoror/monitoror/pkg/monitoror/builder"
	. "github.com/monitoror/monitoror/pkg/monitoror/validator"
)

// Versions
const (
	CurrentVersion = Version3

	Version1 = 1
	Version2 = 2
	Version3 = 3
)

const (
	// Custom Tile Type
	EmptyTileType models.TileType = "EMPTY"
	GroupTileType models.TileType = "GROUP"

	// Store key prefix
	DynamicTileStoreKeyPrefix = "monitoror.config.dynamicTile.key"
)

var SupportedVersions = map[int]bool{
	Version3: true,
}

type (
	configUsecase struct {
		repository config.Repository

		tileConfigs        map[models.TileType]map[string]*TileConfig
		dynamicTileConfigs map[models.TileType]map[string]*DynamicTileConfig

		// jobs cache. used in case of timeout
		dynamicTileStore          cache.Store
		downstreamStoreExpiration time.Duration
	}

	// TileConfig struct is used by GetConfig endpoint to check / hydrate config
	TileConfig struct {
		Validator Validator
		Path      string
	}

	// DynamicTileConfig struct is used by GetConfig endpoint to check / hydrate config
	DynamicTileConfig struct {
		Validator Validator
		Builder   DynamicTileBuilder
	}
)

func NewConfigUsecase(repository config.Repository, store cache.Store, downstreamStoreExpiration int) config.Usecase {
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

func (cu *configUsecase) RegisterTile(tileType models.TileType, validator Validator, path string) {
	cu.RegisterTileWithConfigVariant(tileType, DefaultVariant, validator, path)
}

func (cu *configUsecase) RegisterTileWithConfigVariant(tileType models.TileType, variant string, validator Validator, path string) {
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

func (cu *configUsecase) RegisterDynamicTile(tileType models.TileType, validator Validator, builder DynamicTileBuilder) {
	cu.RegisterDynamicTileWithConfigVariant(tileType, DefaultVariant, validator, builder)
}

func (cu *configUsecase) RegisterDynamicTileWithConfigVariant(tileType models.TileType, variant string, validator Validator, builder DynamicTileBuilder) {
	// Used for authorized type
	cu.tileConfigs[tileType] = nil

	value, exists := cu.dynamicTileConfigs[tileType]
	if !exists {
		value = make(map[string]*DynamicTileConfig)
	}

	value[variant] = &DynamicTileConfig{
		Validator: validator,
		Builder:   builder,
	}
	cu.dynamicTileConfigs[tileType] = value
}
