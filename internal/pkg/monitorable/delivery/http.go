package delivery

import (
	"github.com/fatih/structs"
	"github.com/labstack/echo/v4"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/params"
	"github.com/monitoror/monitoror/internal/pkg/validator"
	"github.com/monitoror/monitoror/internal/pkg/validator/validate"
	coreModels "github.com/monitoror/monitoror/models"
	pkgStructs "github.com/monitoror/monitoror/pkg/structs"
)

func BindAndValidateParams(ctx echo.Context, p params.Validator) error {
	// Bind struct into query string using echo.Context.Bind
	if err := ctx.Bind(p); err != nil {
		return coreModels.ParamsError
	}

	var errors []validator.Error
	// use "validate" tag in struct definition to validate params
	errors = append(errors, validate.Struct(p)...)
	// use "Validate" function to to validate params
	errors = append(errors, p.Validate()...)

	if len(errors) > 0 {
		err := errors[0]
		// Lookup fin struct to find field and replace it by query tag value
		for _, field := range structs.Fields(p) {
			if field.Name() == err.GetFieldName() {
				// Replace FieldName By query FieldName
				if queryTagValue := pkgStructs.GetQueryFieldName(field); queryTagValue != "" {
					err.SetFieldName(queryTagValue)
				}
				break
			}
		}

		// Return Monitorable Error
		return &coreModels.MonitororError{Message: err.Error()}
	}

	return nil
}
