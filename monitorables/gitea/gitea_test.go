package gitea

import (
	"os"
	"testing"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/test"

	"github.com/stretchr/testify/assert"
)

func TestNewMonitorable(t *testing.T) {
	// init Store
	store, mockMonitorableHelper := test.InitMockAndStore()

	// init Env
	// OK (default URL with token)
	_ = os.Setenv("MO_MONITORABLE_GITEA_TOKEN", "xxx")
	// OK (self-hosted instance with anonymous access)
	_ = os.Setenv("MO_MONITORABLE_GITEA_VARIANT1_URL", "https://gitea.example.com/")
	// Url broken
	_ = os.Setenv("MO_MONITORABLE_GITEA_VARIANT2_URL", "url%sgitea.example.com/")
	// No URL and no token: not configured
	_ = os.Setenv("MO_MONITORABLE_GITEA_VARIANT3_TIMEOUT", "1000")

	// NewMonitorable
	monitorable := NewMonitorable(store)
	assert.NotNil(t, monitorable)

	// GetDisplayName
	assert.NotNil(t, monitorable.GetDisplayName())

	// GetVariantsNames and check
	if assert.Len(t, monitorable.GetVariantsNames(), 4) {
		valid, errors := monitorable.Validate("variant1")
		assert.True(t, valid)
		assert.Empty(t, errors)
		_, errors = monitorable.Validate("variant2")
		assert.NotEmpty(t, errors)
		valid, errors = monitorable.Validate("variant3")
		assert.False(t, valid)
		assert.Empty(t, errors)
	}

	// Enable
	for _, variantName := range monitorable.GetVariantsNames() {
		if valid, _ := monitorable.Validate(variantName); valid {
			monitorable.Enable(variantName)
		}
	}

	// Test calls
	mockMonitorableHelper.RouterAssertNumberOfCalls(t, 2, 6)
	mockMonitorableHelper.TileSettingsManagerAssertNumberOfCalls(t, 3, 1, 6, 2)
}
