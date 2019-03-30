package errors

import (
	"fmt"
)

type QueryParamsError struct {
	err error
}

func NewQueryParamsError(err error) *QueryParamsError {
	return &QueryParamsError{err}
}

func (qpe *QueryParamsError) Error() string {
	if qpe.err != nil {
		return fmt.Sprintf("unable to parse/validate queryParams into struct, %v", qpe.err)
	} else {
		return fmt.Sprintf("unable to parse/validate queryParams into struct")
	}
}
