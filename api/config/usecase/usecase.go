package usecase

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/monitoror/monitoror/api/config"
	"github.com/monitoror/monitoror/api/config/models"
	coreModels "github.com/monitoror/monitoror/models"

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

		configData *ConfigData

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
		configData:         initConfigData(),
		generatorTileStore: store.CacheStore,
		cacheExpiration:    time.Millisecond * time.Duration(store.CoreConfig.DownstreamCacheExpiration),
		initialMaxDelay:    store.CoreConfig.InitialMaxDelay,
	}
}

// --- Utility functions ---
func keys(m interface{}) string {
	keys := reflect.ValueOf(m).MapKeys()
	strKeys := make([]string, len(keys))

	for i := 0; i < len(keys); i++ {
		strKeys[i] = fmt.Sprintf(`%v`, keys[i])
	}

	return strings.Join(strKeys, ", ")
}

func stringify(v interface{}) string {
	bytes, _ := json.Marshal(v)
	return string(bytes)
}
