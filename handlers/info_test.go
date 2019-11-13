package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/monitoror/monitoror/cli/version"
	"github.com/monitoror/monitoror/models"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func initInfoEcho() (ctx echo.Context, res *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/api/v1/info", nil)
	res = httptest.NewRecorder()
	ctx = e.NewContext(req, res)

	return
}

func TestGetInfo(t *testing.T) {
	// Init
	ctx, res := initInfoEcho()
	handler := NewHTTPInfoDelivery()

	// Create expected value
	json, err := json.Marshal(models.NewInfoResponse(version.Version, version.GitCommit, version.BuildTime))
	assert.NoError(t, err, "unable to marshal InfoResponse")

	// Test
	if assert.NoError(t, handler.GetInfo(ctx)) {
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(json), strings.TrimSpace(res.Body.String()))
	}
}
