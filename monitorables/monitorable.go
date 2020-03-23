package monitorables

import (
	"github.com/monitoror/monitoror/cli"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/service/store"
)

type Monitorable interface {
	//GetDisplayName return monitorable name display in console
	GetDisplayName() string

	//GetVariants return variantlist extract from config
	GetVariants() []coreModels.Variant

	//Validate test if config variant is valid
	// return false if empty and error if config have an error (ex: wrong url format)
	Validate(variant coreModels.Variant) (bool, error)

	//Enable monitorable variant (add route to echo and enable tile for config verify / hydrate)
	Enable(variant coreModels.Variant)
}

type (
	Manager struct {
		store *store.Store

		monitorables []Monitorable
	}
)

func NewMonitorableManager(store *store.Store) *Manager {
	manager := &Manager{store: store}

	// Register all monitorables
	manager.init()

	return manager
}

func (m *Manager) register(monitorable Monitorable) {
	m.monitorables = append(m.monitorables, monitorable)
}

func (m *Manager) EnableMonitorables() {
	//TODO: LOGS

	for _, monitorable := range m.monitorables {
		for _, variant := range monitorable.GetVariants() {
			monitorable.Enable(variant)
		}
	}
}
