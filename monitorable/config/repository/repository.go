package repository

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/monitoror/monitoror/monitorable/config/models"

	"github.com/monitoror/monitoror/monitorable/config"
)

type (
	configRepository struct {
		httpClient *http.Client
	}
)

func NewConfigRepository() config.Repository {
	//TODO: Add possibility to disable SSL check?
	return &configRepository{httpClient: http.DefaultClient}
}

func ReadConfig(reader io.Reader) (config *models.Config, err error) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return
	}

	if err = json.Unmarshal(bytes, &config); err != nil {
		return
	}

	return
}
