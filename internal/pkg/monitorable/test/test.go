package test

import (
	"testing"

	uiConfigMocks "github.com/monitoror/monitoror/api/config/mocks"
	coreConfig "github.com/monitoror/monitoror/config"
	serviceMocks "github.com/monitoror/monitoror/service/mocks"
	"github.com/monitoror/monitoror/service/store"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
)

type (
	MockMonitorableHelper interface {
		RouterAssertNumberOfCalls(t *testing.T, group int, get int)
		TileSettingsManagerAssertNumberOfCalls(t *testing.T, register int, registerGenerator int, enable int, enableGenerator int)
	}

	mockMonitorable struct {
		mockRouter      *serviceMocks.MonitorableRouter
		mockRouterGroup *serviceMocks.MonitorableRouterGroup

		mockTileSettingManager   *uiConfigMocks.TileSettingManager
		mockTileEnabler          *uiConfigMocks.TileEnabler
		mockTileGeneratorEnabler *uiConfigMocks.TileGeneratorEnabler
	}
)

func InitMockAndStore() (*store.Store, MockMonitorableHelper) {
	mockRouterGroup := new(serviceMocks.MonitorableRouterGroup)
	mockRouterGroup.On("GET", mock.AnythingOfType("string"), mock.AnythingOfType("echo.HandlerFunc"), mock.Anything).Return(&echo.Route{Path: "/path"})

	mockRouter := new(serviceMocks.MonitorableRouter)
	mockRouter.On("Group", mock.AnythingOfType("string"), mock.AnythingOfType("models.VariantName")).Return(mockRouterGroup)

	mockTileEnabler := new(uiConfigMocks.TileEnabler)
	mockTileEnabler.On("Enable",
		mock.AnythingOfType("models.VariantName"),
		mock.Anything, //	I didn't find a way to test that it's an validator.SimpleValidator interface	:(
		mock.AnythingOfType("string"),
		mock.AnythingOfType("int"),
	)
	mockTileGeneratorEnabler := new(uiConfigMocks.TileGeneratorEnabler)
	mockTileGeneratorEnabler.On("Enable",
		mock.AnythingOfType("models.VariantName"),
		mock.Anything, //	I didn't find a way to test that it's an validator.SimpleValidator interface :(
		mock.AnythingOfType("models.TileGeneratorFunction"),
	)

	mockTileSettingManager := new(uiConfigMocks.TileSettingManager)
	mockTileSettingManager.On("Register",
		mock.AnythingOfType("models.TileType"),
		mock.AnythingOfType("models.RawVersion"),
		mock.AnythingOfType("[]models.VariantName"),
	).Return(mockTileEnabler)
	mockTileSettingManager.On("RegisterGenerator",
		mock.AnythingOfType("models.TileType"),
		mock.AnythingOfType("models.RawVersion"),
		mock.AnythingOfType("[]models.VariantName"),
	).Return(mockTileGeneratorEnabler)

	return &store.Store{
			CoreConfig:         &coreConfig.Config{},
			MonitorableRouter:  mockRouter,
			TileSettingManager: mockTileSettingManager,
		},
		&mockMonitorable{
			mockRouter:               mockRouter,
			mockRouterGroup:          mockRouterGroup,
			mockTileSettingManager:   mockTileSettingManager,
			mockTileEnabler:          mockTileEnabler,
			mockTileGeneratorEnabler: mockTileGeneratorEnabler,
		}
}

func (m *mockMonitorable) RouterAssertNumberOfCalls(t *testing.T, group int, get int) {
	m.mockRouter.AssertNumberOfCalls(t, "Group", group)
	m.mockRouterGroup.AssertNumberOfCalls(t, "GET", get)
}

func (m *mockMonitorable) TileSettingsManagerAssertNumberOfCalls(t *testing.T, register int, registerGenerator int, enable int, enableGenerator int) {
	m.mockTileSettingManager.AssertNumberOfCalls(t, "Register", register)
	m.mockTileSettingManager.AssertNumberOfCalls(t, "RegisterGenerator", registerGenerator)
	m.mockTileEnabler.AssertNumberOfCalls(t, "Enable", enable)
	m.mockTileGeneratorEnabler.AssertNumberOfCalls(t, "Enable", enableGenerator)
}
