package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/test"
)

func TestChecksParams_Validate(t *testing.T) {
	param := &ChecksParams{Owner: "test", Repository: "test", Ref: "master"}
	test.AssertParams(t, param, 0)

	param = &ChecksParams{Owner: "test", Repository: "test"}
	test.AssertParams(t, param, 1)

	param = &ChecksParams{Owner: "test"}
	test.AssertParams(t, param, 2)

	param = &ChecksParams{}
	test.AssertParams(t, param, 3)
}

func TestBuildParams_String(t *testing.T) {
	param := &ChecksParams{Owner: "test", Repository: "test", Ref: "master"}
	assert.Equal(t, "CHECKS-test-test-master", fmt.Sprint(param))
}
