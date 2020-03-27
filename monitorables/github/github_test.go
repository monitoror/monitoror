package github

import (
	"os"
	"testing"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/test"
	"github.com/stretchr/testify/assert"
)

func TestNewMonitorable(t *testing.T) {
	// init Store
	mockRouter, mockRouterGroup, mockConfigManager, s := test.InitMockAndStore()

	// init Env
	_ = os.Setenv("MO_MONITORABLE_GITHUB_VARIANT0_TOKEN", "xxx")

	// NewMonitorable
	monitorable := NewMonitorable(s)
	assert.NotNil(t, monitorable)

	// GetDisplayName
	assert.NotNil(t, monitorable.GetDisplayName())

	// GetVariants and check
	assert.Len(t, monitorable.GetVariants(), 2)

	// Enable
	for _, variant := range monitorable.GetVariants() {
		if valid, _ := monitorable.Validate(variant); valid {
			monitorable.Enable(variant)
		}
	}

	// Test calls
	mockRouter.AssertNumberOfCalls(t, "Group", 1)
	mockRouterGroup.AssertNumberOfCalls(t, "GET", 2)
	mockConfigManager.AssertNumberOfCalls(t, "RegisterTile", 3)
	mockConfigManager.AssertNumberOfCalls(t, "EnableTile", 2)
	mockConfigManager.AssertNumberOfCalls(t, "EnableDynamicTile", 1)
}
