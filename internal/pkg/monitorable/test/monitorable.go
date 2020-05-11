package test

import (
	"testing"

	coreConfig "github.com/monitoror/monitoror/config"
	registryMocks "github.com/monitoror/monitoror/registry/mocks"
	serviceMock "github.com/monitoror/monitoror/service/mocks"
	"github.com/monitoror/monitoror/store"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
)

type (
	MockMonitorableHelper interface {
		RouterAssertNumberOfCalls(t *testing.T, group int, get int)
		TileSettingsManagerAssertNumberOfCalls(t *testing.T, register int, registerGenerator int, enable int, enableGenerator int)
	}

	mockMonitorable struct {
		mockRouter      *serviceMock.MonitorableRouter
		mockRouterGroup *serviceMock.MonitorableRouterGroup

		mockRegistry         *registryMocks.Registry
		mockTileEnabler      *registryMocks.TileEnabler
		mockGeneratorEnabler *registryMocks.GeneratorEnabler
	}
)

func InitMockAndStore() (*store.Store, MockMonitorableHelper) {
	mockRouterGroup := new(serviceMock.MonitorableRouterGroup)
	mockRouterGroup.On("GET", mock.AnythingOfType("string"), mock.AnythingOfType("echo.HandlerFunc"), mock.Anything).Return(&echo.Route{Path: "/path"})

	mockRouter := new(serviceMock.MonitorableRouter)
	mockRouter.On("Group", mock.AnythingOfType("string"), mock.AnythingOfType("models.VariantName")).Return(mockRouterGroup)

	mockTileEnabler := new(registryMocks.TileEnabler)
	mockTileEnabler.On("Enable",
		mock.AnythingOfType("models.VariantName"),
		mock.Anything, //	I didn't find a way to test that it's an params.Validator interface	:(
		mock.AnythingOfType("string"),
	)
	mockGeneratorEnabler := new(registryMocks.GeneratorEnabler)
	mockGeneratorEnabler.On("Enable",
		mock.AnythingOfType("models.VariantName"),
		mock.Anything, //	I didn't find a way to test that it's an params.Validator interface :(
		mock.AnythingOfType("models.TileGeneratorFunction"),
	)

	mockRegistry := new(registryMocks.Registry)
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
			CoreConfig:        &coreConfig.CoreConfig{},
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
