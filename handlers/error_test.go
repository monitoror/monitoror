package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/monitoror/monitoror/models/tiles"

	mErrors "github.com/monitoror/monitoror/models/errors"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func initErrorEcho() (ctx echo.Context, res *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/error", nil)
	res = httptest.NewRecorder()
	ctx = e.NewContext(req, res)

	return
}

func TestHttpError_404(t *testing.T) {
	// Init
	ctx, res := initErrorEcho()

	// Parameters
	err := echo.NewHTTPError(http.StatusNotFound, "not found")

	// Expected
	apiError := ApiError{
		Code:    http.StatusNotFound,
		Message: "Not Found",
	}
	json, e := json.Marshal(apiError)
	assert.NoError(t, e, "unable to marshal tile")

	// Test
	HttpErrorHandler(err, ctx)

	assert.Equal(t, http.StatusNotFound, res.Code)
	assert.Equal(t, string(json), strings.TrimSpace(res.Body.String()))
}

func TestHttpError_500(t *testing.T) {
	// Init
	ctx, res := initErrorEcho()

	// Parameters
	err := errors.New("boom")

	// Expected
	apiError := ApiError{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
	}
	json, e := json.Marshal(apiError)
	assert.NoError(t, e, "unable to marshal tile")

	// Test
	HttpErrorHandler(err, ctx)

	assert.Equal(t, http.StatusInternalServerError, res.Code)
	assert.Equal(t, string(json), strings.TrimSpace(res.Body.String()))
}

func TestHttpError_TileError(t *testing.T) {
	// Init
	ctx, res := initErrorEcho()

	// Parameters
	err := mErrors.NewSystemError("BOOM", nil)

	// Expected
	tile := tiles.NewErrorTile("System Error", err.Error())
	j, e := json.Marshal(tile)
	assert.NoError(t, e, "unable to marshal tile")

	// Test
	HttpErrorHandler(err, ctx)

	assert.Equal(t, http.StatusInternalServerError, res.Code)
	assert.Equal(t, string(j), strings.TrimSpace(res.Body.String()))
}
