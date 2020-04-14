package http

import (
	netHttp "net/http"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/delivery"
	"github.com/monitoror/monitoror/monitorables/http/api"
	"github.com/monitoror/monitoror/monitorables/http/api/models"

	"github.com/labstack/echo/v4"
)

//nolint:golint
type HTTPDelivery struct {
	httpUsecase api.Usecase
}

func NewHTTPDelivery(p api.Usecase) *HTTPDelivery {
	return &HTTPDelivery{p}
}

func (h *HTTPDelivery) GetHTTPStatus(c echo.Context) error {
	// Bind / Check Params
	params := &models.HTTPStatusParams{}
	if err := delivery.BindAndValidateRequestParams(c, params); err != nil {
		return err
	}

	tile, err := h.httpUsecase.HTTPStatus(params)
	if err != nil {
		return err
	}

	return c.JSON(netHttp.StatusOK, tile)
}

func (h *HTTPDelivery) GetHTTPRaw(c echo.Context) error {
	// Bind / Check Params
	params := &models.HTTPRawParams{}
	if err := delivery.BindAndValidateRequestParams(c, params); err != nil {
		return err
	}

	tile, err := h.httpUsecase.HTTPRaw(params)
	if err != nil {
		return err
	}

	return c.JSON(netHttp.StatusOK, tile)
}

func (h *HTTPDelivery) GetHTTPFormatted(c echo.Context) error {
	// Bind / Check Params
	params := &models.HTTPFormattedParams{}
	if err := delivery.BindAndValidateRequestParams(c, params); err != nil {
		return err
	}

	tile, err := h.httpUsecase.HTTPFormatted(params)
	if err != nil {
		return err
	}

	return c.JSON(netHttp.StatusOK, tile)
}
