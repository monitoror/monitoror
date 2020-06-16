package models

import (
	"testing"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/test"
)

func TestBuildGeneratorParams_Validate(t *testing.T) {
	param := &BuildGeneratorParams{
		Job:     "test",
		Match:   ".*",
		Unmatch: "master",
	}
	test.AssertParams(t, param, 0)

	param = &BuildGeneratorParams{}
	test.AssertParams(t, param, 1)

	param = &BuildGeneratorParams{Job: "test", Match: "("}
	test.AssertParams(t, param, 1)

	param = &BuildGeneratorParams{Job: "test", Unmatch: "("}
	test.AssertParams(t, param, 1)
}
