package models

import (
	"testing"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/test"
)

func TestMonitorParams_Validate(t *testing.T) {
	param := &MonitorParams{}
	test.AssertParams(t, param, 1)

	param = &MonitorParams{ID: "8176df0f-5059-402b-a3d2-f7821257d724"}
	test.AssertParams(t, param, 0)
}
