package monitorables

import (
	"errors"
	"testing"

	"github.com/monitoror/monitoror/config"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/service/store"
)

type monitorableMock struct {
	variants      []coreModels.VariantName
	validateBool  bool
	validateError error
}

func (m *monitorableMock) GetDisplayName() string                { return "Monitorable mock" }
func (m *monitorableMock) GetVariants() []coreModels.VariantName { return m.variants }
func (m *monitorableMock) Validate(_ coreModels.VariantName) (bool, error) {
	return m.validateBool, m.validateError
}
func (m *monitorableMock) Enable(_ coreModels.VariantName) {}

func TestManager_EnableMonitorables(t *testing.T) {
	mockMonitorable1 := &monitorableMock{
		variants:      []coreModels.VariantName{coreModels.DefaultVariant},
		validateBool:  true,
		validateError: nil,
	}
	mockMonitorable2 := &monitorableMock{
		variants:      []coreModels.VariantName{coreModels.DefaultVariant},
		validateBool:  false,
		validateError: errors.New("boom"),
	}

	manager := NewMonitorableManager(&store.Store{
		CoreConfig: &config.Config{
			Env: "production",
		},
	})

	manager.register(mockMonitorable1)
	manager.register(mockMonitorable2)

	manager.EnableMonitorables()
}
