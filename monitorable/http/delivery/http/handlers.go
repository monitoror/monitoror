package http

import (
	"net/http"

	. "github.com/monitoror/monitoror/models"

	ping "github.com/monitoror/monitoror/monitorable/http"
	"github.com/monitoror/monitoror/monitorable/http/models"

	"github.com/labstack/echo/v4"
)

type httpHttpDelivery struct {
	httpUsecase ping.Usecase
}

func NewHttpHttpDelivery(p ping.Usecase) *httpHttpDelivery {
	return &httpHttpDelivery{p}
}

func (h *httpHttpDelivery) GetHttpAny(c echo.Context) error {
	// Bind / Check Params
	params := &models.HttpAnyParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return QueryParamsError
	}

	tile, err := h.httpUsecase.HttpAny(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}

func (h *httpHttpDelivery) GetHttpRaw(c echo.Context) error {
	// Bind / Check Params
	params := &models.HttpRawParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return QueryParamsError
	}

	tile, err := h.httpUsecase.HttpRaw(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}

func (h *httpHttpDelivery) GetHttpJson(c echo.Context) error {
	// Bind / Check Params
	params := &models.HttpFormattedDataParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return QueryParamsError
	}

	tile, err := h.httpUsecase.HttpJson(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}

func (h *httpHttpDelivery) GetHttpYaml(c echo.Context) error {
	// Bind / Check Params
	params := &models.HttpFormattedDataParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return QueryParamsError
	}

	tile, err := h.httpUsecase.HttpYaml(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
