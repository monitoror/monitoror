package http

import (
	"net/http"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/delivery"
	"github.com/monitoror/monitoror/monitorables/youtrack/api"
	"github.com/monitoror/monitoror/monitorables/youtrack/api/models"

	"github.com/labstack/echo/v4"
)

type YoutrackDelivery struct {
	youtrackUsecase api.Usecase
}

func NewYoutrackDelivery(p api.Usecase) *YoutrackDelivery {
	return &YoutrackDelivery{p}
}

func (gd *YoutrackDelivery) GetCountIssues(c echo.Context) error {
	// Bind / check Params
	params := &models.IssuesCountParams{}
	if err := delivery.BindAndValidateParams(c, params); err != nil {
		return err
	}

	tile, err := gd.youtrackUsecase.CountIssues(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
