package test

import (
	coreConfigMocks "github.com/monitoror/monitoror/api/config/mocks"
	coreConfig "github.com/monitoror/monitoror/config"
	serviceMocks "github.com/monitoror/monitoror/service/mocks"
	"github.com/monitoror/monitoror/service/store"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
)

func InitMockAndStore() (*serviceMocks.MonitorableRouter, *serviceMocks.MonitorableRouterGroup, *coreConfigMocks.Manager, *store.Store) {
	mockRouterGroup := new(serviceMocks.MonitorableRouterGroup)
	mockRouterGroup.On("GET",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("echo.HandlerFunc"),
		mock.Anything,
	).Return(&echo.Route{Path: "/path"})

	mockRouter := new(serviceMocks.MonitorableRouter)
	mockRouter.On("Group",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("models.VariantName"),
	).Return(mockRouterGroup)

	mockConfigManager := new(coreConfigMocks.Manager)
	mockConfigManager.On("RegisterTile",
		mock.AnythingOfType("models.TileType"),
		mock.AnythingOfType("[]models.VariantName"),
		mock.AnythingOfType("models.RawVersion"),
	)
	mockConfigManager.On("EnableTile",
		mock.AnythingOfType("models.TileType"),
		mock.AnythingOfType("models.VariantName"),
		mock.Anything, //	I didn't find a way to test that it's an validator.SimpleValidator interface	:(
		mock.AnythingOfType("string"),
		mock.AnythingOfType("int"),
	)
	mockConfigManager.On("EnableDynamicTile",
		mock.AnythingOfType("models.TileType"),
		mock.AnythingOfType("models.VariantName"),
		mock.Anything, //	I didn't find a way to test that it's an validator.SimpleValidator interface :(
		mock.AnythingOfType("config.DynamicTileBuilder"),
	)

	return mockRouter, mockRouterGroup, mockConfigManager, &store.Store{
		CoreConfig:        &coreConfig.Config{},
		MonitorableRouter: mockRouter,
		UIConfigManager:   mockConfigManager,
	}
}
