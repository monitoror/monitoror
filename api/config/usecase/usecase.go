package usecase

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/monitoror/monitoror/api/config"
	"github.com/monitoror/monitoror/models"

	"github.com/jsdidierlaurent/echo-middleware/cache"
)

// Versions
const (
	CurrentVersion = Version1000
	MinimalVersion = Version1000

	Version1000 = "1.0" // Initial version
)

const (
	EmptyTileType models.TileType = "EMPTY"
	GroupTileType models.TileType = "GROUP"

	DynamicTileStoreKeyPrefix = "monitoror.config.dynamicTile.key"
)

type (
	configUsecase struct {
		repository config.Repository

		configData *ConfigData

		// dynamic tile cache. used in case of timeout
		dynamicTileStore cache.Store
		cacheExpiration  time.Duration
	}
)

func NewConfigUsecase(repository config.Repository, store cache.Store, cacheExpiration int) config.Usecase {
	tileConfigs := make(map[models.TileType]map[string]*TileConfig)

	// Used for authorized type
	tileConfigs[EmptyTileType] = nil
	tileConfigs[GroupTileType] = nil

	return &configUsecase{
		repository:       repository,
		configData:       initConfigData(),
		dynamicTileStore: store,
		cacheExpiration:  time.Millisecond * time.Duration(cacheExpiration),
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
