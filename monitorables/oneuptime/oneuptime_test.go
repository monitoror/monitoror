package oneuptime

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
	// OK
	_ = os.Setenv("MO_MONITORABLE_ONEUPTIME_PROJECTID", "project")
	_ = os.Setenv("MO_MONITORABLE_ONEUPTIME_TOKEN", "xxx")
	// Missing Token
	_ = os.Setenv("MO_MONITORABLE_ONEUPTIME_VARIANT1_PROJECTID", "project")
	// Url broken
	_ = os.Setenv("MO_MONITORABLE_ONEUPTIME_VARIANT2_URL", "url%soneuptime.example.com")
	_ = os.Setenv("MO_MONITORABLE_ONEUPTIME_VARIANT2_PROJECTID", "project")
	_ = os.Setenv("MO_MONITORABLE_ONEUPTIME_VARIANT2_TOKEN", "xxx")
	// No project and no token
	_ = os.Setenv("MO_MONITORABLE_ONEUPTIME_VARIANT3_TIMEOUT", "1000")

	// NewMonitorable
	monitorable := NewMonitorable(store)
	assert.NotNil(t, monitorable)

	// GetDisplayName
	assert.NotNil(t, monitorable.GetDisplayName())

	// GetVariantsNames and check
	if assert.Len(t, monitorable.GetVariantsNames(), 4) {
		_, errors := monitorable.Validate("variant1")
		assert.NotEmpty(t, errors)
		_, errors = monitorable.Validate("variant2")
		assert.NotEmpty(t, errors)
		valid, errors := monitorable.Validate("variant3")
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
	mockMonitorableHelper.RouterAssertNumberOfCalls(t, 1, 1)
	mockMonitorableHelper.TileSettingsManagerAssertNumberOfCalls(t, 1, 1, 1, 1)
}
