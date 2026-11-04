package repository

import (
	"fmt"
	"net/http"

	"github.com/monitoror/monitoror/api/config/models"
)

func (cr *configRepository) GetConfigFromURL(url string) (config *models.Config, err error) {
	resp, err := cr.httpClient.Get(url)
	if err != nil {
		return nil, &models.ConfigFileNotFoundError{Err: err, PathOrURL: url}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &models.ConfigFileNotFoundError{Err: fmt.Errorf("unexpected status %d", resp.StatusCode), PathOrURL: url}
	}

	config, err = ReadConfig(resp.Body)

	return
}
