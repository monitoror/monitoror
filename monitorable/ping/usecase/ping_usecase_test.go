package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/jsdidierlaurent/monitoror/models/tiles"
	. "github.com/jsdidierlaurent/monitoror/monitorable/ping"
	"github.com/jsdidierlaurent/monitoror/monitorable/ping/mocks"
	"github.com/jsdidierlaurent/monitoror/monitorable/ping/model"

	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

var hostname = "test.com"

func TestUsecase_Ping_Success(t *testing.T) {
	// Init
	mockRepo := new(mocks.Repository)
	mockRepo.On("Ping", AnythingOfType("string")).Return(&model.Ping{
		Average: time.Second,
		Min:     time.Second,
		Max:     time.Second,
	}, nil)
	usecase := NewPingUsecase(mockRepo)

	// Params
	param := &model.PingParams{
		Hostname: hostname,
	}

	// Expected
	eTile := tiles.NewHealthTile(PingTileSubType)
	eTile.Label = hostname
	eTile.Status = tiles.SuccessStatus
	eTile.Message = "1s"

	// Test
	rTile, err := usecase.Ping(param)

	if assert.NoError(t, err) {
		assert.Equal(t, eTile, rTile)
		mockRepo.AssertNumberOfCalls(t, "Ping", 1)
		mockRepo.AssertExpectations(t)
	}
}

func TestPing_Fail(t *testing.T) {
	// Init
	mockRepo := new(mocks.Repository)
	mockRepo.On("Ping", AnythingOfType("string")).Return(nil, errors.New("ping error"))

	usecase := NewPingUsecase(mockRepo)

	// Params
	param := &model.PingParams{
		Hostname: hostname,
	}

	// Expected
	eTile := tiles.NewHealthTile(PingTileSubType)
	eTile.Label = hostname
	eTile.Status = tiles.FailStatus

	// Test
	rTile, err := usecase.Ping(param)

	if assert.NoError(t, err) {
		assert.Equal(t, eTile, rTile)
		mockRepo.AssertNumberOfCalls(t, "Ping", 1)
		mockRepo.AssertExpectations(t)
	}
}
