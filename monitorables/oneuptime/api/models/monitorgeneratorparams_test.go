package models

import (
	"testing"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/test"
)

func TestMonitorGeneratorParams_Validate(t *testing.T) {
	param := &MonitorGeneratorParams{}
	test.AssertParams(t, param, 0)

	param = &MonitorGeneratorParams{SortBy: "name"}
	test.AssertParams(t, param, 0)

	param = &MonitorGeneratorParams{SortBy: "test"}
	test.AssertParams(t, param, 1)
}
