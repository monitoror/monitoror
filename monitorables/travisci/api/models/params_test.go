package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/test"
)

func TestBuildParams_Validate(t *testing.T) {
	param := &BuildParams{Owner: "test", Repository: "test", Branch: "master"}
	test.AssertParams(t, param, 0)

	param = &BuildParams{Owner: "test", Repository: "test"}
	test.AssertParams(t, param, 1)

	param = &BuildParams{Owner: "test"}
	test.AssertParams(t, param, 2)

	param = &BuildParams{}
	test.AssertParams(t, param, 3)
}

func TestBuildParams_String(t *testing.T) {
	param := &BuildParams{Repository: "test", Owner: "test", Branch: "test"}
	assert.Equal(t, "BUILD-test-test-test", fmt.Sprint(param))
}
