package usecase

import (
	"errors"
	"testing"
	"time"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/pingdom/api"
	"github.com/monitoror/monitoror/monitorables/pingdom/api/mocks"
	"github.com/monitoror/monitoror/monitorables/pingdom/api/models"

	"github.com/AlekSi/pointer"
	"github.com/jsdidierlaurent/echo-middleware/cache"
	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func initUsecase(mockRepository *mocks.Repository) api.Usecase {
	store := cache.NewGoCacheStore(time.Minute*5, time.Second)
	pu := NewPingdomUsecase(mockRepository, store, 1000)
	return pu
}

func TestPingdomUsecase_Check_NoBulk_Error(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetCheck", AnythingOfType("int")).
		Return(nil, errors.New("boom"))

	pu := initUsecase(mockRepository)

	tile, err := pu.Check(&models.CheckParams{ID: pointer.ToInt(1000)})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &coreModels.MonitororError{}, err)
		assert.Equal(t, "unable to find check", err.Error())
		mockRepository.AssertNumberOfCalls(t, "GetCheck", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestPingdomUsecase_Check_NoBulk(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetCheck", AnythingOfType("int")).
		Return(&models.Check{ID: 1000, Status: "up", Name: "Check 1"}, nil)

	pu := initUsecase(mockRepository)

	tile, err := pu.Check(&models.CheckParams{ID: pointer.ToInt(1000)})
	if assert.NoError(t, err) {
		assert.NotNil(t, tile)
		assert.Equal(t, coreModels.SuccessStatus, tile.Status)
		assert.Equal(t, "Check 1", tile.Label)
		mockRepository.AssertNumberOfCalls(t, "GetCheck", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestPingdomUsecase_Check_NoBulk_WithCache(t *testing.T) {
	mockRepository := new(mocks.Repository)

	pu := initUsecase(mockRepository)
	castedTu := pu.(*pingdomUsecase)

	// Force cache
	key := castedTu.getCheckStoreKey(1000)
	_ = castedTu.store.Set(key,
		models.Check{ID: 1000, Name: "Check 1", Status: "paused"},
		time.Second,
	)

	tile, err := pu.Check(&models.CheckParams{ID: pointer.ToInt(1000)})
	if assert.NoError(t, err) {
		assert.NotNil(t, tile)
		assert.Equal(t, coreModels.DisabledStatus, tile.Status)
		assert.Equal(t, "Check 1", tile.Label)
		mockRepository.AssertNumberOfCalls(t, "GetCheck", 0)
		mockRepository.AssertExpectations(t)
	}
}

func TestPingdomUsecase_Check_Bulk_Error(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetChecks", AnythingOfType("string")).
		Return(nil, errors.New("boom"))

	pu := initUsecase(mockRepository)
	castedTu := pu.(*pingdomUsecase)

	// Force cache
	_ = castedTu.store.Set(castedTu.getTagsByIDStoreKey(1000), "", time.Minute)

	tile, err := pu.Check(&models.CheckParams{ID: pointer.ToInt(1000)})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &coreModels.MonitororError{}, err)
		assert.Equal(t, "unable to find checks", err.Error())
		mockRepository.AssertNumberOfCalls(t, "GetChecks", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestPingdomUsecase_Check_Bulk(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetChecks", AnythingOfType("string")).
		Return([]models.Check{
			{ID: 1000, Status: "up", Name: "Check 1"},
			{ID: 1100, Status: "down", Name: "Check 2"},
			{ID: 1200, Status: "paused", Name: "Check 3"},
		}, nil)

	pu := initUsecase(mockRepository)
	castedTu := pu.(*pingdomUsecase)

	// Force cache
	_ = castedTu.store.Set(castedTu.getTagsByIDStoreKey(1100), "", time.Minute)

	tile, err := pu.Check(&models.CheckParams{ID: pointer.ToInt(1100)})
	if assert.NoError(t, err) {
		assert.NotNil(t, tile)
		assert.Equal(t, coreModels.FailedStatus, tile.Status)
		assert.Equal(t, "Check 2", tile.Label)
		mockRepository.AssertNumberOfCalls(t, "GetChecks", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestPingdomUsecase_Check_Bulk_WithCache(t *testing.T) {
	mockRepository := new(mocks.Repository)

	pu := initUsecase(mockRepository)
	castedTu := pu.(*pingdomUsecase)

	// Force cache
	key := castedTu.getChecksStoreKey("")

	_ = castedTu.store.Set(castedTu.getTagsByIDStoreKey(1000), "", time.Minute)
	_ = castedTu.store.Set(key,
		[]models.Check{
			{ID: 1000, Status: "up", Name: "Check 1"},
			{ID: 1100, Status: "down", Name: "Check 2"},
			{ID: 1200, Status: "paused", Name: "Check 3"},
		},
		time.Second,
	)

	tile, err := pu.Check(&models.CheckParams{ID: pointer.ToInt(1000)})
	if assert.NoError(t, err) {
		assert.NotNil(t, tile)
		assert.Equal(t, coreModels.SuccessStatus, tile.Status)
		assert.Equal(t, "Check 1", tile.Label)
		mockRepository.AssertNumberOfCalls(t, "GetChecks", 0)
		mockRepository.AssertExpectations(t)
	}
}

func TestPingdomUsecase_Checks_Error(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetChecks", AnythingOfType("string")).Return(nil, errors.New("boom"))

	pu := initUsecase(mockRepository)

	results, err := pu.Checks(&models.ChecksParams{SortBy: "name"})
	if assert.Error(t, err) {
		assert.Nil(t, results)
		mockRepository.AssertNumberOfCalls(t, "GetChecks", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestPingdomUsecase_Checks(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetChecks", AnythingOfType("string")).
		Return([]models.Check{
			{ID: 1000, Status: "up", Name: "Check 2"},
			{ID: 1100, Status: "down", Name: "Check 1"},
			{ID: 1200, Status: "paused", Name: "Check 3"},
		}, nil)

	pu := initUsecase(mockRepository)

	results, err := pu.Checks(&models.ChecksParams{SortBy: "name"})
	if assert.NoError(t, err) {
		assert.NotNil(t, results)
		assert.Len(t, results, 2)
		assert.Equal(t, "Check 1", results[0].Label)
		assert.Equal(t, "Check 2", results[1].Label)
		mockRepository.AssertNumberOfCalls(t, "GetChecks", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestPingdomUsecase_ParseStatus(t *testing.T) {
	assert.Equal(t, coreModels.SuccessStatus, parseStatus("up"))
	assert.Equal(t, coreModels.FailedStatus, parseStatus("down"))
	assert.Equal(t, coreModels.DisabledStatus, parseStatus("paused"))
	assert.Equal(t, coreModels.UnknownStatus, parseStatus(""))
}
