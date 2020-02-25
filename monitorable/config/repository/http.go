package repository

import (
	"net/http"

	"github.com/monitoror/monitoror/monitorable/config/models"
)

func (cr *configRepository) GetConfigFromURL(url string) (config *models.Config, err error) {
	resp, err := cr.httpClient.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, ErrConfigFileNotFound
	}
	defer resp.Body.Close()

	config, err = ReadConfig(resp.Body)
	return
}
