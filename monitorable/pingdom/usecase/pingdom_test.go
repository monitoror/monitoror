package usecase

import (
	"errors"
	"testing"
	"time"

	. "github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/pingdom/mocks"
	pingdomModels "github.com/monitoror/monitoror/monitorable/pingdom/models"

	"github.com/AlekSi/pointer"
	"github.com/jsdidierlaurent/echo-middleware/cache"
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

	tile, err := pu.Check(&pingdomModels.CheckParams{ID: pointer.ToInt(1000)})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &models.MonitororError{}, err)
		assert.Equal(t, "unable to find check", err.Error())
		mockRepository.AssertNumberOfCalls(t, "GetCheck", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestPingdomUsecase_Check_NoBulk(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetCheck", AnythingOfType("int")).
		Return(&pingdomModels.Check{ID: 1000, Status: "up", Name: "Check 1"}, nil)

	store := cache.NewGoCacheStore(time.Minute*5, time.Second)
	config := &Pingdom{CacheExpiration: 1000}
	pu := NewPingdomUsecase(mockRepository, config, store)

	tile, err := pu.Check(&pingdomModels.CheckParams{ID: pointer.ToInt(1000)})
	if assert.NoError(t, err) {
		assert.NotNil(t, tile)
		assert.Equal(t, models.SuccessStatus, tile.Status)
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
		pingdomModels.Check{ID: 1000, Name: "Check 1", Status: "paused"},
		time.Second,
	)

	tile, err := pu.Check(&pingdomModels.CheckParams{ID: pointer.ToInt(1000)})
	if assert.NoError(t, err) {
		assert.NotNil(t, tile)
		assert.Equal(t, models.DisabledStatus, tile.Status)
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
	_ = castedTu.store.Set(castedTu.getTagsByIDStoreKey(1000), "", time.Minute)

	tile, err := pu.Check(&pingdomModels.CheckParams{ID: pointer.ToInt(1000)})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.IsType(t, &models.MonitororError{}, err)
		assert.Equal(t, "unable to find checks", err.Error())
		mockRepository.AssertNumberOfCalls(t, "GetChecks", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestPingdomUsecase_Check_Bulk(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetChecks", AnythingOfType("string")).
		Return([]pingdomModels.Check{
			{ID: 1000, Status: "up", Name: "Check 1"},
			{ID: 1100, Status: "down", Name: "Check 2"},
			{ID: 1200, Status: "paused", Name: "Check 3"},
		}, nil)

	store := cache.NewGoCacheStore(time.Minute*5, time.Second)
	config := &Pingdom{CacheExpiration: 1000}
	pu := NewPingdomUsecase(mockRepository, config, store)
	castedTu := pu.(*pingdomUsecase)

	// Force cache
	_ = castedTu.store.Set(castedTu.getTagsByIDStoreKey(1100), "", time.Minute)

	tile, err := pu.Check(&pingdomModels.CheckParams{ID: pointer.ToInt(1100)})
	if assert.NoError(t, err) {
		assert.NotNil(t, tile)
		assert.Equal(t, models.FailedStatus, tile.Status)
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

	_ = castedTu.store.Set(castedTu.getTagsByIDStoreKey(1000), "", time.Minute)
	_ = castedTu.store.Set(key,
		[]pingdomModels.Check{
			{ID: 1000, Status: "up", Name: "Check 1"},
			{ID: 1100, Status: "down", Name: "Check 2"},
			{ID: 1200, Status: "paused", Name: "Check 3"},
		},
		time.Second,
	)

	tile, err := pu.Check(&pingdomModels.CheckParams{ID: pointer.ToInt(1000)})
	if assert.NoError(t, err) {
		assert.NotNil(t, tile)
		assert.Equal(t, models.SuccessStatus, tile.Status)
		assert.Equal(t, "Check 1", tile.Label)
		mockRepository.AssertNumberOfCalls(t, "GetChecks", 0)
		mockRepository.AssertExpectations(t)
	}
}

func TestPingdomUsecase_Checks_Error(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetChecks", AnythingOfType("string")).Return(nil, errors.New("boom"))

	store := cache.NewGoCacheStore(time.Minute*5, time.Second)
	config := &Pingdom{CacheExpiration: 1000}
	pu := NewPingdomUsecase(mockRepository, config, store)

	results, err := pu.Checks(&pingdomModels.ChecksParams{SortBy: "name"})
	if assert.Error(t, err) {
		assert.Nil(t, results)
		mockRepository.AssertNumberOfCalls(t, "GetChecks", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestPingdomUsecase_Checks(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetChecks", AnythingOfType("string")).
		Return([]pingdomModels.Check{
			{ID: 1000, Status: "up", Name: "Check 2"},
			{ID: 1100, Status: "down", Name: "Check 1"},
			{ID: 1200, Status: "paused", Name: "Check 3"},
		}, nil)

	store := cache.NewGoCacheStore(time.Minute*5, time.Second)
	config := &Pingdom{CacheExpiration: 1000}
	pu := NewPingdomUsecase(mockRepository, config, store)

	results, err := pu.Checks(&pingdomModels.ChecksParams{SortBy: "name"})
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
	assert.Equal(t, models.SuccessStatus, parseStatus("up"))
	assert.Equal(t, models.FailedStatus, parseStatus("down"))
	assert.Equal(t, models.DisabledStatus, parseStatus("paused"))
	assert.Equal(t, models.UnknownStatus, parseStatus(""))
}
