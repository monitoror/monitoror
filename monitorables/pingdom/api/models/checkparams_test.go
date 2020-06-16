package models

import (
	"testing"

	"github.com/AlekSi/pointer"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/test"
)

func TestCheckParams_Validate(t *testing.T) {
	param := &CheckParams{}
	test.AssertParams(t, param, 1)

	param = &CheckParams{ID: pointer.ToInt(10)}
	test.AssertParams(t, param, 0)
}
