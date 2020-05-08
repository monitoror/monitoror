package repository

import (
	"os"

	"github.com/monitoror/monitoror/api/config/models"
	"github.com/monitoror/monitoror/internal/pkg/path"
)

func (cr *configRepository) GetConfigFromPath(baseDir, filePath string) (config *models.Config, err error) {
	filePath = path.ToAbsolute(baseDir, filePath)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, &models.ConfigFileNotFoundError{Err: err, PathOrURL: filePath}
	}
	defer file.Close()

	config, err = ReadConfig(file)

	return
}
