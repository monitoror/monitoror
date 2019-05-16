package errors

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/monitoror/monitoror/models/tiles"

	"github.com/stretchr/testify/assert"
)

func TestNoBuildError(t *testing.T) {
	// Init
	ctx, res := initErrorEcho()

	// Parameters
	tile := tiles.NewBuildTile("TEST")
	err := NewNoBuildError(tile)

	// Expected
	expectedTile := tile
	expectedTile.Status = tiles.WarningStatus
	expectedTile.Message = err.Error()
	j, e := json.Marshal(expectedTile)
	assert.NoError(t, e, "unable to marshal tile")

	// Test
	err.Send(ctx)
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, string(j), strings.TrimSpace(res.Body.String()))
}
