package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/monitoror/monitoror/models"
	. "github.com/monitoror/monitoror/monitorable/ping"
	"github.com/monitoror/monitoror/monitorable/ping/mocks"
	pingModels "github.com/monitoror/monitoror/monitorable/ping/models"

	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func TestUsecase_Ping_Success(t *testing.T) {
	// Init
	mockRepo := new(mocks.Repository)
	mockRepo.On("ExecutePing", AnythingOfType("string")).Return(&pingModels.Ping{
		Average: time.Second,
		Min:     time.Second,
		Max:     time.Second,
	}, nil)
	usecase := NewPingUsecase(mockRepo)

	// Params
	param := &pingModels.PingParams{
		Hostname: "monitoror.example.com",
	}

	// Expected
	eTile := models.NewTile(PingTileType).WithValue(models.MillisecondUnit)
	eTile.Label = param.Hostname
	eTile.Status = models.SuccessStatus
	eTile.Value.Values = append(eTile.Value.Values, "1000")

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
	param := &pingModels.PingParams{
		Hostname: "monitoror.example.com",
	}

	// Expected
	eTile := models.NewTile(PingTileType)
	eTile.Label = param.Hostname
	eTile.Status = models.FailedStatus

	// Test
	rTile, err := usecase.Ping(param)

	if assert.NoError(t, err) {
		assert.Equal(t, eTile, rTile)
		mockRepo.AssertNumberOfCalls(t, "ExecutePing", 1)
		mockRepo.AssertExpectations(t)
	}
}
