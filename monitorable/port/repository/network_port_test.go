package repository

import (
	"errors"
	"testing"

	"github.com/jsdidierlaurent/monitoror/config"
	"github.com/jsdidierlaurent/monitoror/pkg/net/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRepository_CheckPort_Success(t *testing.T) {
	mockConn := new(mocks.Conn)
	mockConn.On("Close").Return(nil)
	mockDialer := new(mocks.Dialer)
	mockDialer.On("Dial", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockConn, nil)

	conf := &config.Config{
		PortConfig: config.PortConfig{
			Timeout: 1000,
		},
	}
	repository := NewNetworkPortRepository(conf)

	systemPortRepository, ok := repository.(*systemPortRepository)
	if assert.True(t, ok) {
		systemPortRepository.dialer = mockDialer

		assert.NoError(t, systemPortRepository.CheckPort("test", 1234))
		mockConn.AssertNumberOfCalls(t, "Close", 1)
		mockConn.AssertExpectations(t)
		mockDialer.AssertNumberOfCalls(t, "Dial", 1)
		mockDialer.AssertExpectations(t)
	}
}

func TestRepository_CheckPort_Failed(t *testing.T) {
	mockDialer := new(mocks.Dialer)
	mockDialer.On("Dial", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, errors.New("check port failed"))

	conf := &config.Config{
		PortConfig: config.PortConfig{
			Timeout: 1000,
		},
	}
	repository := NewNetworkPortRepository(conf)

	systemPortRepository, ok := repository.(*systemPortRepository)
	if assert.True(t, ok) {
		systemPortRepository.dialer = mockDialer

		assert.Error(t, systemPortRepository.CheckPort("test", 1234))
		mockDialer.AssertNumberOfCalls(t, "Dial", 1)
		mockDialer.AssertExpectations(t)
	}
}
