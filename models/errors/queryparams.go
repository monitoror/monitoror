package errors

import (
	"fmt"

	. "github.com/jsdidierlaurent/monitowall/models/tiles"
)

type QueryParamsError struct {
	Tile *Tile
	err  error
}

func NewQueryParamsError(tile *Tile, err error) *QueryParamsError {
	return &QueryParamsError{tile, err}
}

func (qpe *QueryParamsError) Error() string {
	if qpe.err != nil {
		return fmt.Sprintf("unable to parse/validate queryParams into struct, %v", qpe.err)
	} else {
		return fmt.Sprintf("unable to parse/validate queryParams into struct")
	}
}
