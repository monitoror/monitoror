package repository

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/test"

	"github.com/stretchr/testify/assert"
)

// /!\ this is an integration test /!\
// Note : It may be necessary to separate them from unit tests

// TestHTTPRepository_Get test if http get works
func TestHTTPRepository_Get(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, "Hello")
	}))
	defer ts.Close()

	repository := NewHTTPRepository(&config.HTTP{SSLVerify: false, Timeout: 2000})
	response, err := repository.Get(ts.URL)

	if assert.NoError(t, err) {
		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, "Hello", strings.TrimSpace(string(response.Body)))
	}
}

func TestHTTPRepository_Get_Error(t *testing.T) {
	repository := NewHTTPRepository(&config.HTTP{SSLVerify: false, Timeout: 2000})
	_, err := repository.Get("http://monitoror.example.com")
	assert.Error(t, err)
}

func TestHTTPRepository_Get_ReadAll_Error(t *testing.T) {
	client := test.NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(iotest.TimeoutReader(strings.NewReader("blabla"))), // Hacked reader to return error on ioutil.ReadAll
			Header:     make(http.Header),
		}
	})
	repository := httpRepository{httpClient: client}

	_, err := repository.Get("http://monitoror.example.com")
	assert.Error(t, err)
}
