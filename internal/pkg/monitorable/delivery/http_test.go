package delivery

import (
	"net/http/httptest"
	"testing"

	"github.com/monitoror/monitoror/api/config/mocks"
	"github.com/monitoror/monitoror/api/config/models"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type FakeValidator map[string]string

func (m *FakeValidator) Validate(_ *models.ConfigVersion) *models.ConfigError { return nil }

func TestBindAndValidateRequestParams(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/api/v1/xxx", nil)
	res := httptest.NewRecorder()
	ctx := e.NewContext(req, res)

	ctx.QueryParams().Add("test", "test")

	mockValidator := new(mocks.ParamsValidator)
	mockValidator.On("Validate", mock.Anything).Return(nil)
	assert.NoError(t, BindAndValidateRequestParams(ctx, mockValidator))

	fake := make(FakeValidator)
	assert.Error(t, BindAndValidateRequestParams(ctx, &fake))

	mockValidator2 := new(mocks.ParamsValidator)
	mockValidator2.On("Validate", mock.Anything).Return(&models.ConfigError{Message: "boom"})
	err := BindAndValidateRequestParams(ctx, mockValidator2)
	assert.Error(t, err)
	assert.Equal(t, "boom", err.Error())
}
