package repository

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/monitoror/monitoror/monitorable/config/models"

	"github.com/monitoror/monitoror/monitorable/config"
)

type (
	configRepository struct {
	}
)

func NewConfigRepository() config.Repository {
	return &configRepository{}
}

func GetConfig(reader io.Reader) (config *models.Config, err error) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return
	}

	return
}
