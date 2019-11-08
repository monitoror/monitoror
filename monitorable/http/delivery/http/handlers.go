package http

import (
	netHttp "net/http"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/http"
	httpModels "github.com/monitoror/monitoror/monitorable/http/models"

	"github.com/labstack/echo/v4"
)

//nolint:golint
type HTTPDelivery struct {
	httpUsecase http.Usecase
}

func NewHTTPDelivery(p http.Usecase) *HTTPDelivery {
	return &HTTPDelivery{p}
}

func (h *HTTPDelivery) GetHTTPAny(c echo.Context) error {
	// Bind / Check Params
	params := &httpModels.HTTPAnyParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return models.QueryParamsError
	}

	tile, err := h.httpUsecase.HTTPAny(params)
	if err != nil {
		return err
	}

	return c.JSON(netHttp.StatusOK, tile)
}

func (h *HTTPDelivery) GetHTTPRaw(c echo.Context) error {
	// Bind / Check Params
	params := &httpModels.HTTPRawParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return models.QueryParamsError
	}

	tile, err := h.httpUsecase.HTTPRaw(params)
	if err != nil {
		return err
	}

	return c.JSON(netHttp.StatusOK, tile)
}

func (h *HTTPDelivery) GetHTTPJson(c echo.Context) error {
	// Bind / Check Params
	params := &httpModels.HTTPJsonParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return models.QueryParamsError
	}

	tile, err := h.httpUsecase.HTTPJson(params)
	if err != nil {
		return err
	}

	return c.JSON(netHttp.StatusOK, tile)
}

func (h *HTTPDelivery) GetHTTPYaml(c echo.Context) error {
	// Bind / Check Params
	params := &httpModels.HTTPYamlParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return models.QueryParamsError
	}

	tile, err := h.httpUsecase.HTTPYaml(params)
	if err != nil {
		return err
	}

	return c.JSON(netHttp.StatusOK, tile)
}
