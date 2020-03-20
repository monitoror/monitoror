package monitorables

import (
	"github.com/monitoror/monitoror/service/store"
)

type Monitorable interface {
	GetVariants() []string

	IsValid(variant string) bool
	Enable(variant string)
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
