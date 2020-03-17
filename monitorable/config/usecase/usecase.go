package usecase

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/monitoror/monitoror/models"
	monitorableConfig "github.com/monitoror/monitoror/monitorable/config"
	"github.com/monitoror/monitoror/pkg/monitoror/builder"
	"github.com/monitoror/monitoror/pkg/monitoror/utils"

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
		repository monitorableConfig.Repository

		tileConfigs        map[models.TileType]map[string]*TileConfig
		dynamicTileConfigs map[models.TileType]map[string]*DynamicTileConfig

		// dynamic tile cache. used in case of timeout
		dynamicTileStore cache.Store
		cacheExpiration  time.Duration
	}

	// TileConfig struct is used by GetConfig endpoint to check / hydrate config
	TileConfig struct {
		Validator       utils.Validator
		Path            string
		InitialMaxDelay int
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
		repository:         repository,
		tileConfigs:        tileConfigs,
		dynamicTileConfigs: dynamicTileConfigs,
		dynamicTileStore:   store,
		cacheExpiration:    time.Millisecond * time.Duration(downstreamStoreExpiration),
	}
}

func (cu *configUsecase) RegisterTile(
	tileType models.TileType, variant string, clientConfigValidator utils.Validator, path string, initialMaxDelay int,
) {
	value, exists := cu.tileConfigs[tileType]
	if !exists {
		value = make(map[string]*TileConfig)
		cu.tileConfigs[tileType] = value
	}

	value[variant] = &TileConfig{
		Path:            path,
		Validator:       clientConfigValidator,
		InitialMaxDelay: initialMaxDelay,
	}
}

func (cu *configUsecase) RegisterDynamicTile(
	tileType models.TileType, variant string, clientConfigValidator utils.Validator, builder builder.DynamicTileBuilder,
) {
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

func (cu *configUsecase) DisableTile(tileType models.TileType, variant string) {
	// TODO
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
