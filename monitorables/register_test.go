package monitorables

import (
	"testing"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/test"
	coreModels "github.com/monitoror/monitoror/models"
	azureDevOpsApi "github.com/monitoror/monitoror/monitorables/azuredevops/api"
	githubApi "github.com/monitoror/monitoror/monitorables/github/api"
	httpApi "github.com/monitoror/monitoror/monitorables/http/api"
	jenkinsApi "github.com/monitoror/monitoror/monitorables/jenkins/api"
	pingApi "github.com/monitoror/monitoror/monitorables/ping/api"
	pingdomApi "github.com/monitoror/monitoror/monitorables/pingdom/api"
	portApi "github.com/monitoror/monitoror/monitorables/port/api"
	travisCIApi "github.com/monitoror/monitoror/monitorables/travisci/api"
	"github.com/monitoror/monitoror/service/registry"

	"github.com/stretchr/testify/assert"
)

func TestManager_RegisterMonitorables(t *testing.T) {
	// init Store
	store, _ := test.InitMockAndStore()
	store.Registry = registry.NewRegistry()

	manager := &Manager{store: store}
	manager.RegisterMonitorables()

	mr := store.Registry.(*registry.MetadataRegistry)

	// ------------ AZURE DEVOPS ------------
	assert.NotNil(t, mr.TileMetadata[azureDevOpsApi.AzureDevOpsBuildTileType])
	assert.NotNil(t, mr.TileMetadata[azureDevOpsApi.AzureDevOpsReleaseTileType])
	// ------------ GITHUB ------------
	assert.NotNil(t, mr.TileMetadata[githubApi.GithubCountTileType])
	assert.NotNil(t, mr.TileMetadata[githubApi.GithubChecksTileType])
	assert.NotNil(t, mr.TileMetadata[githubApi.GithubPullRequestTileType])
	assert.NotNil(t, mr.GeneratorMetadata[coreModels.NewGeneratorTileType(githubApi.GithubPullRequestTileType)])
	// ------------ HTTP ------------
	assert.NotNil(t, mr.TileMetadata[httpApi.HTTPStatusTileType])
	assert.NotNil(t, mr.TileMetadata[httpApi.HTTPRawTileType])
	assert.NotNil(t, mr.TileMetadata[httpApi.HTTPFormattedTileType])
	// ------------ JENKINS ------------
	assert.NotNil(t, mr.TileMetadata[jenkinsApi.JenkinsBuildTileType])
	assert.NotNil(t, mr.GeneratorMetadata[coreModels.NewGeneratorTileType(jenkinsApi.JenkinsBuildTileType)])
	// ------------ PING ------------
	assert.NotNil(t, mr.TileMetadata[pingApi.PingTileType])
	// ------------ PINGDOM ------------
	assert.NotNil(t, mr.TileMetadata[pingdomApi.PingdomCheckTileType])
	assert.NotNil(t, mr.GeneratorMetadata[coreModels.NewGeneratorTileType(pingdomApi.PingdomCheckTileType)])
	// ------------ PORT ------------
	assert.NotNil(t, mr.TileMetadata[portApi.PortTileType])
	// ------------ TRAVIS CI ------------
	assert.NotNil(t, mr.TileMetadata[travisCIApi.TravisCIBuildTileType])
}
