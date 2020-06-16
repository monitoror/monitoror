package repository

import (
	"net/http"

	"github.com/monitoror/monitoror/api/config/models"
)

func (cr *configRepository) GetConfigFromURL(url string) (config *models.Config, err error) {
	resp, err := cr.httpClient.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, &models.ConfigFileNotFoundError{Err: err, PathOrURL: url}
	}
	defer resp.Body.Close()

	config, err = ReadConfig(resp.Body)

	return
}
