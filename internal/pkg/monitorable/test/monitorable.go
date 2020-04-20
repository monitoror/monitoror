package test

import (
	"testing"

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

		mockRegistry         *serviceMocks.Registry
		mockTileEnabler      *serviceMocks.TileEnabler
		mockGeneratorEnabler *serviceMocks.GeneratorEnabler
	}
)

func InitMockAndStore() (*store.Store, MockMonitorableHelper) {
	mockRouterGroup := new(serviceMocks.MonitorableRouterGroup)
	mockRouterGroup.On("GET", mock.AnythingOfType("string"), mock.AnythingOfType("echo.HandlerFunc"), mock.Anything).Return(&echo.Route{Path: "/path"})

	mockRouter := new(serviceMocks.MonitorableRouter)
	mockRouter.On("Group", mock.AnythingOfType("string"), mock.AnythingOfType("models.VariantName")).Return(mockRouterGroup)

	mockTileEnabler := new(serviceMocks.TileEnabler)
	mockTileEnabler.On("Enable",
		mock.AnythingOfType("models.VariantName"),
		mock.Anything, //	I didn't find a way to test that it's an params.Validator interface	:(
		mock.AnythingOfType("string"),
	)
	mockGeneratorEnabler := new(serviceMocks.GeneratorEnabler)
	mockGeneratorEnabler.On("Enable",
		mock.AnythingOfType("models.VariantName"),
		mock.Anything, //	I didn't find a way to test that it's an params.Validator interface :(
		mock.AnythingOfType("models.TileGeneratorFunction"),
	)

	mockRegistry := new(serviceMocks.Registry)
	mockRegistry.On("RegisterTile",
		mock.AnythingOfType("models.TileType"),
		mock.AnythingOfType("versions.RawVersion"),
		mock.AnythingOfType("[]models.VariantName"),
	).Return(mockTileEnabler)
	mockRegistry.On("RegisterGenerator",
		mock.AnythingOfType("models.TileType"),
		mock.AnythingOfType("versions.RawVersion"),
		mock.AnythingOfType("[]models.VariantName"),
	).Return(mockGeneratorEnabler)

	return &store.Store{
			CoreConfig:        &coreConfig.Config{},
			MonitorableRouter: mockRouter,
			Registry:          mockRegistry,
		},
		&mockMonitorable{
			mockRouter:           mockRouter,
			mockRouterGroup:      mockRouterGroup,
			mockRegistry:         mockRegistry,
			mockTileEnabler:      mockTileEnabler,
			mockGeneratorEnabler: mockGeneratorEnabler,
		}
}

func (m *mockMonitorable) RouterAssertNumberOfCalls(t *testing.T, group int, get int) {
	m.mockRouter.AssertNumberOfCalls(t, "Group", group)
	m.mockRouterGroup.AssertNumberOfCalls(t, "GET", get)
}

func (m *mockMonitorable) TileSettingsManagerAssertNumberOfCalls(t *testing.T, registerTile int, registerGenerator int, enableTile int, enableGenerator int) {
	m.mockRegistry.AssertNumberOfCalls(t, "RegisterTile", registerTile)
	m.mockRegistry.AssertNumberOfCalls(t, "RegisterGenerator", registerGenerator)
	m.mockTileEnabler.AssertNumberOfCalls(t, "Enable", enableTile)
	m.mockGeneratorEnabler.AssertNumberOfCalls(t, "Enable", enableGenerator)
}
