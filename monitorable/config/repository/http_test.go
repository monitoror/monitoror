package repository

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// /!\ this is an integration test /!\
// Note : It may be necessary to separate them from unit tests

// TestConfigRepository_GetConfigFromUrl test if http get works
func TestConfigRepository_GetConfigFromUrl(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, `{}`)
	}))
	defer ts.Close()

	repository := NewConfigRepository()
	_, err := repository.GetConfigFromUrl(ts.URL)
	assert.NoError(t, err)
}

// TestConfigRepository_GetConfigFromUrl test if http get works
func TestConfigRepository_GetConfigFromUrl_Error(t *testing.T) {
	repository := NewConfigRepository()
	_, err := repository.GetConfigFromUrl("http://monitoror.test")
	assert.Error(t, err)
}
