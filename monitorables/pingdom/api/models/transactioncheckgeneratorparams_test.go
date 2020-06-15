package models

import (
	"testing"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/test"
)

func TestTransactionCheckGeneratorParams_Validate(t *testing.T) {
	param := &CheckGeneratorParams{}
	test.AssertParams(t, param, 0)

	param = &CheckGeneratorParams{SortBy: "name"}
	test.AssertParams(t, param, 0)

	param = &CheckGeneratorParams{SortBy: "test"}
	test.AssertParams(t, param, 1)
}
