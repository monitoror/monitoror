package errors

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	. "github.com/monitoror/monitoror/models/tiles"
)

type NoBuildError struct {
	BuildTile *BuildTile
}

func NewNoBuildError(buildTile *BuildTile) *NoBuildError {
	return &NoBuildError{buildTile}
}

func (nbe *NoBuildError) Error() string {
	return fmt.Sprintf("unable to found build")
}

func (nbe *NoBuildError) Send(ctx echo.Context) {
	tile := nbe.BuildTile
	tile.Status = WarningStatus
	tile.Message = nbe.Error()

	_ = ctx.JSON(http.StatusOK, tile)
}
