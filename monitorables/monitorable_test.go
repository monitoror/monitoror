package monitorables

import (
	"errors"
	"testing"

	cliMocks "github.com/monitoror/monitoror/cli/mocks"
	"github.com/monitoror/monitoror/config"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/service/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type monitorableMock struct {
	displayName   string
	variants      []coreModels.VariantName
	validateBool  bool
	validateError error
}

func (m *monitorableMock) GetDisplayName() string                     { return m.displayName }
func (m *monitorableMock) GetVariantsNames() []coreModels.VariantName { return m.variants }
func (m *monitorableMock) Validate(_ coreModels.VariantName) (bool, error) {
	return m.validateBool, m.validateError
}
func (m *monitorableMock) Enable(_ coreModels.VariantName) {}

func TestManager_EnableMonitorables(t *testing.T) {
	cliMock := new(cliMocks.CLI)
	cliMock.On("PrintMonitorableHeader")
	cliMock.On("PrintMonitorable",
		mock.AnythingOfType("string"),
		mock.Anything,
		mock.Anything,
	)
	cliMock.On("PrintMonitorableFooter",
		mock.AnythingOfType("bool"),
		mock.AnythingOfType("int"),
	)

	mockMonitorable1 := &monitorableMock{
		displayName:   "Monitorable mock 1",
		variants:      []coreModels.VariantName{coreModels.DefaultVariant},
		validateBool:  true,
		validateError: nil,
	}
	mockMonitorable2 := &monitorableMock{
		displayName:   "Monitorable mock 2 (faker)",
		variants:      []coreModels.VariantName{coreModels.DefaultVariant},
		validateBool:  false,
		validateError: errors.New("boom"),
	}

	manager := NewMonitorableManager(&store.Store{
		CoreConfig: &config.Config{
			Env: "production",
		},
		Cli: cliMock,
	})

	manager.register(mockMonitorable1)
	assert.Len(t, manager.monitorables, 1)
	manager.register(mockMonitorable2)
	assert.Len(t, manager.monitorables, 2)

	manager.EnableMonitorables()

	cliMock.AssertNumberOfCalls(t, "PrintMonitorableHeader", 1)
	cliMock.AssertNumberOfCalls(t, "PrintMonitorable", 2)
	cliMock.AssertNumberOfCalls(t, "PrintMonitorableFooter", 1)
	cliMock.AssertCalled(t, "PrintMonitorableFooter", true, 0)

	// Count non-enabled monitorables
	mockMonitorable3 := &monitorableMock{
		displayName:   "Monitorable mock 3",
		variants:      []coreModels.VariantName{},
		validateBool:  false,
		validateError: nil,
	}
	manager.register(mockMonitorable3)
	assert.Len(t, manager.monitorables, 3)

	manager.EnableMonitorables()
	cliMock.AssertCalled(t, "PrintMonitorableFooter", true, 1)
}
