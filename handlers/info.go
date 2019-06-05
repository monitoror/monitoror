package handlers

import (
	"net/http"

	"github.com/monitoror/monitoror/cli/version"
	"github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/models"

	"github.com/labstack/echo/v4"
)

type httpInfoDelivery struct {
	config *config.Config
}

func NewHttpInfoDelivery(config *config.Config) *httpInfoDelivery {
	return &httpInfoDelivery{config}
}

func (h *httpInfoDelivery) GetInfo(c echo.Context) error {
	response := models.NewInfoResponse(version.Version, version.GitCommit, version.BuildTime, h.config)
	return c.JSON(http.StatusOK, response)
}
