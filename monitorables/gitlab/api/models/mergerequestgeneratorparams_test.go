package models

import (
	"testing"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/test"

	"github.com/AlekSi/pointer"
)

func TestMergeRequestGeneratorParams_Validate(t *testing.T) {
	param := &MergeRequestGeneratorParams{ProjectID: pointer.ToInt(10)}
	test.AssertParams(t, param, 0)

	param = &MergeRequestGeneratorParams{}
	test.AssertParams(t, param, 1)
}
