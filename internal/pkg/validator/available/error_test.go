package available

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/monitoror/monitoror/internal/pkg/validator"
)

func TestAvailableError(t *testing.T) {
	for _, testcase := range []struct {
		err      *availableError
		message  string
		expected string
	}{
		{
			err:      &availableError{errorID: validator.ErrorSince, fieldName: "test", versionParam: "1.0"},
			message:  `"test" field is only available since version "1.0".`,
			expected: "version >= 1.0",
		},
		{
			err:      &availableError{errorID: validator.ErrorUntil, fieldName: "test", versionParam: "1.0"},
			message:  `"test" field is only available until version "1.0".`,
			expected: "version <= 1.0",
		},
		{
			err:      &availableError{errorID: 99999},
			message:  "",
			expected: "",
		},
	} {
		assert.Equal(t, testcase.message, testcase.err.Error())
		assert.Equal(t, testcase.expected, testcase.err.Expected())
	}
}

func TestAvailableError_SetFieldName(t *testing.T) {
	err := &availableError{fieldName: "TEST"}
	assert.Equal(t, "TEST", err.GetFieldName())
	err.SetFieldName("TEST2")
	assert.Equal(t, "TEST2", err.GetFieldName())
}
