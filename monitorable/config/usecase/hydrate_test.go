package usecase

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/monitoror/monitoror/monitorable/config/repository"
	"github.com/monitoror/monitoror/monitorable/jenkins"

	"github.com/stretchr/testify/assert"
)

func TestUsecase_Hydrate(t *testing.T) {
	input := `
{
  "columns": 4,
  "tiles": [
    { "type": "EMPTY" },
    { "type": "PING", "params": { "hostname": "aserver.com" } },
    { "type": "PORT", "params": { "hostname": "bserver.com", "port": 22 } },
    { "type": "GROUP", "label": "...", "tiles": [
      { "type": "PING", "params": { "hostname": "aserver.com" } },
      { "type": "PORT", "params": { "hostname": "bserver.com", "port": 22 } }
    ]},
		{ "type": "JENKINS-BUILD", "params": { "job": "test" } },
		{ "type": "JENKINS-BUILD", "configVariant": "variant1", "params": { "job": "test" } }
  ]
}
`

	usecase := initConfigUsecase()
	reader := ioutil.NopCloser(strings.NewReader(input))
	config, err := repository.GetConfig(reader)
	assert.NoError(t, err)

	err = usecase.Hydrate(config, "http://localhost:8080")
	assert.NoError(t, err)

	assert.Equal(t, "http://localhost:8080/ping?hostname=aserver.com", config.Tiles[1].Url)
	assert.Equal(t, "http://localhost:8080/port?hostname=bserver.com&port=22", config.Tiles[2].Url)

	group := config.Tiles[3].Tiles
	assert.Equal(t, "http://localhost:8080/ping?hostname=aserver.com", group[0].Url)
	assert.Equal(t, "http://localhost:8080/port?hostname=bserver.com&port=22", group[1].Url)

	assert.Equal(t, "http://localhost:8080/jenkins/default?job=test", config.Tiles[4].Url)
	assert.Equal(t, "http://localhost:8080/jenkins/variant1?job=test", config.Tiles[5].Url)
}

func TestUsecase_Hydrate_WithDynamic(t *testing.T) {
	input := `
{
  "columns": 4,
  "tiles": [
    { "type": "JENKINS-MULTIBRANCH" },
    { "type": "GROUP", "label": "...", "tiles": [
      { "type": "PING", "params": { "hostname": "aserver.com" } },
			{ "type": "JENKINS-MULTIBRANCH" }
    ]},
    { "type": "GROUP", "label": "...", "tiles": [
    	{ "type": "JENKINS-MULTIBRANCH"}
    ]}
  ]
}
`
	usecase := initConfigUsecase()

	reader := ioutil.NopCloser(strings.NewReader(input))
	config, err := repository.GetConfig(reader)
	assert.NoError(t, err)

	err = usecase.Hydrate(config, "http://localhost:8080")
	assert.NoError(t, err)

	assert.Equal(t, 3, len(config.Tiles))
	assert.Equal(t, jenkins.JenkinsBuildTileType, config.Tiles[0].Type)
	assert.Equal(t, "http://localhost:8080/jenkins/default?job=test", config.Tiles[0].Url)
	assert.Equal(t, jenkins.JenkinsBuildTileType, config.Tiles[1].Tiles[1].Type)
	assert.Equal(t, "http://localhost:8080/jenkins/default?job=test", config.Tiles[1].Tiles[1].Url)
	assert.Equal(t, jenkins.JenkinsBuildTileType, config.Tiles[2].Tiles[0].Type)
	assert.Equal(t, "http://localhost:8080/jenkins/default?job=test", config.Tiles[2].Tiles[0].Url)
}

func TestUsecase_Hydrate_WithDynamic_WithError(t *testing.T) {
	input := `
{
  "columns": 4,
  "tiles": [
    { "type": "JENKINS-MULTIBRANCH"},
    { "type": "GROUP", "label": "...", "tiles": [
      { "type": "PING", "params": { "hostname": "aserver.com" } },
			{ "type": "JENKINS-MULTIBRANCH" }
    ]},
    { "type": "GROUP", "label": "...", "tiles": [
    	{ "type": "JENKINS-MULTIBRANCH", "configVariant": "variant1"}
    ]}
  ]
}
`
	usecase := initConfigUsecase()

	reader := ioutil.NopCloser(strings.NewReader(input))
	config, err := repository.GetConfig(reader)
	assert.NoError(t, err)

	err = usecase.Hydrate(config, "http://localhost:8080")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), `Error while listing JENKINS-MULTIBRANCH dynamic tiles.`)
}
