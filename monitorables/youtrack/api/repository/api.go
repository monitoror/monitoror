package repository

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/monitoror/monitoror/monitorables/youtrack/api"
	"github.com/monitoror/monitoror/monitorables/youtrack/config"
)

type (
	youtrackRepository struct {
		client *http.Client
		config *config.Youtrack
	}

	countRequest struct {
		Query string `json:"query"`
	}

	countResponse struct {
		Count int `json:"count"`
	}
)

// YouTrack answers -1 while the count is still being computed, so the request is retried.
var (
	CountMaxAttempts = 5
	CountRetryDelay  = 200 * time.Millisecond
)

func NewYoutrackRepository(config *config.Youtrack) api.Repository {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !config.SSLVerify},
	}
	client := &http.Client{Transport: tr, Timeout: time.Duration(config.Timeout) * time.Millisecond}

	// Remove last /
	config.URL = strings.TrimRight(config.URL, "/")

	return &youtrackRepository{client, config}
}

func (yr *youtrackRepository) CountIssues(query string) (int, error) {
	// Marshalling a struct of strings cannot fail
	body, _ := json.Marshal(countRequest{Query: query})

	for attempt := 1; ; attempt++ {
		count, err := yr.requestCount(body)
		if err != nil {
			return 0, err
		}

		if count >= 0 {
			return count, nil
		}

		if attempt >= CountMaxAttempts {
			return 0, fmt.Errorf("youtrack is still computing issues count after %d attempts", attempt)
		}
		time.Sleep(CountRetryDelay)
	}
}

func (yr *youtrackRepository) requestCount(body []byte) (int, error) {
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/issuesGetter/count?fields=count", yr.config.URL), bytes.NewReader(body))
	if err != nil {
		return 0, err
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", yr.config.Token))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := yr.client.Do(request)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}

	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("youtrack returned status %d: %s", response.StatusCode, strings.TrimSpace(string(data)))
	}

	result := &countResponse{}
	if err = json.Unmarshal(data, result); err != nil {
		return 0, err
	}

	return result.Count, nil
}
