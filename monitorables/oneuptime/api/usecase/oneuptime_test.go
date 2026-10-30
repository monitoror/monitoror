package usecase

import (
	"errors"
	"testing"
	"time"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/oneuptime/api"
	"github.com/monitoror/monitoror/monitorables/oneuptime/api/mocks"
	"github.com/monitoror/monitoror/monitorables/oneuptime/api/models"

	"github.com/jsdidierlaurent/echo-middleware/cache"
	"github.com/stretchr/testify/assert"
)

func initUsecase(mockRepository api.Repository) api.Usecase {
	store := cache.NewGoCacheStore(time.Minute*5, time.Second)
	return NewOneuptimeUsecase(mockRepository, store, 1000)
}

func TestOneuptimeUsecase_Monitor_LoadMonitors_Error(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetMonitors").Return(nil, errors.New("boom"))

	tile, err := initUsecase(mockRepository).Monitor(&models.MonitorParams{ID: "1"})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &coreModels.MonitororError{}, err)
		assert.Equal(t, "unable to load monitors", err.Error())
		mockRepository.AssertNumberOfCalls(t, "GetMonitors", 1)
	}
}

func TestOneuptimeUsecase_Monitor_NotFound_Error(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetMonitors").Return([]models.Monitor{{ID: "1", Name: "Monitor 1"}}, nil)

	tile, err := initUsecase(mockRepository).Monitor(&models.MonitorParams{ID: "2"})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &coreModels.MonitororError{}, err)
		assert.Equal(t, "monitor 2 not found", err.Error())
	}
}

func TestOneuptimeUsecase_Monitor_NotFound_UsesCache(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetMonitors").Return([]models.Monitor{{ID: "1", Name: "Monitor 1"}}, nil)

	usecase := initUsecase(mockRepository)
	for i := 0; i < 3; i++ {
		_, err := usecase.Monitor(&models.MonitorParams{ID: "2"})
		assert.Error(t, err)
	}
	mockRepository.AssertNumberOfCalls(t, "GetMonitors", 1)
}

func TestOneuptimeUsecase_Monitor_WithCache(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetMonitors").Return([]models.Monitor{
		{ID: "1", Name: "Monitor 1", IsOperational: true},
		{ID: "2", Name: "Monitor 2", IsOffline: true},
	}, nil)

	usecase := initUsecase(mockRepository)

	tile, err := usecase.Monitor(&models.MonitorParams{ID: "1"})
	if assert.NoError(t, err) {
		assert.Equal(t, api.OneuptimeMonitorTileType, tile.Type)
		assert.Equal(t, "Monitor 1", tile.Label)
		assert.Equal(t, coreModels.SuccessStatus, tile.Status)
	}

	tile, err = usecase.Monitor(&models.MonitorParams{ID: "2"})
	if assert.NoError(t, err) {
		assert.Equal(t, "Monitor 2", tile.Label)
		assert.Equal(t, coreModels.FailedStatus, tile.Status)
	}

	// Second call is served by cache
	_, err = usecase.(*oneuptimeUsecase).loadMonitors()
	assert.NoError(t, err)
	mockRepository.AssertNumberOfCalls(t, "GetMonitors", 1)
}

func TestOneuptimeUsecase_MonitorGenerator_Error(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetMonitors").Return(nil, errors.New("boom"))

	results, err := initUsecase(mockRepository).MonitorGenerator(&models.MonitorGeneratorParams{})
	if assert.Error(t, err) {
		assert.Nil(t, results)
		assert.Equal(t, "unable to list monitors", err.Error())
	}
}

func TestOneuptimeUsecase_MonitorGenerator(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetMonitors").Return([]models.Monitor{
		{ID: "1", Name: "Monitor 2", Labels: []string{"label1"}},
		{ID: "2", Name: "Monitor 1", Labels: []string{"label1"}},
		{ID: "3", Name: "Monitor 3", Labels: []string{"label2"}},
		{ID: "4", Name: "Monitor 4", Labels: []string{"label1"}, IsDisabled: true},
	}, nil)

	results, err := initUsecase(mockRepository).MonitorGenerator(&models.MonitorGeneratorParams{SortBy: "name", Labels: " label1 , label3,"})
	if assert.NoError(t, err) {
		assert.Len(t, results, 2)
		assert.Equal(t, "Monitor 1", results[0].Label)
		assert.Equal(t, "2", results[0].Params.(models.MonitorParams).ID)
		assert.Equal(t, "Monitor 2", results[1].Label)
		assert.Equal(t, "1", results[1].Params.(models.MonitorParams).ID)
	}
}

func TestOneuptimeUsecase_MonitorGenerator_WithoutFilter(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetMonitors").Return([]models.Monitor{
		{ID: "1", Name: "Monitor 2"},
		{ID: "2", Name: "Monitor 1", Labels: []string{"label1"}},
	}, nil)

	results, err := initUsecase(mockRepository).MonitorGenerator(&models.MonitorGeneratorParams{})
	if assert.NoError(t, err) {
		assert.Len(t, results, 2)
		assert.Equal(t, "Monitor 2", results[0].Label)
		assert.Equal(t, "Monitor 1", results[1].Label)
	}
}

func TestOneuptimeUsecase_ParseLabels(t *testing.T) {
	assert.Nil(t, parseLabels(""))
	assert.Nil(t, parseLabels(" , "))
	assert.Equal(t, []string{"a", "b c"}, parseLabels(" a ,b c,"))
}

func TestOneuptimeUsecase_ParseMonitorStatus(t *testing.T) {
	assert.Equal(t, coreModels.SuccessStatus, parseMonitorStatus(&models.Monitor{IsOperational: true}))
	assert.Equal(t, coreModels.FailedStatus, parseMonitorStatus(&models.Monitor{IsOffline: true}))
	assert.Equal(t, coreModels.WarningStatus, parseMonitorStatus(&models.Monitor{}))
	assert.Equal(t, coreModels.DisabledStatus, parseMonitorStatus(&models.Monitor{IsDisabled: true, IsOperational: true}))
}
