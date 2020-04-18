package validate

import (
	"fmt"

	pkgValidator "github.com/monitoror/monitoror/internal/pkg/validator"
)

type (
	validateError struct {
		errorID   pkgValidator.ErrorID
		fieldName string

		// tagParam used in Expected
		tagParam string
	}
)

func (e *validateError) Error() string {
	switch e.errorID {
	case pkgValidator.ErrorRequired:
		return fmt.Sprintf(`Required %q field is missing.`, e.fieldName)
	case pkgValidator.ErrorOneOf:
		return fmt.Sprintf(`Invalid %q field. Must be one of [%s].`, e.fieldName, e.tagParam)
	case pkgValidator.ErrorEq:
		return fmt.Sprintf(`Invalid %q field. Must be equal to %s.`, e.fieldName, e.tagParam)
	case pkgValidator.ErrorNE:
		return fmt.Sprintf(`Invalid %q field. Must be not equal to %s.`, e.fieldName, e.tagParam)
	case pkgValidator.ErrorGT:
		return fmt.Sprintf(`Invalid %q field. Must be greater than %s.`, e.fieldName, e.tagParam)
	case pkgValidator.ErrorGTE:
		return fmt.Sprintf(`Invalid %q field. Must be greater or equal to %s.`, e.fieldName, e.tagParam)
	case pkgValidator.ErrorLT:
		return fmt.Sprintf(`Invalid %q field. Must be lower than %s.`, e.fieldName, e.tagParam)
	case pkgValidator.ErrorLTE:
		return fmt.Sprintf(`Invalid %q field. Must be lower or equal to %s.`, e.fieldName, e.tagParam)
	case pkgValidator.ErrorNotEmpty:
		return fmt.Sprintf(`Invalid %q field. Must be not empty.`, e.fieldName)
	case pkgValidator.ErrorURL:
		return fmt.Sprintf(`Invalid %q field. Must be a valid URL.`, e.fieldName)
	case pkgValidator.ErrorHTTP:
		return fmt.Sprintf(`Invalid %q field. Must be start with "http://" or "https://".`, e.fieldName)
	case pkgValidator.ErrorRegex:
		return fmt.Sprintf(`Invalid %q field. Must be a valid golang regex.`, e.fieldName)
	default:
		return ""
	}
}

func (e *validateError) GetErrorID() pkgValidator.ErrorID { return e.errorID }
func (e *validateError) SetFieldName(f string)            { e.fieldName = f }
func (e *validateError) GetFieldName() string             { return e.fieldName }

func (e *validateError) Expected() string {
	switch e.errorID {
	case pkgValidator.ErrorOneOf:
		return e.tagParam
	case pkgValidator.ErrorEq:
		return fmt.Sprintf(`%s = %s`, e.fieldName, e.tagParam)
	case pkgValidator.ErrorNE:
		return fmt.Sprintf(`%s != %s`, e.fieldName, e.tagParam)
	case pkgValidator.ErrorGT:
		return fmt.Sprintf(`%s > %s`, e.fieldName, e.tagParam)
	case pkgValidator.ErrorGTE:
		return fmt.Sprintf(`%s >= %s`, e.fieldName, e.tagParam)
	case pkgValidator.ErrorLT:
		return fmt.Sprintf(`%s < %s`, e.fieldName, e.tagParam)
	case pkgValidator.ErrorLTE:
		return fmt.Sprintf(`%s <= %s`, e.fieldName, e.tagParam)
	default:
		return ""
	}
}
