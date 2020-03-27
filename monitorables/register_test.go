package monitorables

import (
	"testing"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/test"
)

func TestManager_RegisterMonitorables(t *testing.T) {
	// init Store
	_, _, mockConfigManager, s := test.InitMockAndStore()
	manager := &Manager{store: s}
	manager.RegisterMonitorables()

	tileTypeCount := 15
	mockConfigManager.AssertNumberOfCalls(t, "RegisterTile", tileTypeCount)
}
