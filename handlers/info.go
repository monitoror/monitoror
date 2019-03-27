package handlers

import (
	"net/http"

	"github.com/jsdidierlaurent/monitowall/configs"
	"github.com/jsdidierlaurent/monitowall/models"

	"github.com/labstack/echo/v4"
)

type httpInfoHandler struct {
	buildInfo *configs.BuildInfo
	config    *configs.Config
}

func HttpInfoHandler(buildInfo *configs.BuildInfo, config *configs.Config) *httpInfoHandler {
	return &httpInfoHandler{buildInfo, config}
}

func (h *httpInfoHandler) GetInfo(c echo.Context) error {
	response := models.InfoResponse{
		BuildInfo: *h.buildInfo,
		Config:    *h.config,
	}
	return c.JSON(http.StatusOK, response)
}
