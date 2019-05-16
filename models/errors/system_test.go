package errors

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/monitoror/monitoror/models/tiles"
	"github.com/stretchr/testify/assert"
)

func TestSystemError(t *testing.T) {
	// Init
	ctx, res := initErrorEcho()

	// Parameters
	message := "system error"
	err := NewSystemError(message, errors.New("BOOM"))

	// Expected
	tile := tiles.NewErrorTile("System Error", err.Error())
	j, e := json.Marshal(tile)
	assert.NoError(t, e, "unable to marshal tile")

	// Test
	err.Send(ctx)
	assert.Equal(t, http.StatusInternalServerError, res.Code)
	assert.Equal(t, string(j), strings.TrimSpace(res.Body.String()))
}
