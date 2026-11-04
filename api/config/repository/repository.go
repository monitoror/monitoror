package repository

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/monitoror/monitoror/api/config"
	"github.com/monitoror/monitoror/api/config/models"
)

type (
	configRepository struct {
		httpClient *http.Client
	}
)

var ConfigHTTPTimeout = 10 * time.Second

func NewConfigRepository() config.Repository {
	//TODO: Add possibility to disable SSL check?
	return &configRepository{httpClient: &http.Client{Timeout: ConfigHTTPTimeout}}
}

func ReadConfig(reader io.Reader) (config *models.Config, err error) {
	bytes, err := io.ReadAll(reader)
	if err != nil {
		return
	}

	if err = json.Unmarshal(bytes, &config); err != nil {
		err = &models.ConfigUnmarshalError{Err: err, RawConfig: string(bytes)}
	}

	return
}
