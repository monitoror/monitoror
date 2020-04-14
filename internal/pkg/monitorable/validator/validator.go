package validator

import (
	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
	uiCoreVersions "github.com/monitoror/monitoror/api/config/versions"
	coreModels "github.com/monitoror/monitoror/models"
)

var configVersion = uiConfigModels.ParseVersion(uiCoreVersions.CurrentVersion)

func Validate(v uiConfigModels.ParamsValidator) error {
	if err := v.Validate(configVersion); err != nil {
		return &coreModels.MonitororError{Message: err.Message}
	}

	return nil
}
