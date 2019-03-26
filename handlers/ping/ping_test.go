package ping

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/jsdidierlaurent/monitowall/renderings"

	"github.com/jsdidierlaurent/monitowall/models/mocks"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

var hostname, message = "test", "1ms"

func initEcho() (rec *httptest.ResponseRecorder, context echo.Context) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec = httptest.NewRecorder()
	context = e.NewContext(req, rec)
	context.QueryParams().Set("hostname", "test")
	return
}

func TestGetPing_Success(t *testing.T) {
	rec, context := initEcho()

	// Mock Configuration
	mockPing := new(mocks.PingModel)
	mockPing.On("Ping", AnythingOfType("string")).Return(message, nil)
	handler := NewHandler(mockPing)

	// Expected Result
	response := newResponse()
	response.Label = hostname
	response.Message = message
	response.Status = SuccessStatus
	json, _ := json.Marshal(response)

	assert.NoError(t, handler.GetPing(context))
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, string(json), strings.TrimSpace(rec.Body.String()))
}

func TestGetPing_Error(t *testing.T) {
	rec, context := initEcho()

	// Mock Configuration
	mockPing := new(mocks.PingModel)
	mockPing.On("Ping", AnythingOfType("string")).Return("", fmt.Errorf("Test"))
	handler := NewHandler(mockPing)

	// Expected Result
	response := newResponse()
	response.Label = hostname
	response.Status = FailStatus
	json, _ := json.Marshal(response)

	assert.NoError(t, handler.GetPing(context))
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, string(json), strings.TrimSpace(rec.Body.String()))
}
