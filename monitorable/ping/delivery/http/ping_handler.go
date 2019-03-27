package http

import (
	"net/http"

	"github.com/jsdidierlaurent/monitowall/monitorable/ping"

	"github.com/labstack/echo/v4"
)

type httpPingHandler struct {
	pingUsecase ping.Usecase
}

func NewHttpPingHandler(p ping.Usecase) *httpPingHandler {
	return &httpPingHandler{p}
}

func (h *httpPingHandler) GetPing(c echo.Context) error {
	tile, err := h.pingUsecase.Ping(c)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
