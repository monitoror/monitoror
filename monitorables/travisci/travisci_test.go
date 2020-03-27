package travisci

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
	// Url broken
	_ = os.Setenv("MO_MONITORABLE_TRAVISCI_VARIANT0_URL", "url%stravis.example.com")

	// NewMonitorable
	monitorable := NewMonitorable(s)
	assert.NotNil(t, monitorable)

	// GetDisplayName
	assert.NotNil(t, monitorable.GetDisplayName())

	// GetVariants and check
	if assert.Len(t, monitorable.GetVariants(), 2) {
		_, err := monitorable.Validate("variant0")
		assert.Error(t, err)
	}

	// Enable
	for _, variant := range monitorable.GetVariants() {
		if valid, _ := monitorable.Validate(variant); valid {
			monitorable.Enable(variant)
		}
	}

	// Test calls
	mockRouter.AssertNumberOfCalls(t, "Group", 1)
	mockRouterGroup.AssertNumberOfCalls(t, "GET", 1)
	mockConfigManager.AssertNumberOfCalls(t, "RegisterTile", 1)
	mockConfigManager.AssertNumberOfCalls(t, "EnableTile", 1)
	mockConfigManager.AssertNumberOfCalls(t, "EnableDynamicTile", 0)
}
