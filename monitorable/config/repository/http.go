package repository

import (
	"github.com/monitoror/monitoror/monitorable/config/models"
)

func (cr *configRepository) GetConfigFromUrl(url string) (config *models.Config, err error) {
	resp, err := cr.httpClient.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	config, err = ReadConfig(resp.Body)
	return
}
