package repository

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/monitoror/monitoror/monitorables/youtrack/config"

	"github.com/stretchr/testify/assert"
)

func initRepository(t *testing.T, handler http.HandlerFunc) *youtrackRepository {
	CountRetryDelay = 0

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	conf := &config.Youtrack{
		URL:       server.URL + "/",
		Token:     "TOKEN",
		Timeout:   config.Default.Timeout,
		SSLVerify: true,
	}
	repository := NewYoutrackRepository(conf)
	assert.Equal(t, server.URL, conf.URL)

	return repository.(*youtrackRepository)
}

func TestYoutrackRepository_CountIssues_Success(t *testing.T) {
	repository := initRepository(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/issuesGetter/count", r.URL.Path)
		assert.Equal(t, "count", r.URL.Query().Get("fields"))
		assert.Equal(t, "Bearer TOKEN", r.Header.Get("Authorization"))

		body := &countRequest{}
		assert.NoError(t, json.NewDecoder(r.Body).Decode(body))
		assert.Equal(t, "project: MO", body.Query)

		_, _ = fmt.Fprint(w, `{"count":42,"$type":"IssueCountResponse"}`)
	})

	count, err := repository.CountIssues("project: MO")
	if assert.NoError(t, err) {
		assert.Equal(t, 42, count)
	}
}

func TestYoutrackRepository_CountIssues_RetryWhileComputing(t *testing.T) {
	calls := 0
	repository := initRepository(t, func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls < 3 {
			_, _ = fmt.Fprint(w, `{"count":-1}`)
			return
		}
		_, _ = fmt.Fprint(w, `{"count":7}`)
	})

	count, err := repository.CountIssues("project: MO")
	if assert.NoError(t, err) {
		assert.Equal(t, 7, count)
		assert.Equal(t, 3, calls)
	}
}

func TestYoutrackRepository_CountIssues_StillComputing(t *testing.T) {
	calls := 0
	repository := initRepository(t, func(w http.ResponseWriter, r *http.Request) {
		calls++
		_, _ = fmt.Fprint(w, `{"count":-1}`)
	})

	_, err := repository.CountIssues("project: MO")
	assert.Error(t, err)
	assert.Equal(t, CountMaxAttempts, calls)
}

func TestYoutrackRepository_CountIssues_HTTPError(t *testing.T) {
	repository := initRepository(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprint(w, `{"error":"Unauthorized"}`)
	})

	_, err := repository.CountIssues("project: MO")
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "401")
	}
}

func TestYoutrackRepository_CountIssues_InvalidJSON(t *testing.T) {
	repository := initRepository(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `not json`)
	})

	_, err := repository.CountIssues("project: MO")
	assert.Error(t, err)
}

func TestYoutrackRepository_CountIssues_ConnectionError(t *testing.T) {
	repository := initRepository(t, func(w http.ResponseWriter, r *http.Request) {})
	repository.config.URL = "http://127.0.0.1:0"

	_, err := repository.CountIssues("project: MO")
	assert.Error(t, err)
}

func TestYoutrackRepository_CountIssues_InvalidURL(t *testing.T) {
	repository := initRepository(t, func(w http.ResponseWriter, r *http.Request) {})
	repository.config.URL = "http://%%"

	_, err := repository.CountIssues("project: MO")
	assert.Error(t, err)
}

func TestYoutrackRepository_CountIssues_TruncatedBody(t *testing.T) {
	repository := initRepository(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "100")
		_, _ = fmt.Fprint(w, `{"count":`)
	})

	_, err := repository.CountIssues("project: MO")
	assert.Error(t, err)
}
