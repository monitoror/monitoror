package available

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/monitoror/monitoror/api/config/versions"
	"github.com/monitoror/monitoror/internal/pkg/validator"
)

type availableTestStruct struct {
	Since string `available:"since=1.0"`
	Until string `available:"until=3.2"`
	Both  string `available:"since=1.0,until=3.2"`
	Other string
}

type erroredStruct struct {
	Error string `available:"until=2017-01-09"`
}

func TestStruct(t *testing.T) {
	s := &availableTestStruct{}

	errors := Struct(s, RawVersion("0.1").ToConfigVersion())
	if assert.Len(t, errors, 2) {
		assert.Equal(t, validator.ErrorSince, errors[0].GetErrorID())
		assert.Equal(t, "Since", errors[0].GetFieldName())
		assert.Equal(t, validator.ErrorSince, errors[1].GetErrorID())
		assert.Equal(t, "Both", errors[1].GetFieldName())
	}

	assert.Len(t, Struct(s, RawVersion("1.0").ToConfigVersion()), 0)
	assert.Len(t, Struct(s, RawVersion("2.0").ToConfigVersion()), 0)
	assert.Len(t, Struct(s, RawVersion("3.2").ToConfigVersion()), 0)

	errors = Struct(s, RawVersion("4.0").ToConfigVersion())
	if assert.Len(t, errors, 2) {
		assert.Equal(t, validator.ErrorUntil, errors[0].GetErrorID())
		assert.Equal(t, "Until", errors[0].GetFieldName())
		assert.Equal(t, validator.ErrorUntil, errors[1].GetErrorID())
		assert.Equal(t, "Both", errors[1].GetFieldName())
	}
}

func TestStruct_Panic(t *testing.T) {
	assert.Panics(t, func() {
		Struct(&erroredStruct{}, RawVersion("0.1").ToConfigVersion())
	})
}
