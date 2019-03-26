package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/jsdidierlaurent/monitowall/models/tiles"
	"github.com/jsdidierlaurent/monitowall/monitorable/ping/mocks"
	"github.com/jsdidierlaurent/monitowall/monitorable/ping/model"
	pkgMock "github.com/jsdidierlaurent/monitowall/pkg/bind/mocks"

	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

var hostname = "test.com"

func TestUsecase_Ping_Success(t *testing.T) {
	// Init
	mockParamBinder := new(pkgMock.Binder)
	mockParamBinder.On("Bind", Anything).Return(func(p interface{}) error {
		param, ok := p.(*model.PingParams)
		assert.True(t, ok)
		param.Hostname = hostname
		return nil
	})

	mockRepo := new(mocks.Repository)
	mockRepo.On("Ping", AnythingOfType("string")).Return(&model.Ping{
		Average: time.Second,
		Min:     time.Second,
		Max:     time.Second,
	}, nil)
	usecase := NewPingUsecase(mockRepo)

	// Expected
	eTile := tiles.NewHealthTile(PingTileSubType)
	eTile.Label = hostname
	eTile.Status = tiles.SuccessStatus
	eTile.Message = "1s"

	// Test
	rTile, err := usecase.Ping(mockParamBinder)

	assert.NoError(t, err)
	assert.Equal(t, eTile, rTile)
}

func TestPing_Fail(t *testing.T) {
	// Init
	mockParamBinder := new(pkgMock.Binder)
	mockParamBinder.On("Bind", Anything).Return(func(p interface{}) error {
		param, ok := p.(*model.PingParams)
		assert.True(t, ok)
		param.Hostname = hostname
		return nil
	})

	mockRepo := new(mocks.Repository)
	mockRepo.On("Ping", AnythingOfType("string")).Return(nil, errors.New("ping error"))

	usecase := NewPingUsecase(mockRepo)

	// Expected
	eTile := tiles.NewHealthTile(PingTileSubType)
	eTile.Label = hostname
	eTile.Status = tiles.FailStatus

	// Test
	rTile, err := usecase.Ping(mockParamBinder)

	assert.NoError(t, err)
	assert.Equal(t, eTile, rTile)
}

func TestPing_Error(t *testing.T) {
	// Init
	mockParamBinder := new(pkgMock.Binder)
	mockParamBinder.On("Bind", Anything).Return(errors.New("bind error"))

	mockRepo := new(mocks.Repository)
	usecase := NewPingUsecase(mockRepo)

	// Expected
	eTile := tiles.NewHealthTile(PingTileSubType)
	eTile.Label = hostname
	eTile.Status = tiles.FailStatus

	// Test
	_, err := usecase.Ping(mockParamBinder)
	assert.Error(t, err)
}
