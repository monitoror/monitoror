package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTile(t *testing.T) {
	tile := NewTile("TEST")
	assert.Equal(t, TileType("TEST"), tile.Type)
}
