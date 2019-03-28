package handlers

import (
	"net/http"

	"github.com/jsdidierlaurent/monitowall/cli/version"
	"github.com/jsdidierlaurent/monitowall/config"
	"github.com/jsdidierlaurent/monitowall/models"

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
