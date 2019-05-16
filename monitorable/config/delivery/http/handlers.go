package http

import (
	"net/http"

	"github.com/monitoror/monitoror/models/errors"
	"github.com/monitoror/monitoror/monitorable/config"
	"github.com/monitoror/monitoror/monitorable/config/models"

	"github.com/labstack/echo/v4"
)

type httpConfigHandler struct {
	configUsecase config.Usecase
}

func NewHttpConfigHandler(cu config.Usecase) *httpConfigHandler {
	return &httpConfigHandler{cu}
}

func (h *httpConfigHandler) GetConfig(c echo.Context) error {
	// Bind / check Params
	params := &models.ConfigParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return errors.NewQueryParamsError(err)
	}

	config, err := h.configUsecase.Config(params)
	if err != nil {
		return err
	}

	if err = h.configUsecase.Verify(config); err != nil {
		return err
	}

	if err = h.configUsecase.Hydrate(config); err != nil {
		return err
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	return c.JSON(http.StatusOK, config)
}
