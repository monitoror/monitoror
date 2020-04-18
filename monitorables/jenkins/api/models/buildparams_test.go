package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/test"
)

func TestBuildParams_Validate(t *testing.T) {
	param := &BuildParams{Job: "test", Branch: "test"}
	test.AssertParams(t, param, 0)

	param = &BuildParams{Job: "test"}
	test.AssertParams(t, param, 0)

	param = &BuildParams{}
	test.AssertParams(t, param, 1)
}

func TestBuildParams_String(t *testing.T) {
	param := &BuildParams{Job: "test", Branch: "test"}
	assert.Equal(t, "BUILD-test-test", fmt.Sprint(param))
}
