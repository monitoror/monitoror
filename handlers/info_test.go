package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jsdidierlaurent/monitowall/configs"
	"github.com/jsdidierlaurent/monitowall/models"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func initEcho() (ctx echo.Context, res *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/api/v1/info", nil)
	res = httptest.NewRecorder()
	ctx = e.NewContext(req, res)

	return
}

func TestGetInfo(t *testing.T) {
	// Init
	ctx, res := initEcho()
	emptyBuildInfo := configs.BuildInfo{}
	emptyConfig := configs.Config{}
	handler := HttpInfoHandler(&emptyBuildInfo, &emptyConfig)

	// Create expected value
	json, err := json.Marshal(&models.InfoResponse{
		BuildInfo: emptyBuildInfo,
		Config:    emptyConfig,
	})

	// Test
	assert.NoError(t, err)
	assert.NoError(t, handler.GetInfo(ctx))
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, string(json), strings.TrimSpace(res.Body.String()))
}
