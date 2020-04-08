package monitorables

import (
	"testing"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/test"
)

func TestManager_RegisterMonitorables(t *testing.T) {
	// init Store
	store, mockMonitorableHelper := test.InitMockAndStore()
	manager := &Manager{store: store}
	manager.RegisterMonitorables()

	tileTypeCount := 12
	tileGeneratorCount := 3
	mockMonitorableHelper.TileSettingsManagerAssertNumberOfCalls(t, tileTypeCount, tileGeneratorCount, 0, 0)
}
