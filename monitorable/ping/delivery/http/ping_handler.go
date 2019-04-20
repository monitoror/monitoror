package http

import (
	"net/http"

	"github.com/monitoror/monitoror/models/errors"
	"github.com/monitoror/monitoror/monitorable/ping/model"

	"github.com/monitoror/monitoror/monitorable/ping"

	"github.com/labstack/echo/v4"
)

type HttpPingHandler struct {
	pingUsecase ping.Usecase
}

func NewHttpPingHandler(p ping.Usecase) *HttpPingHandler {
	return &HttpPingHandler{p}
}

func (h *HttpPingHandler) GetPing(c echo.Context) error {
	// Bind / Validate Params
	params := &model.PingParams{}
	err := c.Bind(params)
	if err != nil || !params.Validate() {
		return errors.NewQueryParamsError(err)
	}

	tile, err := h.pingUsecase.Ping(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
