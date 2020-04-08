package delivery

import (
	"github.com/labstack/echo/v4"

	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
	"github.com/monitoror/monitoror/internal/pkg/monitorable/validator"
	coreModels "github.com/monitoror/monitoror/models"
)

func BindAndValidateRequestParams(ctx echo.Context, v uiConfigModels.ParamsValidator) error {
	if err := ctx.Bind(v); err != nil {
		return coreModels.ParamsError
	}

	if err := validator.Validate(v); err != nil {
		return err
	}

	return nil
}
