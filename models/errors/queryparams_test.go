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

func TestQueryParamsError(t *testing.T) {
	// Init
	ctx, res := initErrorEcho()

	// Parameters
	err := NewQueryParamsError(errors.New("BOOM"))

	// Expected
	tile := tiles.NewErrorTile("Wrong Configuration", err.Error())
	j, e := json.Marshal(tile)
	assert.NoError(t, e, "unable to marshal tile")

	// Test
	err.Send(ctx)

	assert.Equal(t, http.StatusBadRequest, res.Code)
	assert.Equal(t, string(j), strings.TrimSpace(res.Body.String()))
}
