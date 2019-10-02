package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/jsdidierlaurent/echo-middleware/cache"
	. "github.com/monitoror/monitoror/config"
	. "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/pingdom/mocks"
	"github.com/monitoror/monitoror/monitorable/pingdom/models"
	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func TestPingdomUsecase_Check_NoBulk_Error(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetCheck", AnythingOfType("int")).
		Return(nil, errors.New("boom"))

	store := cache.NewGoCacheStore(time.Minute*5, time.Second)
	config := &Pingdom{CacheExpiration: 1000}
	pu := NewPingdomUsecase(mockRepository, config, store)

	tile, err := pu.Check(&models.CheckParams{Id: pointer.ToInt(1000)})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &MonitororError{}, err)
		assert.Equal(t, "unable to found check", err.Error())
		mockRepository.AssertNumberOfCalls(t, "GetCheck", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestPingdomUsecase_Check_NoBulk(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetCheck", AnythingOfType("int")).
		Return(&models.Check{Id: 1000, Status: "up", Name: "Check 1"}, nil)

	store := cache.NewGoCacheStore(time.Minute*5, time.Second)
	config := &Pingdom{CacheExpiration: 1000}
	pu := NewPingdomUsecase(mockRepository, config, store)

	tile, err := pu.Check(&models.CheckParams{Id: pointer.ToInt(1000)})
	if assert.NoError(t, err) {
		assert.NotNil(t, tile)
		assert.Equal(t, SuccessStatus, tile.Status)
		assert.Equal(t, "Check 1", tile.Label)
		mockRepository.AssertNumberOfCalls(t, "GetCheck", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestPingdomUsecase_Check_NoBulk_WithCache(t *testing.T) {
	mockRepository := new(mocks.Repository)

	store := cache.NewGoCacheStore(time.Minute*5, time.Second)
	config := &Pingdom{CacheExpiration: 1000}
	pu := NewPingdomUsecase(mockRepository, config, store)
	castedTu := pu.(*pingdomUsecase)

	// Force cache
	key := castedTu.getCheckStoreKey(1000)
	_ = castedTu.store.Set(key,
		models.Check{Id: 1000, Name: "Check 1", Status: "paused"},
		time.Second,
	)

	tile, err := pu.Check(&models.CheckParams{Id: pointer.ToInt(1000)})
	if assert.NoError(t, err) {
		assert.NotNil(t, tile)
		assert.Equal(t, DisabledStatus, tile.Status)
		assert.Equal(t, "Check 1", tile.Label)
		mockRepository.AssertNumberOfCalls(t, "GetCheck", 0)
		mockRepository.AssertExpectations(t)
	}
}

func TestPingdomUsecase_Check_Bulk_Error(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetChecks", AnythingOfType("string")).
		Return(nil, errors.New("boom"))

	store := cache.NewGoCacheStore(time.Minute*5, time.Second)
	config := &Pingdom{CacheExpiration: 1000}
	pu := NewPingdomUsecase(mockRepository, config, store)
	castedTu := pu.(*pingdomUsecase)

	// Force cache
	_ = castedTu.store.Set(castedTu.getTagsByIdStoreKey(1000), "", time.Minute)

	tile, err := pu.Check(&models.CheckParams{Id: pointer.ToInt(1000)})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &MonitororError{}, err)
		assert.Equal(t, "unable to found checks", err.Error())
		mockRepository.AssertNumberOfCalls(t, "GetChecks", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestPingdomUsecase_Check_Bulk(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetChecks", AnythingOfType("string")).
		Return([]models.Check{
			{Id: 1000, Status: "up", Name: "Check 1"},
			{Id: 1100, Status: "down", Name: "Check 2"},
			{Id: 1200, Status: "paused", Name: "Check 3"},
		}, nil)

	store := cache.NewGoCacheStore(time.Minute*5, time.Second)
	config := &Pingdom{CacheExpiration: 1000}
	pu := NewPingdomUsecase(mockRepository, config, store)
	castedTu := pu.(*pingdomUsecase)

	// Force cache
	_ = castedTu.store.Set(castedTu.getTagsByIdStoreKey(1100), "", time.Minute)

	tile, err := pu.Check(&models.CheckParams{Id: pointer.ToInt(1100)})
	if assert.NoError(t, err) {
		assert.NotNil(t, tile)
		assert.Equal(t, FailedStatus, tile.Status)
		assert.Equal(t, "Check 2", tile.Label)
		mockRepository.AssertNumberOfCalls(t, "GetChecks", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestPingdomUsecase_Check_Bulk_WithCache(t *testing.T) {
	mockRepository := new(mocks.Repository)

	store := cache.NewGoCacheStore(time.Minute*5, time.Second)
	config := &Pingdom{CacheExpiration: 1000}
	pu := NewPingdomUsecase(mockRepository, config, store)
	castedTu := pu.(*pingdomUsecase)

	// Force cache
	key := castedTu.getChecksStoreKey("")

	_ = castedTu.store.Set(castedTu.getTagsByIdStoreKey(1000), "", time.Minute)
	_ = castedTu.store.Set(key,
		[]models.Check{
			{Id: 1000, Status: "up", Name: "Check 1"},
			{Id: 1100, Status: "down", Name: "Check 2"},
			{Id: 1200, Status: "paused", Name: "Check 3"},
		},
		time.Second,
	)

	tile, err := pu.Check(&models.CheckParams{Id: pointer.ToInt(1000)})
	if assert.NoError(t, err) {
		assert.NotNil(t, tile)
		assert.Equal(t, SuccessStatus, tile.Status)
		assert.Equal(t, "Check 1", tile.Label)
		mockRepository.AssertNumberOfCalls(t, "GetChecks", 0)
		mockRepository.AssertExpectations(t)
	}
}

func TestPingdomUsecase_ListDynamicTile_Error(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetChecks", AnythingOfType("string")).Return(nil, errors.New("boom"))

	store := cache.NewGoCacheStore(time.Minute*5, time.Second)
	config := &Pingdom{CacheExpiration: 1000}
	pu := NewPingdomUsecase(mockRepository, config, store)

	results, err := pu.ListDynamicTile(&models.ChecksParams{SortBy: "name"})
	if assert.Error(t, err) {
		assert.Nil(t, results)
		mockRepository.AssertNumberOfCalls(t, "GetChecks", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestPingdomUsecase_ListDynamicTile(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetChecks", AnythingOfType("string")).
		Return([]models.Check{
			{Id: 1000, Status: "up", Name: "Check 2"},
			{Id: 1100, Status: "down", Name: "Check 1"},
			{Id: 1200, Status: "paused", Name: "Check 3"},
		}, nil)

	store := cache.NewGoCacheStore(time.Minute*5, time.Second)
	config := &Pingdom{CacheExpiration: 1000}
	pu := NewPingdomUsecase(mockRepository, config, store)

	results, err := pu.ListDynamicTile(&models.ChecksParams{SortBy: "name"})
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
	assert.Equal(t, SuccessStatus, parseStatus("up"))
	assert.Equal(t, FailedStatus, parseStatus("down"))
	assert.Equal(t, DisabledStatus, parseStatus("paused"))
	assert.Equal(t, UnknownStatus, parseStatus(""))
}
