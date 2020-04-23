package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"

	"github.com/monitoror/monitoror/api/config/mocks"
	"github.com/monitoror/monitoror/api/config/models"
)

func initEcho() (ctx echo.Context, res *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/api/v1/config", nil)
	res = httptest.NewRecorder()
	ctx = e.NewContext(req, res)

	return
}

func TestDelivery_ConfigHandler_Success(t *testing.T) {
	// Init
	ctx, res := initEcho()
	ctx.QueryParams().Set("url", "monitoror.example.com")

	config := &models.ConfigBag{}

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("GetConfig", Anything).Return(config)
	mockUsecase.On("Verify", Anything)
	mockUsecase.On("Hydrate", Anything, Anything)
	handler := NewConfigDelivery(mockUsecase)

	// Expected
	json, err := json.Marshal(config)
	assert.NoError(t, err, "unable to marshal config")

	// Test
	if assert.NoError(t, handler.GetConfig(ctx)) {
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(json), strings.TrimSpace(res.Body.String()))
		mockUsecase.AssertNumberOfCalls(t, "GetConfig", 1)
		mockUsecase.AssertNumberOfCalls(t, "Verify", 1)
		mockUsecase.AssertNumberOfCalls(t, "Hydrate", 1)
		mockUsecase.AssertExpectations(t)
	}
}

func TestDelivery_ConfigHandler_ErrorVerify(t *testing.T) {
	// Init
	ctx, res := initEcho()
	ctx.QueryParams().Set("url", "monitoror.example.com")

	conf := &models.ConfigBag{}

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("GetConfig", Anything).Return(conf)
	mockUsecase.On("Verify", Anything).Run(func(args Arguments) {
		conf.AddErrors(models.ConfigError{ID: "", Message: "boom", Data: models.ConfigErrorData{}})
	})
	handler := NewConfigDelivery(mockUsecase)

	// Test
	if assert.NoError(t, handler.GetConfig(ctx)) {
		strConf, _ := json.Marshal(conf)
		assert.Equal(t, string(strConf), strings.TrimSpace(res.Body.String()))
		assert.Equal(t, http.StatusOK, res.Code)

		mockUsecase.AssertNumberOfCalls(t, "GetConfig", 1)
		mockUsecase.AssertNumberOfCalls(t, "Verify", 1)
		mockUsecase.AssertExpectations(t)
	}
}

func TestDelivery_ConfigHandler_ErrorHydrate(t *testing.T) {
	// Init
	ctx, res := initEcho()
	ctx.QueryParams().Set("url", "monitoror.example.com")

	conf := &models.ConfigBag{}

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("GetConfig", Anything).Return(conf)
	mockUsecase.On("Verify", Anything)
	mockUsecase.On("Hydrate", Anything, Anything).Run(func(args Arguments) {
		conf.AddErrors(models.ConfigError{ID: "", Message: "boom", Data: models.ConfigErrorData{}})
	})
	handler := NewConfigDelivery(mockUsecase)

	// Test
	if assert.NoError(t, handler.GetConfig(ctx)) {
		strConf, _ := json.Marshal(conf)
		assert.Equal(t, string(strConf), strings.TrimSpace(res.Body.String()))
		assert.Equal(t, http.StatusOK, res.Code)

		mockUsecase.AssertNumberOfCalls(t, "GetConfig", 1)
		mockUsecase.AssertNumberOfCalls(t, "Verify", 1)
		mockUsecase.AssertNumberOfCalls(t, "Hydrate", 1)
		mockUsecase.AssertExpectations(t)
	}
}
