package http

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
		if valid, _ := monitorable.Validate(variant); valid {
			monitorable.Enable(variant)
		}
	}

	// Test calls
	mockRouter.AssertNumberOfCalls(t, "Group", 1)
	mockRouterGroup.AssertNumberOfCalls(t, "GET", 3)
	mockConfigManager.AssertNumberOfCalls(t, "RegisterTile", 3)
	mockConfigManager.AssertNumberOfCalls(t, "EnableTile", 3)
	mockConfigManager.AssertNumberOfCalls(t, "EnableDynamicTile", 0)
}
