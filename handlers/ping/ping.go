package ping

import (
	"net/http"

	. "github.com/jsdidierlaurent/monitowall/renderings"

	"github.com/jsdidierlaurent/monitowall/models"
	"github.com/labstack/echo/v4"
)

type handler struct {
	PingModel models.PingModel
}

func NewHandler(p models.PingModel) *handler {
	return &handler{p}
}

func newResponse() (response *HealthCheckResponse) {
	response = NewHealthCheckResponse()
	response.Type = TypePing
	return
}

var count = 0

func (h *handler) GetPing(c echo.Context) (err error) {
	// Query Param
	hostname := c.QueryParam("hostname")

	// Response
	response := newResponse()
	response.Label = hostname
	
	// Execute Ping and return result
	message, err := h.PingModel.Ping(hostname)
	if err != nil {
		response.Status = FailStatus
	} else {
		response.Status = SuccessStatus
		response.Message = message
	}
	return c.JSON(http.StatusOK, response)
}
