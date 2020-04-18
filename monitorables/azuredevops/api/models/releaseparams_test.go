package models

import (
	"fmt"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/test"
)

func TestReleaseParams_Validate(t *testing.T) {
	param := &ReleaseParams{}
	test.AssertParams(t, param, 2)

	param.Project = "test"
	test.AssertParams(t, param, 1)

	param.Definition = pointer.ToInt(1)
	test.AssertParams(t, param, 0)
}

func TestReleaseParams_String(t *testing.T) {
	param := &ReleaseParams{
		Project:    "test",
		Definition: pointer.ToInt(1),
	}
	assert.Equal(t, "RELEASE-test-1", fmt.Sprint(param))
}
