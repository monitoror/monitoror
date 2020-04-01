package monitorables

import (
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/service/store"
)

type Monitorable interface {
	//GetDisplayName return monitorable name display in console
	GetDisplayName() string

	//GetVariants return variant list extract from config
	GetVariants() []coreModels.VariantName

	//Validate test if config variant is valid
	// return false if empty and error if config have an error (ex: wrong url format)
	Validate(variant coreModels.VariantName) (bool, error)

	//Enable monitorable variant (add route to echo and enable tile for config verify / hydrate)
	Enable(variant coreModels.VariantName)
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
		erroredVariants := make(map[coreModels.VariantName]error)

		for _, variant := range monitorable.GetVariants() {
			valid, err := monitorable.Validate(variant)
			if err != nil {
				erroredVariants[variant] = err
			}

			if valid {
				monitorable.Enable(variant)
				enabledVariants = append(enabledVariants, variant)
			}
		}

		if len(enabledVariants) == 0 && len(erroredVariants) == 0 {
			nonEnabledMonitorableCount++
		}

		m.store.Cli.PrintMonitorable(monitorable.GetDisplayName(), enabledVariants, erroredVariants)
	}

	m.store.Cli.PrintMonitorableFooter(m.store.CoreConfig.Env == "production", nonEnabledMonitorableCount)
}
