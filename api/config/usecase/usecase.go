package usecase

import (
	"time"

	"github.com/monitoror/monitoror/api/config"
	"github.com/monitoror/monitoror/api/config/models"
	coreConfig "github.com/monitoror/monitoror/config"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/service/registry"
	"github.com/monitoror/monitoror/service/store"

	"github.com/jsdidierlaurent/echo-middleware/cache"
)

const (
	EmptyTileType coreModels.TileType = "EMPTY"
	GroupTileType coreModels.TileType = "GROUP"

	TileGeneratorStoreKeyPrefix = "monitoror.config.tileGenerator.key"
)

type (
	configUsecase struct {
		repository config.Repository

		registry *registry.MetadataRegistry

		// namedConfigs used in GetConfig
		namedConfigs map[coreConfig.ConfigName]string

		// generator tile cache. used in case of timeout
		generatorTileStore cache.Store
		cacheExpiration    time.Duration

		initialMaxDelay int
	}
)

func NewConfigUsecase(repository config.Repository, store *store.Store) config.Usecase {
	tileConfigs := make(map[coreModels.TileType]map[string]*models.TileConfig)

	// Used for authorized type
	tileConfigs[EmptyTileType] = nil
	tileConfigs[GroupTileType] = nil

	return &configUsecase{
		repository:         repository,
		registry:           store.Registry.(*registry.MetadataRegistry),
		namedConfigs:       store.CoreConfig.NamedConfigs,
		generatorTileStore: store.CacheStore,
		cacheExpiration:    time.Millisecond * time.Duration(store.CoreConfig.DownstreamCacheExpiration),
		initialMaxDelay:    store.CoreConfig.InitialMaxDelay,
	}
}
