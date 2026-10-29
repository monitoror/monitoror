package repository

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/monitoror/monitoror/monitorables/oneuptime/api/models"
	"github.com/monitoror/monitoror/monitorables/oneuptime/config"

	"github.com/stretchr/testify/assert"
)

func initRepository(t *testing.T, handler http.HandlerFunc) *oneuptimeRepository {
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	conf := &config.Oneuptime{
		URL:       server.URL + "/",
		ProjectID: "PROJECT",
		Token:     "TOKEN",
		Timeout:   config.Default.Timeout,
		SSLVerify: true,
	}
	repository := NewOneuptimeRepository(conf)
	assert.Equal(t, server.URL, conf.URL)

	return repository.(*oneuptimeRepository)
}

func TestOneuptimeRepository_GetMonitors_Success(t *testing.T) {
	repository := initRepository(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/monitor/get-list", r.URL.Path)
		assert.Equal(t, "TOKEN", r.Header.Get("ApiKey"))
		assert.Equal(t, "PROJECT", r.Header.Get("ProjectID"))

		body, _ := io.ReadAll(r.Body)
		assert.Equal(t, listBody, string(body))

		_, _ = fmt.Fprint(w, `{"data":[
			{"_id":"1","name":"Monitor 1","disableActiveMonitoring":false,"currentMonitorStatus":{"isOperationalState":true,"isOfflineState":false},"labels":[{"name":"label1"},{"name":"label2"}]},
			{"_id":"2","name":"Monitor 2","disableActiveMonitoring":true,"currentMonitorStatus":{"isOperationalState":false,"isOfflineState":true},"labels":[]},
			{"_id":"3","name":"Monitor 3","disableActiveMonitoring":false}
		],"count":3,"skip":0,"limit":100}`)
	})

	monitors, err := repository.GetMonitors()
	if assert.NoError(t, err) {
		assert.Equal(t, []models.Monitor{
			{ID: "1", Name: "Monitor 1", IsOperational: true, Labels: []string{"label1", "label2"}},
			{ID: "2", Name: "Monitor 2", IsDisabled: true, IsOffline: true},
			{ID: "3", Name: "Monitor 3"},
		}, monitors)
	}
}

func TestOneuptimeRepository_GetMonitors_Pagination(t *testing.T) {
	PageSize = 2
	t.Cleanup(func() { PageSize = 100 })

	var skips []string
	repository := initRepository(t, func(w http.ResponseWriter, r *http.Request) {
		skip := r.URL.Query().Get("skip")
		skips = append(skips, skip)
		assert.Equal(t, "2", r.URL.Query().Get("limit"))

		switch skip {
		case "0":
			_, _ = fmt.Fprint(w, `{"data":[{"_id":"1"},{"_id":"2"}],"count":3}`)
		default:
			_, _ = fmt.Fprint(w, `{"data":[{"_id":"3"}],"count":3}`)
		}
	})

	monitors, err := repository.GetMonitors()
	if assert.NoError(t, err) {
		assert.Len(t, monitors, 3)
		assert.Equal(t, []string{"0", "2"}, skips)
	}
}

func TestOneuptimeRepository_GetMonitors_StopOnCount(t *testing.T) {
	PageSize = 2
	t.Cleanup(func() { PageSize = 100 })

	calls := 0
	repository := initRepository(t, func(w http.ResponseWriter, r *http.Request) {
		calls++
		_, _ = fmt.Fprint(w, `{"data":[{"_id":"1"},{"_id":"2"}],"count":2}`)
	})

	monitors, err := repository.GetMonitors()
	if assert.NoError(t, err) {
		assert.Len(t, monitors, 2)
		assert.Equal(t, 1, calls)
	}
}

func TestOneuptimeRepository_GetMonitors_HTTPError(t *testing.T) {
	repository := initRepository(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprint(w, `{"error":"Invalid API Key"}`)
	})

	_, err := repository.GetMonitors()
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "401")
	}
}

func TestOneuptimeRepository_GetMonitors_InvalidJSON(t *testing.T) {
	repository := initRepository(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `not json`)
	})

	_, err := repository.GetMonitors()
	assert.Error(t, err)
}

func TestOneuptimeRepository_GetMonitors_TruncatedBody(t *testing.T) {
	repository := initRepository(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "100")
		_, _ = fmt.Fprint(w, `{"data":`)
	})

	_, err := repository.GetMonitors()
	assert.Error(t, err)
}

func TestOneuptimeRepository_GetMonitors_InvalidURL(t *testing.T) {
	repository := initRepository(t, func(w http.ResponseWriter, r *http.Request) {})
	repository.config.URL = "http://%%"

	_, err := repository.GetMonitors()
	assert.Error(t, err)
}

func TestOneuptimeRepository_GetMonitors_ConnectionError(t *testing.T) {
	repository := initRepository(t, func(w http.ResponseWriter, r *http.Request) {})
	repository.config.URL = "http://127.0.0.1:0"

	_, err := repository.GetMonitors()
	assert.Error(t, err)
}
