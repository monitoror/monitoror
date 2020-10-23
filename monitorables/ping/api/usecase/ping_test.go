package usecase

import (
	"errors"
	"testing"
	"time"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/ping/api"
	"github.com/monitoror/monitoror/monitorables/ping/api/mocks"
	"github.com/monitoror/monitoror/monitorables/ping/api/models"

	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func TestUsecase_Ping_Success(t *testing.T) {
	// Init
	mockRepo := new(mocks.Repository)
	mockRepo.On("ExecutePing", AnythingOfType("string")).Return(&models.Ping{
		Average: time.Second,
		Min:     time.Second,
		Max:     time.Second,
	}, nil)
	usecase := NewPingUsecase(mockRepo)

	// Params
	param := &models.PingParams{
		Hostname: "monitoror.example.com",
	}

	// Expected
	eTile := coreModels.NewTile(api.PingTileType).WithMetrics(coreModels.MillisecondUnit)
	eTile.Label = param.Hostname
	eTile.Status = coreModels.SuccessStatus
	eTile.Metrics.Values = append(eTile.Metrics.Values, "1000")

	// Test
	rTile, err := usecase.Ping(param)

	if assert.NoError(t, err) {
		assert.Equal(t, eTile, rTile)
		mockRepo.AssertNumberOfCalls(t, "ExecutePing", 1)
		mockRepo.AssertExpectations(t)
	}
}

func TestUsecase_Ping_Fail(t *testing.T) {
	// Init
	mockRepo := new(mocks.Repository)
	mockRepo.On("ExecutePing", AnythingOfType("string")).Return(nil, errors.New("ping error"))

	usecase := NewPingUsecase(mockRepo)

	// Params
	param := &models.PingParams{
		Hostname: "monitoror.example.com",
	}

	// Expected
	eTile := coreModels.NewTile(api.PingTileType)
	eTile.Label = param.Hostname
	eTile.Status = coreModels.FailedStatus

	// Test
	rTile, err := usecase.Ping(param)

	if assert.NoError(t, err) {
		assert.Equal(t, eTile, rTile)
		mockRepo.AssertNumberOfCalls(t, "ExecutePing", 1)
		mockRepo.AssertExpectations(t)
	}
}
