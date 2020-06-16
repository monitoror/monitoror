package models

import (
	"testing"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/test"
)

func TestPingParams_Validate(t *testing.T) {
	param := &PingParams{Hostname: "test"}
	test.AssertParams(t, param, 0)

	param = &PingParams{}
	test.AssertParams(t, param, 1)
}
