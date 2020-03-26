package ping

import (
	"testing"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/test"
	"github.com/stretchr/testify/assert"
)

func TestNewMonitorable(t *testing.T) {
	// init Store
	mockRouter, mockRouterGroup, mockConfigManager, s := test.InitMockAndStore()

	// NewMonitorable
	monitorable := NewMonitorable(s)
	assert.NotNil(t, monitorable)

	// GetDisplayName
	assert.NotNil(t, monitorable.GetDisplayName())

	// GetVariants and check
	assert.Len(t, monitorable.GetVariants(), 1)

	// Enable
	for _, variant := range monitorable.GetVariants() {
		_, _ = monitorable.Validate(variant) // Skip validate because is always false in test
		monitorable.Enable(variant)
	}

	// Test calls
	mockRouter.AssertNumberOfCalls(t, "Group", 1)
	mockRouterGroup.AssertNumberOfCalls(t, "GET", 1)
	mockConfigManager.AssertNumberOfCalls(t, "RegisterTile", 1)
	mockConfigManager.AssertNumberOfCalls(t, "EnableTile", 1)
	mockConfigManager.AssertNumberOfCalls(t, "EnableDynamicTile", 0)
}
