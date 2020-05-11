package registry

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/monitoror/monitoror/api/config/versions"
	"github.com/monitoror/monitoror/models"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/models/mocks"
)

func TestRegistry_Tiles(t *testing.T) {
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

func TestEnabler_Enable_Panic(t *testing.T) {
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

func TestRegistry_Monitorable(t *testing.T) {
	mockMonitorable := new(mocks.Monitorable)
	mockMonitorable.On("GetDisplayName").Return("Monitorable Mock")
	mockMonitorable.On("GetVariantsNames").Return([]models.VariantName{models.DefaultVariantName, "variant1", "variant2"})
	mockMonitorable.On("Validate", mock.AnythingOfType("models.VariantName")).Return(true, nil).Once()
	mockMonitorable.On("Validate", mock.AnythingOfType("models.VariantName")).Return(false, nil).Once()
	mockMonitorable.On("Validate", mock.AnythingOfType("models.VariantName")).Return(false, []error{errors.New("boom")}).Once()
	mockMonitorable.On("Enable", mock.AnythingOfType("models.VariantName"))

	registry := NewRegistry()
	registry.RegisterMonitorable(mockMonitorable)

	if assert.Len(t, registry.GetMonitorables(), 1) {
		m := registry.GetMonitorables()[0]
		assert.Equal(t, "Monitorable Mock", m.Monitorable.GetDisplayName())

		assert.Equal(t, models.DefaultVariantName, m.VariantsMetadata[0].VariantName)
		assert.Equal(t, true, m.VariantsMetadata[0].Enabled)
		assert.Equal(t, models.VariantName("variant1"), m.VariantsMetadata[1].VariantName)
		assert.Equal(t, false, m.VariantsMetadata[1].Enabled)
		assert.Equal(t, models.VariantName("variant2"), m.VariantsMetadata[2].VariantName)
		assert.Equal(t, false, m.VariantsMetadata[2].Enabled)
	}
}
