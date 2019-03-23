package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jsdidierlaurent/monitowall/config"
	"github.com/jsdidierlaurent/monitowall/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetPing_unit(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set(middlewares.BuildInfoContextKey, &config.BuildInfo{})
	c.Set(middlewares.ConfigContextKey, &config.Config{})

	var infoJSON = `{"build-info":{"git-commit":"","version":"","build-time":"","os":"","arch":""},"configuration":{"port":0,"cache":{"duration":0,"cleanup-Interval":0},"gitlab":{}}}`

	if assert.NoError(t, GetInfo(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, infoJSON, strings.TrimSpace(rec.Body.String()))
	}
}
