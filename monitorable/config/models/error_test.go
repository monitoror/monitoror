package models

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestConfigError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/api/v1/config", nil)
	res := httptest.NewRecorder()
	ctx := e.NewContext(req, res)

	err := NewConfigError()

	err.Add("Bug1")
	err.Add("Bug2")
	err.Send(ctx)

	assert.Equal(t, 2, err.Count())
	assert.Equal(t, "invalid configuration:\n - Bug1\n - Bug2", err.Error())

	assert.Equal(t, http.StatusBadRequest, res.Code)
	assert.Equal(t, `["Bug1","Bug2"]`, strings.TrimSpace(res.Body.String()))
}
