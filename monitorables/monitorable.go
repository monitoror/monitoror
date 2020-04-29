package monitorables

import (
	"github.com/monitoror/monitoror/cli"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/store"
)

type Monitorable interface {
	//GetDisplayName return monitorable name display in console
	GetDisplayName() string

	//GetVariantsNames return variant list extract from config
	GetVariantsNames() []coreModels.VariantName

	//Validate test if config variant is valid
	// return false if empty and error if config have an error (ex: wrong url format)
	Validate(variantName coreModels.VariantName) (bool, []error)

	//Enable monitorable variant (add route to echo and enable tile for config verify / hydrate)
	Enable(variantName coreModels.VariantName)
}

type (
	Manager struct {
		store *store.Store

		monitorables []Monitorable
	}
)

func NewMonitorableManager(store *store.Store) *Manager {
	return &Manager{store: store}
}

func (m *Manager) register(monitorable Monitorable) {
	m.monitorables = append(m.monitorables, monitorable)
}

func (m *Manager) EnableMonitorables() {
	m.store.Cli.PrintMonitorableHeader()

	nonEnabledMonitorableCount := 0

	for _, monitorable := range m.monitorables {
		var enabledVariants []coreModels.VariantName
		var erroredVariants []cli.ErroredVariant

		for _, variantName := range monitorable.GetVariantsNames() {
			valid, err := monitorable.Validate(variantName)
			if err != nil {
				erroredVariants = append(erroredVariants, cli.ErroredVariant{VariantName: variantName, Errors: err})
			}

			if valid {
				monitorable.Enable(variantName)
				enabledVariants = append(enabledVariants, variantName)
			}
		}

		if len(enabledVariants) == 0 && len(erroredVariants) == 0 {
			nonEnabledMonitorableCount++
		}

		m.store.Cli.PrintMonitorable(monitorable.GetDisplayName(), enabledVariants, erroredVariants)
	}

	m.store.Cli.PrintMonitorableFooter(m.store.CoreConfig.Env == "production", nonEnabledMonitorableCount)
}
