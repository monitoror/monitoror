package errors

import (
	"fmt"

	. "github.com/monitoror/monitoror/models/tiles"
)

type NoBuildError struct {
	BuildTile *BuildTile
}

func NewNoBuildError(buildTile *BuildTile) *NoBuildError {
	return &NoBuildError{buildTile}
}

func (nbe *NoBuildError) Error() string {
	return fmt.Sprintf("no build found for %s", nbe.BuildTile.Label)
}
