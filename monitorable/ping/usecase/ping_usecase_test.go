package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/monitoror/monitoror/models/tiles"
	. "github.com/monitoror/monitoror/monitorable/ping"
	"github.com/monitoror/monitoror/monitorable/ping/mocks"
	"github.com/monitoror/monitoror/monitorable/ping/model"

	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func TestUsecase_Ping_Success(t *testing.T) {
	// Init
	mockRepo := new(mocks.Repository)
	mockRepo.On("Ping", Anything, AnythingOfType("string")).Return(&model.Ping{
		Average: time.Second,
		Min:     time.Second,
		Max:     time.Second,
	}, nil)
	usecase := NewPingUsecase(mockRepo)

	// Params
	param := &model.PingParams{
		Hostname: "test.com",
	}

	// Expected
	eTile := tiles.NewHealthTile(PingTileSubType)
	eTile.Label = param.Hostname
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

func TestUsecase_Ping_Fail(t *testing.T) {
	// Init
	mockRepo := new(mocks.Repository)
	mockRepo.On("Ping", Anything, AnythingOfType("string")).Return(nil, errors.New("ping error"))

	usecase := NewPingUsecase(mockRepo)

	// Params
	param := &model.PingParams{
		Hostname: "test.com",
	}

	// Expected
	eTile := tiles.NewHealthTile(PingTileSubType)
	eTile.Label = param.Hostname
	eTile.Status = tiles.FailedStatus

	// Test
	rTile, err := usecase.Ping(param)

	if assert.NoError(t, err) {
		assert.Equal(t, eTile, rTile)
		mockRepo.AssertNumberOfCalls(t, "Ping", 1)
		mockRepo.AssertExpectations(t)
	}
}
