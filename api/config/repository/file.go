package repository

import (
	"errors"
	"os"

	"github.com/monitoror/monitoror/api/config/models"
)

func (cr *configRepository) GetConfigFromPath(path string) (config *models.Config, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, &models.ConfigFileNotFoundError{Err: err, PathOrURL: path}
	}
	defer file.Close()

	config, err = ReadConfig(file)

	// Remove RawConfig by security on GetConfigFromPath.
	// This can be leak files if monitoror as to high right on system.
	// TODO: Remove this when directory traversal will be fix: https://github.com/monitoror/monitoror/issues/222
	var cue *models.ConfigUnmarshalError
	if errors.As(err, &cue) {
		cue.RawConfig = ""
	}

	return
}
