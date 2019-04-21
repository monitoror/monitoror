package repository

import (
	"context"
	"errors"
	"testing"

	pkgNet "github.com/monitoror/monitoror/pkg/net"

	. "github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/pkg/net/mocks"

	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func initRepository(t *testing.T, dialer pkgNet.Dialer) *portRepository {
	conf := InitConfig()
	repository := NewPortRepository(conf)

	systemPortRepository, ok := repository.(*portRepository)
	if assert.True(t, ok) {
		systemPortRepository.dialer = dialer
		return systemPortRepository
	}
	return nil
}

func TestRepository_OpenSocket_Success(t *testing.T) {
	mockConn := new(mocks.Conn)
	mockConn.On("Close").Return(nil)
	mockDialer := new(mocks.Dialer)
	mockDialer.On("Dial", AnythingOfType("string"), AnythingOfType("string")).Return(mockConn, nil)

	repository := initRepository(t, mockDialer)
	if repository != nil {
		assert.NoError(t, repository.OpenSocket(context.Background(), "test", 1234))
		mockConn.AssertNumberOfCalls(t, "Close", 1)
		mockConn.AssertExpectations(t)
		mockDialer.AssertNumberOfCalls(t, "Dial", 1)
		mockDialer.AssertExpectations(t)
	}
}

func TestRepository_OpenSocket_Failed(t *testing.T) {
	mockDialer := new(mocks.Dialer)
	mockDialer.On("Dial", AnythingOfType("string"), AnythingOfType("string")).Return(nil, errors.New("check port failed"))

	repository := initRepository(t, mockDialer)
	if repository != nil {
		assert.Error(t, repository.OpenSocket(context.Background(), "test", 1234))
		mockDialer.AssertNumberOfCalls(t, "Dial", 1)
		mockDialer.AssertExpectations(t)
	}
}
