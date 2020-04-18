package validate

import (
	"testing"

	"github.com/monitoror/monitoror/internal/pkg/validator"

	"github.com/stretchr/testify/assert"
)

func TestValidateError(t *testing.T) {
	for _, testcase := range []struct {
		err      *validateError
		message  string
		expected string
	}{
		{
			err:      &validateError{errorID: validator.ErrorRequired, fieldName: "test"},
			message:  `Required "test" field is missing.`,
			expected: "",
		},
		{
			err:      &validateError{errorID: validator.ErrorOneOf, fieldName: "test", tagParam: "test"},
			message:  `Invalid "test" field. Must be one of [test].`,
			expected: `test`,
		},
		{
			err:      &validateError{errorID: validator.ErrorEq, fieldName: "test", tagParam: "1"},
			message:  `Invalid "test" field. Must be equal to 1.`,
			expected: `test = 1`,
		},
		{
			err:      &validateError{errorID: validator.ErrorNE, fieldName: "test", tagParam: "1"},
			message:  `Invalid "test" field. Must be not equal to 1.`,
			expected: `test != 1`,
		},
		{
			err:      &validateError{errorID: validator.ErrorGT, fieldName: "test", tagParam: "1"},
			message:  `Invalid "test" field. Must be greater than 1.`,
			expected: `test > 1`,
		},
		{
			err:      &validateError{errorID: validator.ErrorGTE, fieldName: "test", tagParam: "1"},
			message:  `Invalid "test" field. Must be greater or equal to 1.`,
			expected: `test >= 1`,
		},
		{
			err:      &validateError{errorID: validator.ErrorLT, fieldName: "test", tagParam: "1"},
			message:  `Invalid "test" field. Must be lower than 1.`,
			expected: `test < 1`,
		},
		{
			err:      &validateError{errorID: validator.ErrorLTE, fieldName: "test", tagParam: "1"},
			message:  `Invalid "test" field. Must be lower or equal to 1.`,
			expected: `test <= 1`,
		},
		{
			err:      &validateError{errorID: validator.ErrorNotEmpty, fieldName: "test"},
			message:  `Invalid "test" field. Must be not empty.`,
			expected: ``,
		},
		{
			err:      &validateError{errorID: validator.ErrorURL, fieldName: "test"},
			message:  `Invalid "test" field. Must be a valid URL.`,
			expected: "",
		},
		{
			err:      &validateError{errorID: validator.ErrorHTTP, fieldName: "test"},
			message:  `Invalid "test" field. Must be start with "http://" or "https://".`,
			expected: "",
		},
		{
			err:      &validateError{errorID: validator.ErrorRegex, fieldName: "test"},
			message:  `Invalid "test" field. Must be a valid golang regex.`,
			expected: "",
		},
		{
			err:      &validateError{errorID: 99999},
			message:  "",
			expected: "",
		},
	} {
		assert.Equal(t, testcase.message, testcase.err.Error())
		assert.Equal(t, testcase.expected, testcase.err.Expected())
	}
}

func TestValidateError_SetFieldName(t *testing.T) {
	err := &validateError{fieldName: "TEST"}
	assert.Equal(t, "TEST", err.GetFieldName())
	err.SetFieldName("TEST2")
	assert.Equal(t, "TEST2", err.GetFieldName())
}
