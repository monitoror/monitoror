package available

import (
	"fmt"

	pkgValidator "github.com/monitoror/monitoror/internal/pkg/validator"
)

type (
	availableError struct {
		errorID   pkgValidator.ErrorID
		fieldName string

		// versionParam used in Expected
		versionParam string
	}
)

func (e *availableError) Error() string {
	switch e.errorID {
	case pkgValidator.ErrorSince:
		return fmt.Sprintf(`%q field is only available since version %q.`, e.fieldName, e.versionParam)
	case pkgValidator.ErrorUntil:
		return fmt.Sprintf(`%q field is only available until version %q.`, e.fieldName, e.versionParam)
	default:
		return ""
	}
}

func (e *availableError) GetErrorID() pkgValidator.ErrorID { return e.errorID }
func (e *availableError) SetFieldName(f string)            { e.fieldName = f }
func (e *availableError) GetFieldName() string             { return e.fieldName }

func (e *availableError) Expected() string {
	switch e.errorID {
	case pkgValidator.ErrorSince:
		return fmt.Sprintf(`version >= %s`, e.versionParam)
	case pkgValidator.ErrorUntil:
		return fmt.Sprintf(`version <= %s`, e.versionParam)
	default:
		return ""
	}
}
