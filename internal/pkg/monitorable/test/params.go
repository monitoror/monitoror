package test

import (
	"testing"

	"github.com/monitoror/monitoror/api/config/versions"
	"github.com/monitoror/monitoror/internal/pkg/monitorable/params"
	"github.com/monitoror/monitoror/internal/pkg/validator"
	"github.com/monitoror/monitoror/internal/pkg/validator/available"
	"github.com/monitoror/monitoror/internal/pkg/validator/validate"

	"github.com/stretchr/testify/assert"
)

func AssertParams(t *testing.T, p params.Validator, errorCount int, optionalConfigVersion ...*versions.ConfigVersion) {
	var errors []validator.Error

	if len(optionalConfigVersion) == 1 {
		// use "available" tag in struct definition to validate params
		errors = append(errors, available.Struct(p, optionalConfigVersion[0])...)
	}

	// use "validate" tag in struct definition to validate params
	errors = append(errors, validate.Struct(p)...)
	// use "Validate" function to to validate params
	errors = append(errors, p.Validate()...)

	assert.Len(t, errors, errorCount)
}
