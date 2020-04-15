package azuredevops

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
	_ = os.Setenv("MO_MONITORABLE_AZUREDEVOPS_VARIANT0_URL", "https://azure.example.com/myProject1")
	_ = os.Setenv("MO_MONITORABLE_AZUREDEVOPS_VARIANT0_TOKEN", "xxx")
	// Missing Token
	_ = os.Setenv("MO_MONITORABLE_AZUREDEVOPS_VARIANT1_URL", "https://azure.example.com/myProject2")
	// Url broken
	_ = os.Setenv("MO_MONITORABLE_AZUREDEVOPS_VARIANT2_URL", "url%sazure.example.com/myProject2")

	// NewMonitorable
	monitorable := NewMonitorable(store)
	assert.NotNil(t, monitorable)

	// GetDisplayName
	assert.NotNil(t, monitorable.GetDisplayName())

	// GetVariantsNames and check
	if assert.Len(t, monitorable.GetVariantsNames(), 4) {
		_, err := monitorable.Validate("variant1")
		assert.Error(t, err)
		_, err = monitorable.Validate("variant2")
		assert.Error(t, err)
	}

	// Enable
	for _, variantName := range monitorable.GetVariantsNames() {
		if valid, _ := monitorable.Validate(variantName); valid {
			monitorable.Enable(variantName)
		}
	}

	// Test calls
	mockMonitorableHelper.RouterAssertNumberOfCalls(t, 1, 2)
	mockMonitorableHelper.TileSettingsManagerAssertNumberOfCalls(t, 2, 0, 2, 0)
}
