package repository

import (
	"errors"
	"testing"

	"github.com/monitoror/monitoror/monitorables/pingdom/config"
	pkgPingdom "github.com/monitoror/monitoror/pkg/gopingdom"
	"github.com/monitoror/monitoror/pkg/gopingdom/mocks"

	"github.com/jsdidierlaurent/go-pingdom/pingdom"
	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func initRepository(t *testing.T, checkAPI pkgPingdom.PingdomCheckAPI, transactionCheckAPI pkgPingdom.PingdomTransactionCheckAPI) *pingdomRepository {
	conf := &config.Pingdom{
		URL:             "https://pingdom.example.com/",
		Token:           config.Default.Token,
		Timeout:         config.Default.Timeout,
		CacheExpiration: config.Default.CacheExpiration,
	}
	repository := NewPingdomRepository(conf)

	assert.Equal(t, "https://pingdom.example.com", conf.URL)

	apiPingdomRepository, ok := repository.(*pingdomRepository)
	if assert.True(t, ok) {
		apiPingdomRepository.pingdomCheckAPI = checkAPI
		apiPingdomRepository.pingdomTransactionCheckAPI = transactionCheckAPI
		return apiPingdomRepository
	}
	return nil
}

func TestPingdomRepository_NewPingdomRepository_Error(t *testing.T) {
	conf := &config.Pingdom{
		URL:             "wrong%url",
		Token:           config.Default.Token,
		Timeout:         config.Default.Timeout,
		CacheExpiration: config.Default.CacheExpiration,
	}

	assert.Panics(t, func() { _ = NewPingdomRepository(conf) })
}

func TestPingdomRepository_GetPingdomCheck_Success(t *testing.T) {
	mock := new(mocks.PingdomCheckAPI)
	mock.On("Read", Anything).Return(&pingdom.CheckResponse{ID: 1000, Name: "Check 1", Status: "up"}, nil)

	repository := initRepository(t, mock, nil)
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

	repository := initRepository(t, mock, nil)
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

	repository := initRepository(t, mock, nil)
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

	repository := initRepository(t, mock, nil)
	_, err := repository.GetChecks("")
	assert.Error(t, err)
	mock.AssertNumberOfCalls(t, "List", 1)
	mock.AssertExpectations(t)
}

func TestPingdomRepository_GetPingdomTransactionCheck_Success(t *testing.T) {
	mock := new(mocks.PingdomTransactionCheckAPI)
	mock.On("Read", Anything).Return(&pingdom.TransactionCheckResponse{ID: 1000, Name: "Check 1", Status: "successful"}, nil)

	repository := initRepository(t, nil, mock)
	check, err := repository.GetTransactionCheck(1000)
	if assert.NoError(t, err) {
		assert.Equal(t, "Check 1", check.Name)
		assert.Equal(t, "successful", check.Status)
	}

	mock.AssertNumberOfCalls(t, "Read", 1)
	mock.AssertExpectations(t)
}

func TestPingdomRepository_GetPingdomTransactionCheck_Error(t *testing.T) {
	mock := new(mocks.PingdomTransactionCheckAPI)
	mock.On("Read", Anything).Return(nil, errors.New("boom"))

	repository := initRepository(t, nil, mock)
	_, err := repository.GetTransactionCheck(1000)
	assert.Error(t, err)
	mock.AssertNumberOfCalls(t, "Read", 1)
	mock.AssertExpectations(t)
}

func TestPingdomRepository_GetPingdomTransactionChecks_Success(t *testing.T) {
	mock := new(mocks.PingdomTransactionCheckAPI)
	mock.On("List", Anything).Return([]pingdom.TransactionCheckResponse{
		{ID: 1000, Name: "Check 1", Status: "successful"},
		{ID: 2000, Name: "Check 2", Status: "successful"},
		{ID: 3000, Name: "Check 3", Status: "failing"},
	}, nil)

	repository := initRepository(t, nil, mock)
	checks, err := repository.GetTransactionChecks("tests")
	if assert.NoError(t, err) {
		assert.Len(t, checks, 3)
	}

	mock.AssertNumberOfCalls(t, "List", 1)
	mock.AssertExpectations(t)
}

func TestPingdomRepository_GetPingdomTransactionChecks_Error(t *testing.T) {
	mock := new(mocks.PingdomTransactionCheckAPI)
	mock.On("List", Anything).Return(nil, errors.New("boom"))

	repository := initRepository(t, nil, mock)
	_, err := repository.GetTransactionChecks("")
	assert.Error(t, err)
	mock.AssertNumberOfCalls(t, "List", 1)
	mock.AssertExpectations(t)
}
