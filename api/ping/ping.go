package ping

import (
	"net/http"

	"github.com/jsdidierlaurent/monitowall/model"
	"github.com/labstack/echo/v4"
)

type handler struct {
	PingModel model.PingModelImpl
}

func NewHandler(p model.PingModelImpl) *handler {
	return &handler{p}
}

func (h *handler) GetPing(c echo.Context) (err error) {
	hostname := c.QueryParam("hostname")
	return c.JSON(http.StatusOK, h.PingModel.Ping(hostname))
}
