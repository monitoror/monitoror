package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/monitoror/monitoror/models"
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

func TestDelivery_ConfigHandler_Success(t *testing.T) {
	// Init
	ctx, res := initEcho()
	ctx.QueryParams().Set("url", "monitoror.example.com")

	config := &models.Config{
		Columns: 2,
	}

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("GetConfig", Anything).Return(config, nil)
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

func TestDelivery_ConfigHandler_QueryParamsError(t *testing.T) {
	// Init
	ctx, _ := initEcho()
	ctx.QueryParams().Del("hostname")

	mockUsecase := new(mocks.Usecase)
	handler := NewConfigDelivery(mockUsecase)

	// Test
	err := handler.GetConfig(ctx)
	assert.Error(t, err)
	assert.IsType(t, &MonitororError{}, err)
}

func TestDelivery_ConfigHandler_ErrorConfig(t *testing.T) {
	// Init
	ctx, _ := initEcho()
	ctx.QueryParams().Set("url", "monitoror.example.com")

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("GetConfig", Anything).Return(nil, errors.New("boom"))
	handler := NewConfigDelivery(mockUsecase)

	// Test
	if assert.Error(t, handler.GetConfig(ctx)) {
		mockUsecase.AssertNumberOfCalls(t, "GetConfig", 1)
		mockUsecase.AssertExpectations(t)
	}
}

func TestDelivery_ConfigHandler_ErrorVerify(t *testing.T) {
	// Init
	ctx, res := initEcho()
	ctx.QueryParams().Set("url", "monitoror.example.com")

	conf := &models.Config{
		Columns: 2,
		Errors:  []string{},
	}
	conf.AddErrors("boom")

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("GetConfig", Anything).Return(conf, nil)
	mockUsecase.On("Verify", Anything)
	handler := NewConfigDelivery(mockUsecase)

	// Expected
	json, err := json.Marshal(conf)
	assert.NoError(t, err, "unable to marshal config")

	// Test
	if assert.NoError(t, handler.GetConfig(ctx)) {
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(json), strings.TrimSpace(res.Body.String()))
		mockUsecase.AssertNumberOfCalls(t, "GetConfig", 1)
		mockUsecase.AssertNumberOfCalls(t, "Verify", 1)
		mockUsecase.AssertExpectations(t)
	}
}

func TestDelivery_ConfigHandler_ErrorHydrate(t *testing.T) {
	// Init
	ctx, res := initEcho()
	ctx.QueryParams().Set("url", "monitoror.example.com")

	conf := &models.Config{
		Columns: 2,
		Errors:  []string{},
	}
	conf.AddWarnings("boom")

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("GetConfig", Anything).Return(conf, nil)
	mockUsecase.On("Verify", Anything)
	mockUsecase.On("Hydrate", Anything, Anything)
	handler := NewConfigDelivery(mockUsecase)

	// Expected
	json, err := json.Marshal(conf)
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
