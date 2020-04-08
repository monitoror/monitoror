package usecase

import (
	"testing"

	"github.com/monitoror/monitoror/api/config/versions"
	coreModels "github.com/monitoror/monitoror/models"

	"github.com/stretchr/testify/assert"
)

func TestSettingManager(t *testing.T) {
	configUsecase := &configUsecase{configData: initConfigData()}
	configUsecase.Register("TEST", versions.CurrentVersion, []coreModels.VariantName{"test-variant"}).
		Enable("test-variant", nil, "test-route", 1000)
	configUsecase.RegisterGenerator("TEST", versions.CurrentVersion, []coreModels.VariantName{"test-variant"}).
		Enable("test-variant", nil, nil)

	data := configUsecase.configData
	assert.Len(t, data.tileSettings, 1)
	assert.Equal(t, versions.CurrentVersion, data.tileSettings["TEST"].MinimalVersion)
	assert.Equal(t, coreModels.TileType("TEST"), data.tileSettings["TEST"].TileType)
	assert.Equal(t, coreModels.VariantName("test-variant"), data.tileSettings["TEST"].Variants["test-variant"].VariantName)
	assert.True(t, data.tileSettings["TEST"].Variants["test-variant"].IsEnabled())
	assert.Nil(t, data.tileSettings["TEST"].Variants["test-variant"].GetValidator())

	assert.Len(t, data.tileGeneratorSettings, 1)
	generatorTileType := coreModels.NewGeneratorTileType("TEST")
	assert.Equal(t, versions.CurrentVersion, data.tileGeneratorSettings[generatorTileType].MinimalVersion)
	assert.Equal(t, generatorTileType, data.tileGeneratorSettings[generatorTileType].TileType)
	assert.Equal(t, coreModels.VariantName("test-variant"), data.tileGeneratorSettings[generatorTileType].Variants["test-variant"].VariantName)
	assert.True(t, data.tileGeneratorSettings[generatorTileType].Variants["test-variant"].IsEnabled())
	assert.Nil(t, data.tileGeneratorSettings[generatorTileType].Variants["test-variant"].GetValidator())
}

func TestTileSetting_Enable_Panic(t *testing.T) {
	configUsecase := &configUsecase{configData: initConfigData()}
	tileEnabler := configUsecase.Register("TEST", versions.CurrentVersion, []coreModels.VariantName{"test-variant"})
	assert.Panics(t, func() {
		tileEnabler.Enable("wrong-variant", nil, "", 100)
	})

	tileGeneratorEnabler := configUsecase.RegisterGenerator("TEST", versions.CurrentVersion, []coreModels.VariantName{"test-variant"})
	assert.Panics(t, func() {
		tileGeneratorEnabler.Enable("wrong-variant", nil, nil)
	})
}
