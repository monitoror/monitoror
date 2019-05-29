package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mErrors "github.com/monitoror/monitoror/models/errors"
	"github.com/monitoror/monitoror/monitorable/config/mocks"
	"github.com/monitoror/monitoror/monitorable/config/models"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func initEcho() (ctx echo.Context, res *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/api/v1/config", nil)
	res = httptest.NewRecorder()
	ctx = e.NewContext(req, res)

	return
}

func TestDelivery_GetConfig_Success(t *testing.T) {
	// Init
	ctx, res := initEcho()
	ctx.QueryParams().Set("url", "test.com")

	config := &models.Config{
		Columns: 2,
	}

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Config", Anything).Return(config, nil)
	mockUsecase.On("Verify", Anything).Return(nil)
	mockUsecase.On("Hydrate", Anything).Return(nil)
	handler := NewHttpConfigDelivery(mockUsecase)

	// Expected
	json, err := json.Marshal(config)
	assert.NoError(t, err, "unable to marshal config")

	// Test
	if assert.NoError(t, handler.GetConfig(ctx)) {
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(json), strings.TrimSpace(res.Body.String()))
		mockUsecase.AssertNumberOfCalls(t, "Config", 1)
		mockUsecase.AssertNumberOfCalls(t, "Verify", 1)
		mockUsecase.AssertNumberOfCalls(t, "Hydrate", 1)
		mockUsecase.AssertExpectations(t)
	}
}

func TestDelivery_GetConfig_QueryParamsError(t *testing.T) {
	// Init
	ctx, _ := initEcho()
	ctx.QueryParams().Del("hostname")

	mockUsecase := new(mocks.Usecase)
	handler := NewHttpConfigDelivery(mockUsecase)

	// Test
	err := handler.GetConfig(ctx)
	assert.Error(t, err)
	assert.IsType(t, &mErrors.QueryParamsError{}, err)
}

func TestDelivery_GetConfig_ErrorConfig(t *testing.T) {
	// Init
	ctx, _ := initEcho()
	ctx.QueryParams().Set("url", "test.com")

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Config", Anything).Return(nil, errors.New("boom"))
	handler := NewHttpConfigDelivery(mockUsecase)

	// Test
	if assert.Error(t, handler.GetConfig(ctx)) {
		mockUsecase.AssertNumberOfCalls(t, "Config", 1)
		mockUsecase.AssertExpectations(t)
	}
}

func TestDelivery_GetConfig_ErrorVerify(t *testing.T) {
	// Init
	ctx, _ := initEcho()
	ctx.QueryParams().Set("url", "test.com")

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Config", Anything).Return(nil, nil)
	mockUsecase.On("Verify", Anything).Return(errors.New("boom"))
	handler := NewHttpConfigDelivery(mockUsecase)

	// Test
	if assert.Error(t, handler.GetConfig(ctx)) {
		mockUsecase.AssertNumberOfCalls(t, "Config", 1)
		mockUsecase.AssertNumberOfCalls(t, "Verify", 1)
		mockUsecase.AssertExpectations(t)
	}
}

func TestDelivery_GetConfig_ErrorHydrate(t *testing.T) {
	// Init
	ctx, _ := initEcho()
	ctx.QueryParams().Set("url", "test.com")

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Config", Anything).Return(nil, nil)
	mockUsecase.On("Verify", Anything).Return(nil)
	mockUsecase.On("Hydrate", Anything).Return(errors.New("boom"))
	handler := NewHttpConfigDelivery(mockUsecase)

	// Test
	if assert.Error(t, handler.GetConfig(ctx)) {
		mockUsecase.AssertNumberOfCalls(t, "Config", 1)
		mockUsecase.AssertNumberOfCalls(t, "Verify", 1)
		mockUsecase.AssertNumberOfCalls(t, "Hydrate", 1)
		mockUsecase.AssertExpectations(t)
	}
}
