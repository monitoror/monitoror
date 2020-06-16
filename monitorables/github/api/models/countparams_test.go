package models

import (
	"testing"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/test"
)

func TestCountParams_Validate(t *testing.T) {
	param := &CountParams{Query: "test"}
	test.AssertParams(t, param, 0)

	param = &CountParams{}
	test.AssertParams(t, param, 1)
}
