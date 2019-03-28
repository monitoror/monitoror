package handlers

import (
	"net/http"

	"github.com/jsdidierlaurent/monitoror/cli/version"
	"github.com/jsdidierlaurent/monitoror/config"
	"github.com/jsdidierlaurent/monitoror/models"

	"github.com/labstack/echo/v4"
)

type httpInfoHandler struct {
	config *config.Config
}

func HttpInfoHandler(config *config.Config) *httpInfoHandler {
	return &httpInfoHandler{config}
}

func (h *httpInfoHandler) GetInfo(c echo.Context) error {
	response := models.NewInfoResponse(version.Version, version.GitCommit, version.BuildTime, h.config)
	return c.JSON(http.StatusOK, response)
}
