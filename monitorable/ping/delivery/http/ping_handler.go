package http

import (
	"net/http"

	"github.com/jsdidierlaurent/monitoror/models/errors"
	"github.com/jsdidierlaurent/monitoror/monitorable/ping/model"

	"github.com/jsdidierlaurent/monitoror/monitorable/ping"

	"github.com/labstack/echo/v4"
)

type httpPingHandler struct {
	pingUsecase ping.Usecase
}

func NewHttpPingHandler(p ping.Usecase) *httpPingHandler {
	return &httpPingHandler{p}
}

func (h *httpPingHandler) GetPing(c echo.Context) error {
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
