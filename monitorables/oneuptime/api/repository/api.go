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

	"github.com/monitoror/monitoror/monitorables/oneuptime/api"
	"github.com/monitoror/monitoror/monitorables/oneuptime/api/models"
	"github.com/monitoror/monitoror/monitorables/oneuptime/config"
)

type (
	oneuptimeRepository struct {
		client *http.Client
		config *config.Oneuptime
	}

	listResponse struct {
		Data  []monitor `json:"data"`
		Count int       `json:"count"`
	}

	monitor struct {
		ID                      string         `json:"_id"`
		Name                    string         `json:"name"`
		DisableActiveMonitoring bool           `json:"disableActiveMonitoring"`
		CurrentMonitorStatus    *monitorStatus `json:"currentMonitorStatus"`
		Labels                  []label        `json:"labels"`
	}

	monitorStatus struct {
		IsOperationalState bool `json:"isOperationalState"`
		IsOfflineState     bool `json:"isOfflineState"`
	}

	label struct {
		Name string `json:"name"`
	}
)

var PageSize = 100

// Only the fields used by Monitoror are requested to OneUptime
const listBody = `{"query":{},"select":{"_id":true,"name":true,"disableActiveMonitoring":true,` +
	`"currentMonitorStatus":{"isOperationalState":true,"isOfflineState":true},"labels":{"name":true}},` +
	`"sort":{"name":"Ascending"}}`

func NewOneuptimeRepository(config *config.Oneuptime) api.Repository {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !config.SSLVerify},
	}
	client := &http.Client{Transport: tr, Timeout: time.Duration(config.Timeout) * time.Millisecond}

	// Remove last /
	config.URL = strings.TrimRight(config.URL, "/")

	return &oneuptimeRepository{client, config}
}

func (r *oneuptimeRepository) GetMonitors() ([]models.Monitor, error) {
	var results []models.Monitor

	for skip := 0; ; skip += PageSize {
		page, err := r.getMonitorsPage(skip)
		if err != nil {
			return nil, err
		}

		for _, m := range page.Data {
			results = append(results, toMonitor(m))
		}

		if len(page.Data) < PageSize || len(results) >= page.Count {
			return results, nil
		}
	}
}

func (r *oneuptimeRepository) getMonitorsPage(skip int) (*listResponse, error) {
	url := fmt.Sprintf("%s/api/monitor/get-list?skip=%d&limit=%d", r.config.URL, skip, PageSize)
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(listBody)))
	if err != nil {
		return nil, err
	}
	request.Header.Set("ApiKey", r.config.Token)
	request.Header.Set("ProjectID", r.config.ProjectID)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := r.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("oneuptime returned status %d: %s", response.StatusCode, strings.TrimSpace(string(data)))
	}

	result := &listResponse{}
	if err = json.Unmarshal(data, result); err != nil {
		return nil, err
	}

	return result, nil
}

func toMonitor(m monitor) models.Monitor {
	result := models.Monitor{
		ID:         m.ID,
		Name:       m.Name,
		IsDisabled: m.DisableActiveMonitoring,
	}

	if m.CurrentMonitorStatus != nil {
		result.IsOperational = m.CurrentMonitorStatus.IsOperationalState
		result.IsOffline = m.CurrentMonitorStatus.IsOfflineState
	}

	for _, l := range m.Labels {
		result.Labels = append(result.Labels, l.Name)
	}

	return result
}
