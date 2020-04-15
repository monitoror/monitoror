package registry

import (
	"testing"

	"github.com/monitoror/monitoror/api/config/versions"
	"github.com/monitoror/monitoror/models"
	coreModels "github.com/monitoror/monitoror/models"

	"github.com/stretchr/testify/assert"
)

func TestSettingManager(t *testing.T) {
	registry := NewRegistry()
	registry.RegisterTile("TEST", versions.CurrentVersion, []coreModels.VariantName{"test-variant"}).
		Enable("test-variant", nil, "test-route")
	registry.RegisterGenerator("TEST", versions.CurrentVersion, []coreModels.VariantName{"test-variant"}).
		Enable("test-variant", nil, nil)

	assert.Len(t, registry.TileMetadata, 1)
	assert.Equal(t, versions.CurrentVersion, registry.TileMetadata["TEST"].GetMinimalVersion())
	assert.Equal(t, []models.VariantName{"test-variant"}, registry.TileMetadata["TEST"].GetVariantsNames())
	assert.Equal(t, coreModels.TileType("TEST"), registry.TileMetadata["TEST"].TileType)
	variant, exists := registry.TileMetadata["TEST"].GetVariant("test-variant")
	if assert.True(t, exists) {
		assert.True(t, variant.IsEnabled())
		assert.Nil(t, variant.GetValidator())
	}

	assert.Len(t, registry.GeneratorMetadata, 1)
	generatorTileType := coreModels.NewGeneratorTileType("TEST")
	assert.Equal(t, versions.CurrentVersion, registry.GeneratorMetadata[generatorTileType].GetMinimalVersion())
	assert.Equal(t, []models.VariantName{"test-variant"}, registry.GeneratorMetadata[generatorTileType].GetVariantsNames())
	assert.Equal(t, generatorTileType, registry.GeneratorMetadata[generatorTileType].TileType)
	variant, exists = registry.GeneratorMetadata[generatorTileType].GetVariant("test-variant")
	if assert.True(t, exists) {
		assert.True(t, variant.IsEnabled())
		assert.Nil(t, variant.GetValidator())
	}
}

func TestTileSetting_Enable_Panic(t *testing.T) {
	registry := NewRegistry()
	tileEnabler := registry.RegisterTile("TEST", versions.CurrentVersion, []coreModels.VariantName{"test-variant"})
	assert.Panics(t, func() {
		tileEnabler.Enable("wrong-variant", nil, "")
	})

	tileGeneratorEnabler := registry.RegisterGenerator("TEST", versions.CurrentVersion, []coreModels.VariantName{"test-variant"})
	assert.Panics(t, func() {
		tileGeneratorEnabler.Enable("wrong-variant", nil, nil)
	})
}
