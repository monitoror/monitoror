package monitorables

import (
	"github.com/monitoror/monitoror/monitorables/azuredevops"
	"github.com/monitoror/monitoror/monitorables/github"
	"github.com/monitoror/monitoror/monitorables/http"
	"github.com/monitoror/monitoror/monitorables/jenkins"
	"github.com/monitoror/monitoror/monitorables/ping"
	"github.com/monitoror/monitoror/monitorables/pingdom"
	"github.com/monitoror/monitoror/monitorables/port"
	"github.com/monitoror/monitoror/monitorables/travisci"
)

func (m *Manager) init() {
	// ------------ AZURE DEVOPS ------------
	m.register(azuredevops.NewMonitorable(m.store))
	// ------------ GITHUB ------------
	m.register(github.NewMonitorable(m.store))
	// ------------ HTTP ------------
	m.register(http.NewMonitorable(m.store))
	// ------------ JENKINS ------------
	m.register(jenkins.NewMonitorable(m.store))
	// ------------ PING ------------
	m.register(ping.NewMonitorable(m.store))
	// ------------ PINGDOM ------------
	m.register(pingdom.NewMonitorable(m.store))
	// ------------ PORT ------------
	m.register(port.NewMonitorable(m.store))
	// ------------ TRAVIS CI ------------
	m.register(travisci.NewMonitorable(m.store))
}
