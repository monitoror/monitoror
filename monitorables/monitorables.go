package monitorables

import (
	"github.com/monitoror/monitoror/monitorables/azuredevops"
	"github.com/monitoror/monitoror/monitorables/github"
	"github.com/monitoror/monitoror/monitorables/gitlab"
	"github.com/monitoror/monitoror/monitorables/http"
	"github.com/monitoror/monitoror/monitorables/jenkins"
	"github.com/monitoror/monitoror/monitorables/ping"
	"github.com/monitoror/monitoror/monitorables/pingdom"
	"github.com/monitoror/monitoror/monitorables/port"
	"github.com/monitoror/monitoror/monitorables/travisci"
	"github.com/monitoror/monitoror/store"
)

func RegisterMonitorables(s *store.Store) {
	// ------------ AZURE DEVOPS ------------
	s.Registry.RegisterMonitorable(azuredevops.NewMonitorable(s))
	// ------------ GITHUB ------------
	s.Registry.RegisterMonitorable(github.NewMonitorable(s))
	// ------------ GITLAB ------------
	s.Registry.RegisterMonitorable(gitlab.NewMonitorable(s))
	// ------------ HTTP ------------
	s.Registry.RegisterMonitorable(http.NewMonitorable(s))
	// ------------ JENKINS ------------
	s.Registry.RegisterMonitorable(jenkins.NewMonitorable(s))
	// ------------ PING ------------
	s.Registry.RegisterMonitorable(ping.NewMonitorable(s))
	// ------------ PINGDOM ------------
	s.Registry.RegisterMonitorable(pingdom.NewMonitorable(s))
	// ------------ PORT ------------
	s.Registry.RegisterMonitorable(port.NewMonitorable(s))
	// ------------ TRAVIS CI ------------
	s.Registry.RegisterMonitorable(travisci.NewMonitorable(s))
}
