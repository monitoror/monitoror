package models

import (
	"fmt"
	"testing"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/test"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
)

func TestBuildParams_Validate(t *testing.T) {
	param := &BuildParams{}
	test.AssertParams(t, param, 2)

	param.Project = "test"
	test.AssertParams(t, param, 1)

	param.Definition = pointer.ToInt(1)
	test.AssertParams(t, param, 0)

	param.Branch = pointer.ToString("test")
	test.AssertParams(t, param, 0)
}

func TestBuildParams_String(t *testing.T) {
	param := &BuildParams{
		Project:    "test",
		Definition: pointer.ToInt(1),
	}
	assert.Equal(t, "BUILD-test-1", fmt.Sprint(param))

	param.Branch = pointer.ToString("test")
	assert.Equal(t, "BUILD-test-1-test", fmt.Sprint(param))
}
