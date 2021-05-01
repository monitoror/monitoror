package repository

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/monitoror/monitoror/monitorables/youtrack/api"
	"github.com/monitoror/monitoror/monitorables/youtrack/api/models"
	"github.com/monitoror/monitoror/monitorables/youtrack/config"
)

type (
	youtrackRepository struct {
		client *http.Client
		config *config.Youtrack
	}
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

func (yr *youtrackRepository) GetIssues(query string) (*models.Issues, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/api/issues", yr.config.URL), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", yr.config.Token))

	q := request.URL.Query()
	q.Add("query", query)
	request.URL.RawQuery = q.Encode()

	result, err := yr.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	issues := &models.Issues{}
	err = json.Unmarshal(body, issues)
	if err != nil {
		return nil, err
	}

	return issues, nil
}
