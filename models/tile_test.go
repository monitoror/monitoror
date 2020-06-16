package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTile_NewTile(t *testing.T) {
	tile := NewTile("TEST")
	assert.Equal(t, TileType("TEST"), tile.Type)
}

func TestTile_WithBuild(t *testing.T) {
	tile := NewTile("TEST").WithBuild()
	assert.NotNil(t, tile.Build)
}

func TestTile_WithValue(t *testing.T) {
	tile := NewTile("TEST").WithValue(MillisecondUnit)
	assert.NotNil(t, tile.Value)
	assert.Equal(t, MillisecondUnit, tile.Value.Unit)
}

func TestNewGeneratorTileType(t *testing.T) {
	generatorTest := NewGeneratorTileType("TEST")
	assert.Equal(t, "GENERATE:TEST", string(generatorTest))
}

func TestTileType_IsGenerator(t *testing.T) {
	generatorTest := NewGeneratorTileType("TEST")
	assert.True(t, generatorTest.IsGenerator())
}

func TestTileType_GetGeneratedTileType(t *testing.T) {
	generatorTest := NewGeneratorTileType("TEST")
	assert.Equal(t, "TEST", string(generatorTest.GetGeneratedTileType()))
}
