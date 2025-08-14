package repository

import (
	"crypto/tls"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/monitoror/monitoror/monitorables/http/api"
	"github.com/monitoror/monitoror/monitorables/http/api/models"
	"github.com/monitoror/monitoror/monitorables/http/config"
)

type (
	httpRepository struct {
		httpClient *http.Client
	}
	clientInterceptor struct {
		Transport http.RoundTripper // A custom RoundTripper that adds headers to each request.
		Headers   []string
	}
)

func NewHTTPRepository(config *config.HTTP) api.Repository {
	var certificates []tls.Certificate

	if config.Certificate != "" && config.Key != "" {
		cert, error := tls.LoadX509KeyPair(config.Certificate, config.Key)

		if error == nil {
			certificates = append(certificates, cert)
		}
	}

	bt := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !config.SSLVerify, Certificates: certificates},
	}

	// Intercept and add headers
	tr := &clientInterceptor{bt, config.Headers}

	client := &http.Client{Transport: tr, Timeout: time.Duration(config.Timeout) * time.Millisecond}

	return &httpRepository{client}
}

func (r *httpRepository) Get(url string) (response *models.Response, err error) {
	resp, err := r.httpClient.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	response = &models.Response{
		StatusCode: resp.StatusCode,
		Body:       bytes,
	}

	return
}

// RoundTrip implements the http.RoundTripper interface.
//
// It adds the headers from Headers to the request, then calls the original
// RoundTripper to make the request.
func (ci *clientInterceptor) RoundTrip(req *http.Request) (*http.Response, error) {

	headers := append(config.Default.Headers, ci.Headers...)

	// Add headers to the request
	for _, header := range headers {

		parts := strings.SplitN(header, ":", 2)
		if len(parts) == 2 {
			req.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		}
	}

	// Call the original RoundTripper
	return ci.Transport.RoundTrip(req)
}
