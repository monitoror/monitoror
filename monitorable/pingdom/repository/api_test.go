package repository

import (
	"errors"
	"testing"

	"github.com/jsdidierlaurent/go-pingdom/pingdom"
	. "github.com/monitoror/monitoror/config"
	pkgPingdom "github.com/monitoror/monitoror/pkg/gopingdom"
	"github.com/monitoror/monitoror/pkg/gopingdom/mocks"

	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func initRepository(t *testing.T, checkAPI pkgPingdom.PingdomCheckAPI) *pingdomRepository {
	conf := InitConfig()
	repository := NewPingdomRepository(conf.Monitorable.Pingdom[DefaultVariant])

	apiPingdomRepository, ok := repository.(*pingdomRepository)
	if assert.True(t, ok) {
		apiPingdomRepository.pingdomCheckAPI = checkAPI
		return apiPingdomRepository
	}
	return nil
}

func TestPingdomRepository_NewPingdomRepository_Error(t *testing.T) {
	conf := InitConfig()
	conf.Monitorable.Pingdom[DefaultVariant].URL = "wrong%url"

	assert.Panics(t, func() { _ = NewPingdomRepository(conf.Monitorable.Pingdom[DefaultVariant]) })
}

func TestPingdomRepository_GetPingdomCheck_Success(t *testing.T) {
	mock := new(mocks.PingdomCheckAPI)
	mock.On("Read", Anything).Return(&pingdom.CheckResponse{ID: 1000, Name: "Check 1", Status: "up"}, nil)

	repository := initRepository(t, mock)
	check, err := repository.GetCheck(1000)
	if assert.NoError(t, err) {
		assert.Equal(t, "Check 1", check.Name)
		assert.Equal(t, "up", check.Status)
	}

	mock.AssertNumberOfCalls(t, "Read", 1)
	mock.AssertExpectations(t)
}

func TestPingdomRepository_GetPingdomCheck_Error(t *testing.T) {
	mock := new(mocks.PingdomCheckAPI)
	mock.On("Read", Anything).Return(nil, errors.New("boom"))

	repository := initRepository(t, mock)
	_, err := repository.GetCheck(1000)
	assert.Error(t, err)
	mock.AssertNumberOfCalls(t, "Read", 1)
	mock.AssertExpectations(t)
}

func TestPingdomRepository_GetPingdomChecks_Success(t *testing.T) {
	mock := new(mocks.PingdomCheckAPI)
	mock.On("List", Anything).Return([]pingdom.CheckResponse{
		{ID: 1000, Name: "Check 1", Status: "up"},
		{ID: 2000, Name: "Check 2", Status: "up"},
		{ID: 3000, Name: "Check 3", Status: "down"},
	}, nil)

	repository := initRepository(t, mock)
	checks, err := repository.GetChecks("tests")
	if assert.NoError(t, err) {
		assert.Len(t, checks, 3)
	}

	mock.AssertNumberOfCalls(t, "List", 1)
	mock.AssertExpectations(t)
}

func TestPingdomRepository_GetPingdomChecks_Error(t *testing.T) {
	mock := new(mocks.PingdomCheckAPI)
	mock.On("List", Anything).Return(nil, errors.New("boom"))

	repository := initRepository(t, mock)
	_, err := repository.GetChecks("")
	assert.Error(t, err)
	mock.AssertNumberOfCalls(t, "List", 1)
	mock.AssertExpectations(t)
}
